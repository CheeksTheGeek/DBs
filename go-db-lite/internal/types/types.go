package types

// // InputBuffer represents the structure for an input buffer
type InputBuffer struct {
	Buffer       []byte
	BufferLength int
	InputLength  int
}

// // Structure for a table column
// type TableColumn struct {
// 	Name         string
// 	Type         DataType
// 	Size         int
// 	IsPrimaryKey bool
// }

// // Structure for a table
// type Table struct {
// 	Name        string
// 	Columns     []TableColumn
// 	ColumnCount int
// }

// create a type for config
type Config struct {
	HomeDir  string
	InMemory bool
	DbName   string
	DbFile   string
}
