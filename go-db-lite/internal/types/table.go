package types

import (
	"encoding/binary"
	"io"
)

type Column struct {
	Name     [64]byte // Fixed-size field for name
	DataType DataType
	Nullable bool
}

func (c *Column) WriteTo(w io.Writer) error {
	if _, err := w.Write(c.Name[:]); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, c.DataType); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, c.Nullable)
}

func (c *Column) ReadFrom(r io.Reader) error {
	if _, err := io.ReadFull(r, c.Name[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &c.DataType); err != nil {
		return err
	}
	return binary.Read(r, binary.LittleEndian, &c.Nullable)
}

type Row struct {
	Values []byte // Store serialized values as a byte slice
}

func (r *Row) WriteTo(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, r.Values)
}

func (r *Row) ReadFrom(reader io.Reader) error {
	var length uint32
	if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
		return err
	}
	r.Values = make([]byte, length)
	_, err := io.ReadFull(reader, r.Values)
	return err
}

type Table struct {
	Name    [64]byte // Fixed-size field for name
	Columns []Column
	Rows    []Row
}

func (t *Table) AddColumn(name string, dataType DataType, nullable bool) {
	var colName [64]byte
	copy(colName[:], name)
	t.Columns = append(t.Columns, Column{Name: colName, DataType: dataType, Nullable: nullable})
}

func (t *Table) AddRow(values []interface{}) {
	serializedValues := serializeValues(values)
	t.Rows = append(t.Rows, Row{Values: serializedValues})
}

func (t *Table) GetColumnNames() []string {
	var columnNames []string
	for _, column := range t.Columns {
		columnNames = append(columnNames, string(column.Name[:]))
	}
	return columnNames
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

func (t *Table) WriteTo(w io.Writer) error {
	// Write table name
	if _, err := w.Write(t.Name[:]); err != nil {
		return err
	}

	// Write number of columns
	if err := binary.Write(w, binary.LittleEndian, uint32(len(t.Columns))); err != nil {
		return err
	}

	// Write columns
	for _, col := range t.Columns {
		if err := col.WriteTo(w); err != nil {
			return err
		}
	}

	// Write number of rows
	if err := binary.Write(w, binary.LittleEndian, uint32(len(t.Rows))); err != nil {
		return err
	}

	// Write rows
	for _, row := range t.Rows {
		if err := row.WriteTo(w); err != nil {
			return err
		}
	}

	return nil
}

func (t *Table) ReadFrom(r io.Reader) error {
	// Read table name
	if _, err := io.ReadFull(r, t.Name[:]); err != nil {
		return err
	}

	// Read number of columns
	var colCount uint32
	if err := binary.Read(r, binary.LittleEndian, &colCount); err != nil {
		return err
	}

	// Read columns
	t.Columns = make([]Column, colCount)
	for i := range t.Columns {
		if err := t.Columns[i].ReadFrom(r); err != nil {
			return err
		}
	}

	// Read number of rows
	var rowCount uint32
	if err := binary.Read(r, binary.LittleEndian, &rowCount); err != nil {
		return err
	}

	// Read rows
	t.Rows = make([]Row, rowCount)
	for i := range t.Rows {
		if err := t.Rows[i].ReadFrom(r); err != nil {
			return err
		}
	}

	return nil
}

func serializeValues(values []interface{}) []byte {
	// Implement serialization logic here
	// This is a placeholder and needs to be implemented based on your specific requirements
	return []byte{}
}
