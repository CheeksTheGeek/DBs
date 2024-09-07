package types

type DataType uint32

const (
	SQL_TYPE_BIT DataType = iota
	SQL_TYPE_TINYINT
	SQL_TYPE_BOOL
	SQL_TYPE_SMALLINT
	SQL_TYPE_MEDIUMINT
	SQL_TYPE_INT
	SQL_TYPE_BIGINT
	SQL_TYPE_FLOAT
	SQL_TYPE_DOUBLE
	SQL_TYPE_DECIMAL
	SQL_TYPE_DATE
	SQL_TYPE_TIME
	SQL_TYPE_DATETIME
	SQL_TYPE_TIMESTAMP
	SQL_TYPE_YEAR
	SQL_TYPE_VARCHAR
	SQL_TYPE_BINARY
	SQL_TYPE_VARBINARY
	SQL_TYPE_TINYBLOB
	SQL_TYPE_TINYTEXT
	SQL_TYPE_TEXT
	SQL_TYPE_BLOB
	SQL_TYPE_CHAR
	SQL_TYPE_MEDIUMTEXT
	SQL_TYPE_MEDIUMBLOB
	SQL_TYPE_LONGTEXT
	SQL_TYPE_LONGBLOB
	SQL_TYPE_ENUM
	SQL_TYPE_SET
	SQL_TYPE_UNKNOWN
)

// mapping of the data types to their corresponding string representations
var dataTypeToString = map[DataType]string{
	SQL_TYPE_BIT:        "BIT",
	SQL_TYPE_TINYINT:    "TINYINT",
	SQL_TYPE_BOOL:       "BOOL",
	SQL_TYPE_SMALLINT:   "SMALLINT",
	SQL_TYPE_MEDIUMINT:  "MEDIUMINT",
	SQL_TYPE_INT:        "INT",
	SQL_TYPE_BIGINT:     "BIGINT",
	SQL_TYPE_FLOAT:      "FLOAT",
	SQL_TYPE_DOUBLE:     "DOUBLE",
	SQL_TYPE_DECIMAL:    "DECIMAL",
	SQL_TYPE_DATE:       "DATE",
	SQL_TYPE_TIME:       "TIME",
	SQL_TYPE_DATETIME:   "DATETIME",
	SQL_TYPE_TIMESTAMP:  "TIMESTAMP",
	SQL_TYPE_YEAR:       "YEAR",
	SQL_TYPE_VARCHAR:    "VARCHAR",
	SQL_TYPE_BINARY:     "BINARY",
	SQL_TYPE_VARBINARY:  "VARBINARY",
	SQL_TYPE_TINYBLOB:   "TINYBLOB",
	SQL_TYPE_TINYTEXT:   "TINYTEXT",
	SQL_TYPE_TEXT:       "TEXT",
	SQL_TYPE_BLOB:       "BLOB",
	SQL_TYPE_CHAR:       "CHAR",
	SQL_TYPE_MEDIUMTEXT: "MEDIUMTEXT",
	SQL_TYPE_MEDIUMBLOB: "MEDIUMBLOB",
	SQL_TYPE_LONGTEXT:   "LONGTEXT",
	SQL_TYPE_LONGBLOB:   "LONGBLOB",
	SQL_TYPE_ENUM:       "ENUM",
	SQL_TYPE_SET:        "SET",
	SQL_TYPE_UNKNOWN:    "UNKNOWN",
}

// GetDataTypeString returns the string representation of the given data type
func (d DataType) GetDataTypeString() string {
	return dataTypeToString[d]
}

// GetDataTypeFromString returns the data type corresponding to the given string representation
func GetDataTypeFromString(dataTypeStr string) DataType {
	for key, value := range dataTypeToString {
		if value == dataTypeStr {
			return key
		}
	}
	return SQL_TYPE_UNKNOWN
}
