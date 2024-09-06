package types

type Column struct {
	Name     string
	DataType DataType
	Nullable bool
}

type Row struct {
	Values []interface{}
}

type Table struct {
	Name    string
	Columns []Column
	Rows    []Row
}

func (t *Table) AddColumn(name string, dataType DataType, nullable bool) {
	t.Columns = append(t.Columns, Column{Name: name, DataType: dataType, Nullable: nullable})
}

func (t *Table) AddRow(values []interface{}) {
	t.Rows = append(t.Rows, Row{Values: values})
}

func (t *Table) GetColumnNames() []string {
	var columnNames []string
	for _, column := range t.Columns {
		columnNames = append(columnNames, column.Name)
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
	t.Name = tableName
	t.Columns = columns
	t.Rows = []Row{}
}
