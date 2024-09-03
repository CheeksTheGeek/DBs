package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type FileHeader struct {
	MagicNumber uint32
	Version     uint16
}

type Database struct {
	FileHeader FileHeader
	Tables     []Table
	Metadata   map[string]string
}

func NewDatabase() *Database {
	return &Database{
		FileHeader: FileHeader{
			MagicNumber: 0xDBDBDBDB,
			Version:     1,
		},
		Tables:   []Table{},
		Metadata: make(map[string]string),
	}
}

func (db *Database) AddTable(table Table) {
	db.Tables = append(db.Tables, table)
}

func (db *Database) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	buffer := new(bytes.Buffer)

	// Write file header
	if err := binary.Write(buffer, binary.LittleEndian, db.FileHeader); err != nil {
		return fmt.Errorf("error writing file header: %v", err)
	}

	// Write metadata
	for key, value := range db.Metadata {
		if _, err := buffer.WriteString(fmt.Sprintf("%s:%s\n", key, value)); err != nil {
			return fmt.Errorf("error writing metadata: %v", err)
		}
	}

	// Write tables
	for _, table := range db.Tables {
		if err := binary.Write(buffer, binary.LittleEndian, table); err != nil {
			return fmt.Errorf("error writing table: %v", err)
		}
	}

	if _, err := file.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func (db *Database) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(file); err != nil {
		return fmt.Errorf("error reading from file: %v", err)
	}

	// Read file header
	if err := binary.Read(buffer, binary.LittleEndian, &db.FileHeader); err != nil {
		return fmt.Errorf("error reading file header: %v", err)
	}

	// Read metadata
	db.Metadata = make(map[string]string)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			db.Metadata[parts[0]] = parts[1]
		}
	}

	// Read tables
	db.Tables = []Table{}
	for buffer.Len() > 0 {
		var table Table
		if err := binary.Read(buffer, binary.LittleEndian, &table); err != nil {
			return fmt.Errorf("error reading table: %v", err)
		}
		db.Tables = append(db.Tables, table)
	}

	return nil
}
