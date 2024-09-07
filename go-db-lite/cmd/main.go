package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/chaitanyasharma/DBs/go-db-lite/internal/ansi"
	"github.com/chaitanyasharma/DBs/go-db-lite/internal/parser"
	"github.com/chaitanyasharma/DBs/go-db-lite/internal/types"
	"golang.org/x/term"
)

var config *types.Config

const InMemoryDBName = "MEMORY"
const DefaultHomeDirName = "data"
const HistoryFileName = ".history"

// create select insert delete update alter

func main() {
	dbFileName := "default.db"
	if len(os.Args) == 2 && !strings.HasPrefix(os.Args[1], "-") {
		dbFileName = os.Args[1]
	}
	// get the database file name from the user, which by default is the first argument, but if that's not a file, assume the user is providing a flag as -file or --file so use the flags package to get the file name
	// first try to get the file name from the arguments
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		// if the file does not exist, try to get the file name from the flags
		flag.StringVar(&dbFileName, "file", dbFileName, "the file name of the database")
	}

	// specify a help flag
	helpFlag := flag.Bool("help", false, "show help")
	flag.BoolVar(helpFlag, "h", false, "show help")
	versionFlag := flag.Bool("version", false, "show version")
	flag.BoolVar(versionFlag, "v", false, "show version")

	inMemoryFlag := flag.Bool("in-memory", false, "use in-memory database")
	operatingDirFlag := flag.String("dir", "./"+DefaultHomeDirName, "the directory where the database files will be stored")

	flag.Parse()

	// confirm that inMemoryFlag and operatingDirFlag are not both set
	if *inMemoryFlag && *operatingDirFlag != DefaultHomeDirName {
		fmt.Println(ansi.BoldText + ansi.Red + "Error: both in-memory and operating directory flags cannot be set" + ansi.Reset)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// if the help flag is set, print the help message
	if *helpFlag {
		fmt.Println("Usage: go-db-lite [flags]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// if the version flag is set, print the version message
	if *versionFlag {
		fmt.Println("go-db-lite version 0.1")
		os.Exit(0)
	}

	// set config values
	config = types.NewConfig(*operatingDirFlag, *inMemoryFlag, dbFileName)
	fmt.Println(ansi.BoldHighIntensityText + ansi.Green + "Mini SQL DB starting...\n" + ansi.Reset)
	fmt.Println(ansi.RegText + ansi.Magenta + "Type 'exit' to quit" + ansi.Reset)
	if *inMemoryFlag {
		config.InMemory = true
		config.DBFileName = InMemoryDBName + ".db"
		fmt.Println("Using in-memory database")
	} else {
		if *operatingDirFlag != "./"+DefaultHomeDirName {
			config.HomeDir = *operatingDirFlag
		} else {
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Println(ansi.BoldText+ansi.Red+"Error getting current directory:"+ansi.Reset, err)
				os.Exit(1)
			}
			config.HomeDir = currentDir + "/" + DefaultHomeDirName
		}
		// if the folder does not exist, create it
		if _, err := os.Stat(config.HomeDir); os.IsNotExist(err) {
			if err := os.MkdirAll(config.HomeDir, 0755); err != nil {
				fmt.Println(ansi.BoldText+ansi.Red+"Error creating directory:"+ansi.Reset, err)
				os.Exit(1)
			}
		}
		fmt.Printf("The directory: %s has been assumed as the operating directory\n", config.HomeDir)
		if _, err := os.Stat(config.GetDBFilePath()); os.IsNotExist(err) {
			newDB := types.NewDatabase()
			if err := newDB.WriteToFile(config.GetDBFilePath()); err != nil {
				fmt.Printf("Error initializing database file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s file has been created and initialized\n", config.DBFileName)
		}
	}

	// Create or open .history file
	historyFile := config.HomeDir + "/" + HistoryFileName
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		_, err = os.Create(historyFile)
		if err != nil {
			fmt.Println(ansi.BoldText+ansi.Red+"Error creating history file:"+ansi.Reset, err)
			os.Exit(1)
		}
	}

	// Read history file into memory
	historyBytes, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error reading history file:"+ansi.Reset, err)
		os.Exit(1)
	}
	history := strings.Split(string(historyBytes), "\n")
	historyIndex := len(history)

	// Open the history file for appending
	historyFileHandle, err := os.OpenFile(historyFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error opening history file:"+ansi.Reset, err)
		os.Exit(1)
	}
	defer historyFileHandle.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset)

		// Enable raw mode
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("\n"+ansi.BoldText+ansi.Red+"Error setting raw mode:"+ansi.Reset, err)
			continue
		}

		var input string
		var cursorPos int

		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println("\n"+ansi.BoldText+ansi.Red+"Error reading input:"+ansi.Reset, err)
				break
			}

			if char == '\033' {
				// This is an escape sequence, likely an arrow key
				if reader.Buffered() >= 2 {
					seq := make([]byte, 2)
					_, err := reader.Read(seq)
					if err != nil {
						fmt.Println("\n"+ansi.BoldText+ansi.Red+"Error reading input:"+ansi.Reset, err)
						break
					}

					if seq[0] == '[' { // This is the start of the escape sequence
						switch seq[1] {
						case 'A': // Up arrow
							if historyIndex > 0 {
								historyIndex--
								input = history[historyIndex]
								fmt.Print("\r\033[K") // Clear the current line
								fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset + input)
								cursorPos = len(input)
							}
							continue
						case 'B': // Down arrow
							if historyIndex < len(history)-1 {
								historyIndex++
								input = history[historyIndex]
							} else {
								historyIndex = len(history)
								input = ""
							}
							fmt.Print("\r\033[K") // Clear the current line
							fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset + input)
							cursorPos = len(input)
							continue
						case 'C': // Right arrow
							if cursorPos < len(input) {
								cursorPos++
								fmt.Print("\033[1C") // Move cursor right
							}
							continue
						case 'D': // Left arrow
							if cursorPos > 0 {
								cursorPos--
								fmt.Print("\033[1D") // Move cursor left
							}
							continue
						}
					}
				}
			} else if char == '\r' || char == '\n' {
				fmt.Print("\r\n") // Move to the next line
				break
			} else if char == 127 || char == '\b' { // Backspace
				if cursorPos > 0 {
					input = input[:cursorPos-1] + input[cursorPos:]
					cursorPos--
					fmt.Print("\r\033[K") // Clear the current line
					fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset + input)
					fmt.Print("\033[" + fmt.Sprint(cursorPos+4) + "G") // Move cursor to the correct position
				}
			} else if char == 3 || char == 26 { // ctrl+c or ctrl+z
				term.Restore(int(os.Stdin.Fd()), oldState)
				Exit(1)
			} else { // regular character
				input = input[:cursorPos] + string(char) + input[cursorPos:]
				cursorPos++
				fmt.Print("\r\033[K") // Clear the current line
				fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset + input)
				fmt.Print("\033[" + fmt.Sprint(cursorPos+4) + "G") // Move cursor to the correct position
			}
		}

		// Disable raw mode
		term.Restore(int(os.Stdin.Fd()), oldState)

		input = strings.TrimSpace(input)

		// Exit command (compared case insensitive)
		if strings.EqualFold(input, "exit") || strings.EqualFold(input, ".exit") {
			Exit(0)
		}

		// Add the command to history file and memory only if it's not empty
		if input != "" {
			if _, err := historyFileHandle.WriteString(input + "\n"); err != nil {
				fmt.Println(ansi.BoldText+ansi.Red+"Error writing to history file:"+ansi.Reset, err)
			}
			history = append(history, input)
			historyIndex = len(history)

			// Prepare the input buffer for the parser
			inputBuffer := &types.InputBuffer{Buffer: []byte(input)}

			command, err := parser.ParseCommand(inputBuffer)
			if err != nil {
				fmt.Println(ansi.BoldText+ansi.Red+"Error parsing command:"+ansi.Reset, err)
				continue
			}

			executeCommand(command, inputBuffer)
		} else {
			fmt.Println(ansi.BoldText + ansi.Red + "Please enter a valid command" + ansi.Reset)
			continue
		}
	}
}

func Exit(exitCode int) {
	fmt.Print("\r\n")
	if exitCode == 0 {
		fmt.Println(ansi.BoldHighIntensityText + ansi.Green + "Until Next Time 👋" + ansi.Reset)
	} else {
		fmt.Println(ansi.BoldHighIntensityText + ansi.Red + "Exiting disgracefully because you pressed ctrl+c or ctrl+z!?>!? 👹" + ansi.Reset)
	}
	fmt.Print("\r\n")
	os.Exit(exitCode)
}

// Simple execution function
func executeCommand(command types.CommandType, inputBuffer *types.InputBuffer) {
	var tlc string
	switch cmd := command.(type) {
	case types.CreateCommand:
		tlc = "Create"
		switch cmd.Command() {
		case types.CmdCreateDatabase:
			dbName := strings.TrimSpace(string(inputBuffer.Buffer[len("create database"):]))
			if _, err := os.Stat(config.HomeDir + "/" + dbName + ".db"); os.IsNotExist(err) && dbName != "" {
				newDB := types.NewDatabase()
				if err := newDB.WriteToFile(config.HomeDir + "/" + dbName + ".db"); err != nil {
					fmt.Println(ansi.BoldText+ansi.Red+"Error creating database file:"+ansi.Reset, err)
					return
				}
				fmt.Println(ansi.RegText+ansi.Green+"Database file created successfully:"+ansi.Reset, dbName+".db")
			} else if dbName == "" {
				fmt.Println(ansi.BoldText + ansi.Red + "Database name cannot be empty" + ansi.Reset)
			} else {
				fmt.Println(ansi.BoldText+ansi.Red+"Database file already exists:"+ansi.Reset, dbName+".db")
			}
		case types.CmdCreateTable:
			tableName, columnsDef := parseTableCommand(inputBuffer.Buffer)
			columns := parseColumns(columnsDef)
			if columns != nil {
				createTable(tableName, columns)
			} else {
				fmt.Println(ansi.BoldText + ansi.Red + "Error parsing columns for table" + ansi.Reset)
			}
			//
		case types.CmdCreateIndex:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdCreateUniqueIndex:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdCreateView:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdCreateProcedure:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdCreateUnknown:
			notImplemented(tlc, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Create Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.DropCommand:
		tlc = "Drop"
		switch cmd.Command() {
		case types.CmdDropTable:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdDropIndex:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdDropView:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdDropProcedure:
			notImplemented(tlc, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Drop Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.AlterCommand:
		tlc = "Alter"
		switch cmd.Command() {
		case types.CmdAlterTable:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdAlterIndex:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdAlterView:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdAlterProcedure:
			notImplemented(tlc, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Alter Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.GrantCommand:
		tlc = "Grant"
		switch cmd.Command() {
		case types.CmdGrant:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdRevoke:
			notImplemented(tlc, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Grant Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.LockCommand:
		tlc = "Lock"
		switch cmd.Command() {
		case types.CmdLockDatabase:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdLockTable:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdLockIndex:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdLockView:
			notImplemented(tlc, cmd.CommandName())
		case types.CmdLockProcedure:
			notImplemented(tlc, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Lock Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.CoreCommand:
		switch cmd.Command() {
		case types.CmdExit:
			Exit(0)
		case types.CmdSelect:
			query := strings.TrimSpace(string(inputBuffer.Buffer))
			if strings.HasPrefix(strings.ToLower(query), "select * from ") {
				tableName := strings.TrimSpace(query[len("select * from "):])
				if tableName == "" {
					fmt.Println(ansi.BoldText + ansi.Red + "Table name cannot be empty" + ansi.Reset)
				} else {
					// executeSelectAllCommand(tableName)
					db := types.NewDatabase()
					if err := db.ReadFromFile(config.GetDBFilePath()); err != nil {
						fmt.Println(ansi.BoldText+ansi.Red+"Error reading database file during select:"+ansi.Reset, err)
						return
					}
					for _, table := range db.Tables {
						if strings.EqualFold(strings.TrimRight(string(table.Name[:]), "\x00"), tableName) {
							table.PrintTable()
						}
					}
				}
			} else {
				fmt.Println(ansi.BoldText + ansi.Red + "Invalid select command. Only 'SELECT * FROM <table_name>' is supported." + ansi.Reset)
			}
		case types.CmdInsert:
			executeInsertCommand(inputBuffer.Buffer)
		case types.CmdUpdate:
			notImplemented("", cmd.CommandName())
		case types.CmdDelete:
			notImplemented("", cmd.CommandName())
		case types.CmdShowDatabases:
			// list all the databases in the operating directory
			printDatabases(config.HomeDir, config.DBFileName)
		case types.CmdUse:
			// see if it matches the format "use <database_name>"
			dbFileName := strings.TrimSpace(string(inputBuffer.Buffer[len("use"):])) + ".db"
			if len(dbFileName) == 0 {
				fmt.Println(ansi.BoldText + ansi.Red + "Database name cannot be empty" + ansi.Reset)
				return
			}
			dbPath := config.HomeDir + "/" + dbFileName
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				fmt.Println(ansi.BoldText+ansi.Red+"Database file does not exist:"+ansi.Reset, dbPath)
				fmt.Println(ansi.RegText + ansi.Green + "Available databases:" + ansi.Reset)
				printDatabases(config.HomeDir, config.DBFileName)
				return
			}
			config.DBFileName = dbFileName
			fmt.Println(ansi.RegText+ansi.Green+"Using database:"+ansi.Reset, dbFileName)
		case types.CmdHelp:
			printHelp()
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unrecognized Core Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.UnknownCommand:
		fmt.Println(ansi.RegText + ansi.Red + "Unknown Command" + ansi.Reset)
	default:
		fmt.Println(ansi.RegText + ansi.Red + "Unrecognized Command" + ansi.Reset)
	}
}

// parseTableCommand parses the table name and columns definition from the input buffer
func parseTableCommand(buffer []byte) (string, string) {
	command := string(buffer)
	parts := strings.SplitN(command, "(", 2)
	if len(parts) < 2 {
		return "", ""
	}
	tableName := strings.TrimSpace(parts[0][len("create table"):])
	columnsDef := strings.TrimSpace(parts[1])
	return tableName, columnsDef
}

// parseColumns parses the columns definition and returns a slice of Column
func parseColumns(columnsDef string) []types.Column {
	columns := []types.Column{}
	columnDefs := strings.Split(columnsDef, ",")
	for _, columnDef := range columnDefs {
		parts := strings.Fields(strings.TrimSpace(columnDef))
		if len(parts) < 2 {
			return nil
		}
		columnName := parts[0]
		dataType := types.GetDataTypeFromString(parts[1])
		nullable := true
		if len(parts) > 2 && strings.ToLower(parts[2]) == "not" && strings.ToLower(parts[3]) == "null" {
			nullable = false
		}
		var nameBytes [64]byte
		copy(nameBytes[:], columnName)
		columns = append(columns, types.Column{Name: nameBytes, DataType: dataType, Nullable: nullable})
	}
	return columns
}

// createTable creates a new table and adds it to the database
func createTable(tableName string, columns []types.Column) {
	var nameBytes [64]byte
	copy(nameBytes[:], strings.ToLower(tableName))
	for i := range columns {
		copy(columns[i].Name[:], strings.ToLower(string(columns[i].Name[:])))
	}
	table := types.Table{Name: nameBytes, Columns: columns}

	// File-based database logic
	dbFileName := config.GetDBFilePath()
	fmt.Println("dbFileName:", dbFileName)
	db := types.NewDatabase()
	if err := db.ReadFromFile(dbFileName); err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error reading database file during create table:"+ansi.Reset, err)
		return
	}
	// Check if the table already exists
	for _, t := range db.Tables {
		if strings.EqualFold(strings.TrimRight(string(t.Name[:]), "\x00"), tableName) {
			fmt.Println(ansi.BoldText+ansi.Red+"Table already exists:"+ansi.Reset, tableName)
			return
		}
	}

	// Add the new table to the database
	db.AddTable(table)
	// Write the updated database back to the file
	if err := db.WriteToFile(dbFileName); err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error writing to database file:"+ansi.Reset, err)
		return
	}

	fmt.Println(ansi.RegText+ansi.Green+"Table created and added to database file:"+ansi.Reset, tableName)
	db.PrintDatabase()
}

// printDatabases prints out all the databases in the operating directory
func printDatabases(homeDir string, currentDB string) {
	files, err := os.ReadDir(homeDir)
	if err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error reading database files during print:"+ansi.Reset, err)
		return
	}
	fmt.Println(ansi.RegText + ansi.Green + "Databases:" + ansi.Reset)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".db") {
			if file.Name() == currentDB {
				fmt.Println(ansi.RegText + ansi.Green + " " + ansi.Reset + ansi.HighIntensityText + ansi.Magenta + "*" + ansi.Reset + " " + ansi.RegText + ansi.Green + file.Name() + ansi.Reset)
			} else {
				fmt.Println(ansi.RegText + ansi.Green + "   " + file.Name() + ansi.Reset)
			}
		}
	}
}

func notImplemented(prefix string, commandName string) {
	commandName = prefix + " " + commandName
	fmt.Println(ansi.BoldText + ansi.Yellow + "Executing the " + commandName + " functionality ..." + ansi.Reset)
	fmt.Println(ansi.BoldText + ansi.Red + "Not Implemented Yet!" + ansi.Reset)
	fmt.Println(ansi.BoldText + ansi.Red + "Come back another day or contribute to the project at " + ansi.RegText + ansi.Cyan + "https://github.com/chaitanyasharma/DBs" + ansi.Reset + ansi.BoldText + ansi.Red + "!" + ansi.Reset)
	fmt.Println(ansi.BoldText + ansi.Red + "Thank you!" + ansi.Reset)
}

func deserializeValues(data []byte, columns []types.Column) []interface{} {
	values := make([]interface{}, len(columns))
	offset := 0

	for i, col := range columns {
		switch col.DataType {
		case types.SQL_TYPE_INT:
			values[i] = int32(binary.LittleEndian.Uint32(data[offset:]))
			offset += 4
		case types.SQL_TYPE_VARCHAR:
			length := int(binary.LittleEndian.Uint16(data[offset:]))
			offset += 2
			values[i] = string(data[offset : offset+length])
			offset += length
		case types.SQL_TYPE_BOOL:
			values[i] = data[offset] != 0
			offset += 1
		case types.SQL_TYPE_FLOAT:
			values[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
			offset += 4
		default:
			values[i] = "UNKNOWN"
		}
	}

	return values
}

func executeInsertCommand(buffer []byte) {
	query := string(buffer)
	// parts := strings.SplitN(query, "values", 2) // we have to make it case insensitive
	parts := strings.SplitN(strings.ToLower(query), "values", 2)
	fmt.Println("query:", query)
	fmt.Println("parts:", parts)
	if len(parts) != 2 {
		fmt.Println(ansi.BoldText + ansi.Red + "Invalid insert command. Use: INSERT INTO table_name (column1, column2, ...) VALUES (value1, value2, ...)" + ansi.Reset)
		return
	}

	// Parse table name and columns
	tableInfo := strings.TrimSpace(parts[0][len("insert into"):])
	tableNameEnd := strings.Index(tableInfo, "(")
	if tableNameEnd == -1 {
		fmt.Println(ansi.BoldText + ansi.Red + "Invalid insert command. Table name or columns not specified." + ansi.Reset)
		return
	}
	tableName := strings.TrimSpace(tableInfo[:tableNameEnd])
	columns := strings.TrimSpace(tableInfo[tableNameEnd+1 : len(tableInfo)-1])
	columnNames := strings.Split(columns, ",")
	for i, col := range columnNames {
		columnNames[i] = strings.TrimSpace(col)
	}

	// Parse values
	values := strings.TrimSpace(parts[1])
	values = strings.TrimPrefix(values, "(")
	values = strings.TrimSuffix(values, ")")
	valueStrings := strings.Split(values, ",")
	for i, val := range valueStrings {
		valueStrings[i] = strings.TrimSpace(val)
	}

	if len(columnNames) != len(valueStrings) {
		fmt.Println(ansi.BoldText + ansi.Red + "Number of columns does not match number of values." + ansi.Reset)
		return
	}

	// Read the database
	fmt.Println("config.GetDBFilePath():", config.GetDBFilePath())
	db := types.NewDatabase()
	if err := db.ReadFromFile(config.GetDBFilePath()); err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error reading database file during insert:"+ansi.Reset, err)
		return
	}
	fmt.Println("db:", db)
	// Find the table
	var table *types.Table
	for i := range db.Tables {
		fmt.Println("db.Tables[i].Name[:]:", string(db.Tables[i].Name[:]))
		fmt.Println("tableName:", tableName)
		if strings.EqualFold(strings.TrimRight(string(db.Tables[i].Name[:]), "\x00"), tableName) {
			table = &db.Tables[i]
			break
		}
	}

	if table == nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Table not found:"+ansi.Reset, tableName)
		return
	}

	// Validate column names
	tableColumns := table.GetColumnNames()
	fmt.Println("tableColumns:", tableColumns)
	for _, col := range columnNames {
		fmt.Println("col:", col)
		if !containsCaseInsensitive(tableColumns, strings.TrimSpace(col)) {
			fmt.Printf(ansi.BoldText+ansi.Red+"Column %s not found in table %s\n"+ansi.Reset, col, tableName)
			return
		}
	}

	// Prepare the row data
	rowData := make([]interface{}, len(table.Columns))
	for i, col := range table.Columns {
		colName := strings.TrimSpace(string(col.Name[:]))
		valueIndex := indexOfCaseInsensitive(columnNames, colName)
		if valueIndex == -1 {
			if !col.Nullable {
				fmt.Printf(ansi.BoldText+ansi.Red+"Non-nullable column %s is missing a value\n"+ansi.Reset, colName)
				return
			}
			rowData[i] = nil
		} else {
			value := valueStrings[valueIndex]
			convertedValue, err := convertValue(value, col.DataType)
			if err != nil {
				fmt.Printf(ansi.BoldText+ansi.Red+"Error converting value for column %s: %v\n"+ansi.Reset, colName, err)
				return
			}
			rowData[i] = convertedValue
		}
	}

	// Add the row to the table
	if err := table.AddRow(rowData); err != nil {
		fmt.Printf(ansi.BoldText+ansi.Red+"Error adding row: %v\n"+ansi.Reset, err)
		return
	}

	// Write the updated database back to the file
	if err := db.WriteToFile(config.GetDBFilePath()); err != nil {
		fmt.Println(ansi.BoldText+ansi.Red+"Error writing to database file:"+ansi.Reset, err)
		return
	}

	fmt.Println(ansi.RegText + ansi.Green + "Row inserted successfully." + ansi.Reset)
}

func containsCaseInsensitive(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func indexOfCaseInsensitive(slice []string, item string) int {
	for i, s := range slice {
		if strings.EqualFold(s, item) {
			return i
		}
	}
	return -1
}

func convertValue(value string, dataType types.DataType) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	switch dataType {
	case types.SQL_TYPE_INT, types.SQL_TYPE_DOUBLE, types.SQL_TYPE_DECIMAL, types.SQL_TYPE_DATE, types.SQL_TYPE_DATETIME, types.SQL_TYPE_TIMESTAMP, types.SQL_TYPE_BLOB:
		return strconv.Atoi(value)
	case types.SQL_TYPE_VARCHAR, types.SQL_TYPE_CHAR, types.SQL_TYPE_TINYTEXT, types.SQL_TYPE_TEXT, types.SQL_TYPE_MEDIUMTEXT, types.SQL_TYPE_LONGTEXT:
		return strings.Trim(value, "'\""), nil
	case types.SQL_TYPE_BOOL:
		return strconv.ParseBool(value)
	case types.SQL_TYPE_FLOAT:
		return strconv.ParseFloat(value, 32)
	default:
		return nil, fmt.Errorf("unsupported data type")
	}
}
