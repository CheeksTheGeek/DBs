import os
import enum


# Enum Definitions for command types
class CoreCommandType(enum.Enum):
    CMD_EXIT = 1
    CMD_INSERT = 2
    CMD_SELECT = 3
    CMD_UPDATE = 4
    CMD_DELETE = 5
    CMD_CREATE = 6
    CMD_DROP = 7
    CMD_ALTER = 8
    CMD_GRANT = 9
    CMD_REVOKE = 10
    CMD_LOCK = 11
    CMD_UNKNOWN = 12


class CreateCommandType(enum.Enum):
    CMD_CREATE_DATABASE = 1
    CMD_CREATE_TABLE = 2
    CMD_CREATE_INDEX = 3
    CMD_CREATE_UNIQUE_INDEX = 4
    CMD_CREATE_VIEW = 5
    CMD_CREATE_OR_REPLACE_VIEW = 6
    CMD_CREATE_PROCEDURE = 7


class CommandLevel(enum.Enum):
    CORE_COMMAND_LEVEL = 1
    CREATE_COMMAND_LEVEL = 2
    UNKNOWN_COMMAND_LEVEL = 3


# Structure for the database (holds tables)
class Database:
    def __init__(self, name):
        self.name = name
        self.tables = {}

    def add_table(self, table):
        if table.name in self.tables:
            print(f"Table '{table.name}' already exists.")
        else:
            self.tables[table.name] = table
            print(f"Table '{table.name}' created successfully.")

    def get_table(self, table_name):
        return self.tables.get(table_name, None)


# Structure for a table (holds rows)
class Table:
    def __init__(self, name, columns):
        self.name = name
        self.columns = columns
        self.rows = []

    def insert_row(self, row_data):
        if len(row_data) != len(self.columns):
            print("Error: Column count does not match.")
            return
        self.rows.append(row_data)
        print(f"Inserted row into '{self.name}': {row_data}")

    def select_all(self):
        print(f"Table: {self.name}")
        for row in self.rows:
            print(row)


# Command Mapping Structure
class CommandMapping:
    def __init__(self, command_str, command, level):
        self.command_str = command_str
        self.command = command
        self.level = level


# Command mappings
core_command_map = [
    CommandMapping("exit", CoreCommandType.CMD_EXIT, CommandLevel.CORE_COMMAND_LEVEL),
    CommandMapping(
        "insert", CoreCommandType.CMD_INSERT, CommandLevel.CORE_COMMAND_LEVEL
    ),
    CommandMapping(
        "select", CoreCommandType.CMD_SELECT, CommandLevel.CORE_COMMAND_LEVEL
    ),
    CommandMapping(
        "create", CoreCommandType.CMD_CREATE, CommandLevel.CORE_COMMAND_LEVEL
    ),
]


create_command_map = [
    CommandMapping(
        "database",
        CreateCommandType.CMD_CREATE_DATABASE,
        CommandLevel.CREATE_COMMAND_LEVEL,
    ),
    CommandMapping(
        "table", CreateCommandType.CMD_CREATE_TABLE, CommandLevel.CREATE_COMMAND_LEVEL
    ),
]


# Simulating a running database environment
db_environment = {"current_db": None}  # The active database


def to_lowercase(s):
    return s.lower()


# Parsing command
def parse_command(command, core_map, sub_map=None):
    command_lower = to_lowercase(command).split()

    if not command_lower:
        return CoreCommandType.CMD_UNKNOWN, CommandLevel.UNKNOWN_COMMAND_LEVEL

    main_command = command_lower[0]

    for mapping in core_map:
        if main_command == mapping.command_str:
            if len(command_lower) > 1 and sub_map:
                sub_command = command_lower[1]
                for sub_mapping in sub_map:
                    if sub_command == sub_mapping.command_str:
                        return sub_mapping.command, sub_mapping.level
            return mapping.command, mapping.level

    return CoreCommandType.CMD_UNKNOWN, CommandLevel.UNKNOWN_COMMAND_LEVEL


# Database operations
def create_database(database_name):
    if (
        db_environment["current_db"]
        and db_environment["current_db"].name == database_name
    ):
        print(f"Database '{database_name}' is already in use.")
    else:
        db_environment["current_db"] = Database(database_name)
        print(f"Database '{database_name}' created and set as active.")


# Table operations
def create_table(table_name, columns_def):
    if db_environment["current_db"] is None:
        print("No active database. Use 'create database <name>' first.")
        return

    columns = [col.strip() for col in columns_def.split(",")]
    new_table = Table(table_name, columns)
    db_environment["current_db"].add_table(new_table)


def insert_row(table_name, row_data):
    if db_environment["current_db"] is None:
        print("No active database.")
        return

    table = db_environment["current_db"].get_table(table_name)
    if table:
        row_data_list = [data.strip() for data in row_data.split(",")]
        table.insert_row(row_data_list)
    else:
        print(f"Table '{table_name}' not found.")


def select_all_from_table(table_name):
    if db_environment["current_db"] is None:
        print("No active database.")
        return

    table = db_environment["current_db"].get_table(table_name)
    if table:
        table.select_all()
    else:
        print(f"Table '{table_name}' not found.")


# Command execution mapping using a dictionary
def execute_core_command(core_command, full_command):
    switcher = {
        CoreCommandType.CMD_EXIT: exit_program,
        CoreCommandType.CMD_CREATE: lambda: parse_and_execute_create_command(
            full_command
        ),
        CoreCommandType.CMD_INSERT: lambda: parse_and_execute_insert_command(
            full_command
        ),
        CoreCommandType.CMD_SELECT: lambda: parse_and_execute_select_command(
            full_command
        ),
        CoreCommandType.CMD_UNKNOWN: lambda: print("Unknown command."),
    }

    # Get the function from the switcher dictionary
    func = switcher.get(core_command, lambda: print("Invalid core command."))

    # Execute the function
    func()


# Parse and execute subcommands for 'CREATE'
def parse_and_execute_create_command(full_command):
    sub_command_str = full_command[len("create ") :]
    create_command, level = parse_command(sub_command_str, create_command_map)

    if level == CommandLevel.CREATE_COMMAND_LEVEL:
        if create_command == CreateCommandType.CMD_CREATE_DATABASE:
            db_name = sub_command_str.split(" ")[1]
            create_database(db_name)
        elif create_command == CreateCommandType.CMD_CREATE_TABLE:
            parts = sub_command_str.split("(", 1)
            if len(parts) < 2:
                print("Error: Invalid table creation syntax.")
                return
            table_name = parts[0].split(" ")[1]
            columns_def = parts[1].strip(")")
            create_table(table_name, columns_def)
        else:
            print(f"Command '{sub_command_str}' is not supported yet.")
    else:
        print(f"Unknown or invalid subcommand for 'create': {sub_command_str}")


# Parse and execute 'INSERT'
def parse_and_execute_insert_command(full_command):
    parts = full_command.split(" ", 2)
    if len(parts) < 3 or "into" not in parts[1]:
        print("Error: Invalid insert syntax. Use 'insert into <table> values (...)'")
        return

    table_name = parts[2].split(" ")[0]
    row_data = parts[2].split("(", 1)[1].strip(")")
    insert_row(table_name, row_data)


# Parse and execute 'SELECT'
def parse_and_execute_select_command(full_command):
    parts = full_command.split(" ", 2)
    if len(parts) < 3 or "from" not in parts[1]:
        print("Error: Invalid select syntax. Use 'select * from <table>'")
        return

    table_name = parts[2]
    select_all_from_table(table_name)


# Function to handle exiting the program
def exit_program():
    print("Exiting...")
    exit()


# Test functions for demonstration
def test_command_execution(command_str):
    # Parse the command
    core_command, level = parse_command(command_str, core_command_map)

    # Execute the command
    execute_core_command(core_command, command_str)


# Example usage
if __name__ == "__main__":
    test_command_execution("create database test_db")
    test_command_execution("create table users (id, name, age)")
    test_command_execution("insert into users values (1, 'Alice', 30)")
    test_command_execution("insert into users values (2, 'Bob', 25)")
    test_command_execution("select * from users")
    test_command_execution("exit")
