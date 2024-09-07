package main

import (
	"fmt"

	"github.com/chaitanyasharma/DBs/go-db-lite/internal/ansi"
)

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
	fmt.Println(ansi.RegBg + ansi.Black + ansi.BoldText + ansi.Yellow + "  │   ── " + ansi.BoldText + ansi.White + "drop view <view_name>" + ansi.Reset)
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
	fmt.Println("This is currently a " + ansi.BoldText + ansi.White + "toy" + ansi.Reset + " project, and I'm learning as I go. So, expect bugs and errors. Please report them at " + ansi.RegText + ansi.Cyan + "https://github.com/chaitanyasharma/DBs/issues" + ansi.Reset + " with the tag 'go-db-lite:'")
}
