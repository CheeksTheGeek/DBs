package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chaitanyasharma/DBs/go-db-lite/internal/ansi"
	"github.com/chaitanyasharma/DBs/go-db-lite/internal/parser"
	"github.com/chaitanyasharma/DBs/go-db-lite/internal/types"
)

var config types.Config

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
	operatingDirFlag := flag.String("dir", ".", "the directory where the database files will be stored")

	flag.Parse()

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

	fmt.Println(ansi.BoldHighIntensityText + ansi.Green + "Mini SQL DB starting...\n" + ansi.Reset)
	fmt.Println(ansi.RegText + ansi.Magenta + "Type 'exit' to quit" + ansi.Reset)
	if *inMemoryFlag {
		fmt.Println("Using in-memory database")
	} else {
		config.HomeDir, _ = os.Getwd()
		if *operatingDirFlag != "." {
			config.HomeDir = *operatingDirFlag
		}
		fmt.Printf("Current directory: %s has been assumed as the operating directory\n", config.HomeDir)
		if _, err := os.Stat(config.HomeDir + "/" + dbFileName); os.IsNotExist(err) {
			fmt.Printf("%s file has been created and the data will persist in that database\n", dbFileName)
		}
	}

	fmt.Println("dbFileName:", dbFileName)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ansi.RegText + ansi.Cyan + "db" + ansi.Reset + ansi.BoldHighIntensityText + ansi.Yellow + " > " + ansi.Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(ansi.BoldText+ansi.Red+"Error reading input:"+ansi.Reset, err)
			continue
		}

		// Remove the newline character from the input
		input = strings.TrimSpace(input)
		// Exit command (compared case insensitive)
		if strings.EqualFold(input, "exit") || strings.EqualFold(input, ".exit") {
			fmt.Println(ansi.BoldHighIntensityText + ansi.Green + "Until Next Time 👋" + ansi.Reset)
			break
		}

		// Prepare the input buffer for the parser
		inputBuffer := &types.InputBuffer{Buffer: []byte(input)} // using our own input buffer so that we can use it in the parser

		command, err := parser.ParseCommand(inputBuffer)
		if err != nil {
			fmt.Println(ansi.BoldText+ansi.Red+"Error parsing command:"+ansi.Reset, err)
			continue
		}

		executeCommand(command, inputBuffer)
	}
}

// Simple execution function
func executeCommand(command types.CommandType, inputBuffer *types.InputBuffer) {
	switch cmd := command.(type) {
	case types.CreateCommand:
		switch cmd.Command() {
		case types.CmdCreateDatabase:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create Database Command:"+ansi.Reset, cmd.CommandName())
			dbName := strings.TrimSpace(string(inputBuffer.Buffer[len("create database"):]))
			// file, err := os.Create(homeDir + "/" + dbName + ".db") instead check if the file already exists
			if _, err := os.Stat(config.HomeDir + "/" + dbName + ".db"); os.IsNotExist(err) {
				file, err := os.Create(config.HomeDir + "/" + dbName + ".db")
				if err != nil {
					fmt.Println(ansi.BoldText+ansi.Red+"Error creating database file:"+ansi.Reset, err)
					return
				}
				defer file.Close()
				fmt.Println(ansi.RegText+ansi.Green+"Database file created successfully:"+ansi.Reset, dbName+".db")
			} else {
				fmt.Println(ansi.BoldText+ansi.Red+"Database file already exists:"+ansi.Reset, dbName+".db")
			}
		case types.CmdCreateTable:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create Table Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdCreateIndex:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create Index Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdCreateUniqueIndex:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create Unique Index Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdCreateView:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create View Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdCreateProcedure:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Create Procedure Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdCreateUnknown:
			fmt.Println(ansi.BoldText+ansi.Red+"Executing Create Unknown Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Create Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.DropCommand:
		switch cmd.Command() {
		case types.CmdDropTable:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Drop Table Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdDropIndex:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Drop Index Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdDropView:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Drop View Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdDropProcedure:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Drop Procedure Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Drop Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.AlterCommand:
		switch cmd.Command() {
		case types.CmdAlterTable:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Alter Table Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdAlterIndex:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Alter Index Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdAlterView:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Alter View Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdAlterProcedure:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Alter Procedure Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Alter Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.GrantCommand:
		switch cmd.Command() {
		case types.CmdGrant:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Grant Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdRevoke:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Revoke Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Grant Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.LockCommand:
		switch cmd.Command() {
		case types.CmdLockDatabase:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Lock Database Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdLockTable:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Lock Table Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdLockIndex:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Lock Index Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdLockView:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Lock View Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdLockProcedure:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Lock Procedure Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unknown Lock Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.CoreCommand:
		switch cmd.Command() {
		case types.CmdExit:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Exit Command:"+ansi.Reset, cmd.CommandName())
			os.Exit(0)
		case types.CmdSelect:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Select Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdInsert:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Insert Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdUpdate:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Update Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdDelete:
			fmt.Println(ansi.RegText+ansi.Green+"Executing Delete Command:"+ansi.Reset, cmd.CommandName())
		case types.CmdShowDatabases:
			// show databases
			// list all the databases in the operating directory
			files, err := os.ReadDir(config.HomeDir)
			if err != nil {
				fmt.Println(ansi.BoldText+ansi.Red+"Error reading database files:"+ansi.Reset, err)
				return
			}
			fmt.Println(ansi.RegText + ansi.Green + "Databases:" + ansi.Reset)
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".db") {
					fmt.Println(ansi.RegText + ansi.Green + "  " + file.Name() + ansi.Reset)
				}
			}
		case types.CmdUse:
			// use <database_name>
			dbName := strings.TrimSpace(string(inputBuffer.Buffer[len("use"):]))
			fmt.Println(ansi.RegText+ansi.Green+"Executing Use Command:"+ansi.Reset, cmd.CommandName(), dbName)
		case types.CmdHelp:
			printHelp()
		case types.CmdUnknownCommand:
			fmt.Println(ansi.BoldText+ansi.Red+"Executing Unknown Command:"+ansi.Reset, cmd.CommandName())
		default:
			fmt.Println(ansi.BoldText+ansi.Red+"Unrecognized Core Command:"+ansi.Reset, cmd.CommandName())
		}
	case types.UnknownCommand:
		fmt.Println(ansi.RegText + ansi.Red + "Unknown Command" + ansi.Reset)
	default:
		fmt.Println(ansi.RegText + ansi.Red + "Unrecognized Command" + ansi.Reset)
	}
}

// printHelp prints out all the commands, as well the version number of the program
func printHelp() {
	fmt.Println(ansi.HighIntensityText + ansi.Blue + "Version: 0.1" + ansi.Reset)

	fmt.Println(ansi.RegBg + ansi.Blue + ansi.BoldText + ansi.Magenta + "Commands:" + ansi.Reset)

	// select
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "select <column_name / *> from <table_name>" + ansi.Reset)

	// insert
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "insert into <table_name> (column_name, ...) values (value, ...)" + ansi.Reset)

	// update
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "update <table_name> set column_name = value where condition" + ansi.Reset)

	// delete
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "delete from <table_name> where condition" + ansi.Reset)

	// create
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "CREATE commands" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "create database <database_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "create table <table_name> (column_name column_type, ...)" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "create index <index_name> on <table_name> (column_name, ...)" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "create unique index <index_name> on <table_name> (column_name, ...)" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "create view <view_name> as select <select_statement>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   └── " + ansi.BoldText + ansi.White + "create procedure <procedure_name> (param_name param_type, ...)" + ansi.Reset)

	// drop
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "DROP commands" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "drop database <database_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "drop table <table_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "drop index <index_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "drop view <view_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   └── " + ansi.BoldText + ansi.White + "drop procedure <procedure_name>" + ansi.Reset)

	// alter
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "ALTER commands" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "alter table <table_name> add column <column_name> <column_type>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "alter table <table_name> drop column <column_name>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   └── " + ansi.BoldText + ansi.White + "alter table <table_name> rename column <old_column_name> <new_column_name>" + ansi.Reset)

	// grant
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  ├── " + ansi.BoldText + ansi.White + "GRANT commands" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ├── " + ansi.BoldText + ansi.White + "grant <privilege> on <object> to <user>" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   └── " + ansi.BoldText + ansi.White + "revoke <privilege> on <object> from <user>" + ansi.Reset)

	// lock
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  └── " + ansi.BoldText + ansi.White + "LOCK commands" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "      ├── " + ansi.BoldText + ansi.White + "lock database" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "      ├── " + ansi.BoldText + ansi.White + "lock table" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "      ├── " + ansi.BoldText + ansi.White + "lock index" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "      ├── " + ansi.BoldText + ansi.White + "lock view" + ansi.Reset)
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "      └── " + ansi.BoldText + ansi.White + "lock procedure" + ansi.Reset)

	fmt.Println(ansi.RegText + ansi.Magenta + "Type 'exit' to quit" + ansi.Reset)

	// fmt.Println("This is currently a toy project, and I'm learning as I go. So, expect bugs and errors. Please report them at " + ansi.RegText + ansi.Cyan + "https://github.com/chaitanyasharma/DBs/issues" + ansi.Reset)
	// bolden the toy word
	fmt.Println("This is currently a " + ansi.BoldText + ansi.White + "toy" + ansi.Reset + " project, and I'm learning as I go. So, expect bugs and errors. Please report them at " + ansi.RegText + ansi.Cyan + "https://github.com/chaitanyasharma/DBs/issues" + ansi.Reset + " with the tag 'go-db-lite:'")
}
