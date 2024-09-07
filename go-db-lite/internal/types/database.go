package types

import (
	"encoding/binary"
	"fmt"
	"os"
)

type FileHeader struct {
	MagicNumber uint32
	Version     uint32
	TableCount  uint32
}

type Database struct {
	Tables     []Table
	Metadata   map[string]string
	FileHeader FileHeader
}

func NewDatabase() *Database {
	return &Database{
		Tables: make([]Table, 0),
	}
}

func (db *Database) WriteToFile(filename string) error {
	// create the file if it does not exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("error creating database file: %v", err)
		}
		defer file.Close()
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating database file: %v", err)
	}
	defer file.Close()

	// Write the number of tables
	if err := binary.Write(file, binary.LittleEndian, uint32(len(db.Tables))); err != nil {
		return fmt.Errorf("error writing number of tables: %v", err)
	}

	// Write each table
	for _, table := range db.Tables {
		if _, err := table.WriteTo(file); err != nil {
			return fmt.Errorf("error writing table: %v", err)
		}
	}

	return nil
}

func (db *Database) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening database file: %v", err)
	}
	defer file.Close()

	// Read the number of tables
	var numTables uint32
	if err := binary.Read(file, binary.LittleEndian, &numTables); err != nil {
		return fmt.Errorf("error reading number of tables: %v", err)
	}

	// Read each table
	db.Tables = make([]Table, numTables)
	for i := range db.Tables {
		if _, err := db.Tables[i].ReadFrom(file); err != nil {
			return fmt.Errorf("error reading table: %v", err)
		}
	}

	return nil
}

// printDatabase prints the database exhaustively, in a beautiful manner
func (db *Database) PrintDatabase() {
	fmt.Println("Database:")
	fmt.Printf("Magic Number: 0x%X\n", db.FileHeader.MagicNumber)
	fmt.Println("Version:", db.FileHeader.Version)
	fmt.Println("\nMetadata:")
	for key, value := range db.Metadata {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println("\nTables:")
	for _, table := range db.Tables {
		fmt.Printf("\nTable: %s\n", table.Name)
		fmt.Printf("Columns: %d\n", table.ColumnCount)

		// Print table metadata
		fmt.Println("Metadata:")
		for key, value := range table.Metadata {
			fmt.Printf("  %s: %s\n", key, value)
		}

		table.PrintTableMetadata()
	}
}

// AddTable adds a table to the database
func (db *Database) AddTable(table Table) {
	db.Tables = append(db.Tables, table)
	db.FileHeader.TableCount++
}
