package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chaitanyasharma/DBs/go-db-lite/internal/ansi"
)

type Column struct {
	Name         [64]byte // Fixed-size field for name
	DataType     DataType
	Nullable     bool
	IsPrimaryKey bool
	Size         int
}

func (c *Column) WriteTo(w io.Writer) (int64, error) {
	var written int64
	n, err := w.Write(c.Name[:])
	written += int64(n)
	if err != nil {
		return written, err
	}
	if err := binary.Write(w, binary.LittleEndian, c.DataType); err != nil {
		return written, err
	}
	return written, binary.Write(w, binary.LittleEndian, c.Nullable)
}

func (c *Column) ReadFrom(r io.Reader) (int64, error) {
	var read int64
	if n, err := io.ReadFull(r, c.Name[:]); err != nil {
		return read, err
	} else {
		read += int64(n)
	}
	if err := binary.Read(r, binary.LittleEndian, &c.DataType); err != nil {
		return read, err
	}
	read += int64(binary.Size(c.DataType))
	return read, binary.Read(r, binary.LittleEndian, &c.Nullable)
}

type Row struct {
	Values []byte // Store serialized values as a byte slice
}

func (r *Row) WriteTo(w io.Writer) (int64, error) {
	var written int64
	n, err := w.Write(r.Values)
	written += int64(n)
	if err != nil {
		return written, err
	}
	return written, binary.Write(w, binary.LittleEndian, r.Values)
}

func (r *Row) ReadFrom(reader io.Reader) (int64, error) {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return 0, err
	}
	r.Values = make([]byte, length)
	_, err := io.ReadFull(reader, r.Values)
	return 0, err
}

type Table struct {
	Name        [64]byte // Fixed-size field for name
	Columns     []Column
	Rows        []Row
	ColumnCount int
	RowCount    int
	Metadata    map[string]string
}

func (t *Table) AddColumn(name string, dataType DataType, nullable bool) {
	var colName [64]byte
	copy(colName[:], name)
	t.Columns = append(t.Columns, Column{Name: colName, DataType: dataType, Nullable: nullable})
}

func (t *Table) AddRow(values []interface{}) error {
	if len(values) != len(t.Columns) {
		return fmt.Errorf("number of values (%d) does not match number of columns (%d)", len(values), len(t.Columns))
	}

	serializedValues, err := serializeValues(values, t.Columns)
	if err != nil {
		return fmt.Errorf("error serializing values: %v", err)
	}

	t.Rows = append(t.Rows, Row{Values: serializedValues})
	return nil
}

func (t *Table) GetColumnNames() []string {
	names := make([]string, len(t.Columns))
	for i, col := range t.Columns {
		names[i] = strings.TrimRight(string(col.Name[:]), "\x00")
	}
	return names
}

func (t *Table) GetRow(index int) Row {
	return t.Rows[index]
}

func (t *Table) GetAllRows() []Row {
	return t.Rows
}

func (t *Table) CreateTable(tableName string, columns []Column) {
	copy(t.Name[:], tableName)
	t.Columns = columns
	t.Rows = []Row{}
}

func (t *Table) WriteTo(w io.Writer) (int64, error) {
	var written int64
	// Write table name
	if n, err := w.Write(t.Name[:]); err != nil {
		return written, err
	} else {
		written += int64(n)
	}
	if _, err := w.Write(t.Name[:]); err != nil {
		return written, err
	}

	// Write number of columns
	if err := binary.Write(w, binary.LittleEndian, uint32(len(t.Columns))); err != nil {
		return written, err
	}

	// Write columns
	for _, col := range t.Columns {
		if _, err := col.WriteTo(w); err != nil {
			return written, err
		}
	}

	// Write number of rows
	if err := binary.Write(w, binary.LittleEndian, uint32(len(t.Rows))); err != nil {
		return written, err
	}

	// Write rows
	for _, row := range t.Rows {
		if _, err := row.WriteTo(w); err != nil {
			return written, err
		}
	}

	return written, nil
}

func (t *Table) ReadFrom(r io.Reader) (int64, error) {
	var read int64
	// Read table name
	if n, err := io.ReadFull(r, t.Name[:]); err != nil {
		return read, err
	} else {
		read += int64(n)
	}

	// Read number of columns
	var colCount uint32
	if err := binary.Read(r, binary.LittleEndian, &colCount); err != nil {
		return read, err
	}
	read += int64(binary.Size(colCount))

	// Read columns
	t.Columns = make([]Column, colCount)
	for i := range t.Columns {
		if _, err := t.Columns[i].ReadFrom(r); err != nil {
			return read, err
		}
	}

	// Read number of rows
	var rowCount uint32
	if err := binary.Read(r, binary.LittleEndian, &rowCount); err != nil {
		return read, err
	}
	read += int64(binary.Size(rowCount))

	// Read rows
	t.Rows = make([]Row, rowCount)
	for i := range t.Rows {
		if _, err := t.Rows[i].ReadFrom(r); err != nil {
			return read, err
		}
	}

	return read, nil
}

func serializeValues(values []interface{}, columns []Column) ([]byte, error) {
	var buf bytes.Buffer
	for i, value := range values {
		dataTypeStr := columns[i].DataType.GetDataTypeString()
		buf.WriteString(dataTypeStr)
		buf.WriteByte(0) // Null terminator for the string

		if value == nil {
			buf.WriteByte(0xFF) // Indicator for null value
			continue
		}

		switch v := value.(type) {
		case int32:
			binary.Write(&buf, binary.LittleEndian, v)
		case string:
			binary.Write(&buf, binary.LittleEndian, uint16(len(v)))
			buf.WriteString(v)
		case bool:
			if v {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		case float32:
			binary.Write(&buf, binary.LittleEndian, v)
		default:
			return nil, fmt.Errorf("unsupported type: %T", v)
		}
	}
	return buf.Bytes(), nil
}

func deserializeValues(serialized []byte, columns []Column) ([]interface{}, error) {
	var values []interface{}
	offset := 0
	for _, col := range columns {
		if offset >= len(serialized) {
			return nil, fmt.Errorf("unexpected end of data for column %s", strings.TrimRight(string(col.Name[:]), "\x00"))
		}

		// Read data type string
		end := bytes.IndexByte(serialized[offset:], 0)
		if end == -1 {
			return nil, fmt.Errorf("malformed data: no null terminator for data type string")
		}
		dataTypeStr := string(serialized[offset : offset+end])
		offset += end + 1 // +1 for null terminator

		dataType := GetDataTypeFromString(dataTypeStr)
		if dataType == SQL_TYPE_UNKNOWN {
			return nil, fmt.Errorf("unknown data type: %s", dataTypeStr)
		}

		// Check for null value
		if serialized[offset] == 0xFF {
			values = append(values, nil)
			offset++
			continue
		}

		switch dataType {
		case SQL_TYPE_INT:
			value := int32(binary.LittleEndian.Uint32(serialized[offset:]))
			values = append(values, value)
			offset += 4
		case SQL_TYPE_VARCHAR:
			length := int(binary.LittleEndian.Uint16(serialized[offset:]))
			offset += 2
			value := string(serialized[offset : offset+length])
			values = append(values, value)
			offset += length
		case SQL_TYPE_BOOL:
			value := serialized[offset] != 0
			values = append(values, value)
			offset++
		case SQL_TYPE_FLOAT:
			value := math.Float32frombits(binary.LittleEndian.Uint32(serialized[offset:]))
			values = append(values, value)
			offset += 4
		default:
			return nil, fmt.Errorf("unsupported data type: %s", dataTypeStr)
		}
	}
	return values, nil
}

// PrintTableMetadata prints the table metadata in a beautiful format
func (t *Table) PrintTableMetadata() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Column Name\tType\tSize\tPrimary Key")
	fmt.Fprintln(w, "------------\t----\t----\t-----------")

	for _, col := range t.Columns {
		primaryKey := "No"
		if col.IsPrimaryKey {
			primaryKey = "Yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			strings.TrimRight(string(col.Name[:]), "\x00"),
			col.DataType.GetDataTypeString(),
			col.Size,
			primaryKey)
	}
	w.Flush()
}

// PrintTable prints the table's contents (actual table, i.e. the column names, and then rows below them as opposed to the metadata)
func (t *Table) PrintTable() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)

	// Determine column widths
	columnWidths := make([]int, len(t.Columns))
	for i, col := range t.Columns {
		columnWidths[i] = len(strings.TrimRight(string(col.Name[:]), "\x00"))
	}
	for _, row := range t.Rows {
		values, err := deserializeValues(row.Values, t.Columns)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i, value := range values {
			width := len(fmt.Sprintf("%v", value))
			if width > columnWidths[i] {
				columnWidths[i] = width
			}
		}
	}

	// Print column names
	for i, col := range t.Columns {
		fmt.Fprintf(w, "%s%s%s\t", ansi.BoldText, ansi.Cyan, padRight(strings.TrimRight(string(col.Name[:]), "\x00"), columnWidths[i]))
	}
	fmt.Fprintln(w, ansi.Reset)

	// Print separator
	for i, width := range columnWidths {
		fmt.Fprint(w, strings.Repeat("-", width))
		if i < len(columnWidths)-1 {
			fmt.Fprint(w, "+")
		}
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range t.Rows {
		values, err := deserializeValues(row.Values, t.Columns)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for i, value := range values {
			fmt.Fprintf(w, "%s%s\t", ansi.Green, padRight(fmt.Sprintf("%v", value), columnWidths[i]))
		}
		fmt.Fprintln(w, ansi.Reset)
	}

	w.Flush()
}

// padRight pads the string with spaces to the specified width
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
