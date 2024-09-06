package types

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
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

	// Write file header
	if err := binary.Write(file, binary.LittleEndian, db.FileHeader); err != nil {
		return fmt.Errorf("error writing file header: %v", err)
	}

	// Write metadata
	for key, value := range db.Metadata {
		if _, err := file.WriteString(fmt.Sprintf("%s:%s\n", key, value)); err != nil {
			return fmt.Errorf("error writing metadata: %v", err)
		}
	}

	// Write tables
	for _, table := range db.Tables {
		if err := table.WriteTo(file); err != nil {
			return fmt.Errorf("error writing table: %v", err)
		}
	}

	return nil
}

func (db *Database) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file header
	if err := binary.Read(file, binary.LittleEndian, &db.FileHeader); err != nil {
		if err == io.EOF {
			// File is empty, initialize with default header
			db.FileHeader = FileHeader{Version: 1, MagicNumber: 0x4744424C} // 'GDBL' in hex
		} else {
			return fmt.Errorf("error reading file header: %v", err)
		}
	}

	// Read metadata
	db.Metadata = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			db.Metadata[parts[0]] = parts[1]
		}
	}

	// Read tables
	db.Tables = []Table{}
	for {
		var table Table
		if err := table.ReadFrom(file); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading table: %v", err)
		}
		db.Tables = append(db.Tables, table)
	}

	return nil
}
