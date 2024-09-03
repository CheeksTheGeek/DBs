import os
import enum

# ANSI escape codes for terminal formatting
ANSI_RESET = "\033[0m"
ANSI_REG_TEXT = "\033[0;3"
ANSI_BOLD_TEXT = "\033[1;3"
ANSI_UNDERLINE_TEXT = "\033[4;3"
ANSI_REG_BG = "\033[4"
ANSI_HIGH_INTENSITY_BG = "\033[0;10"
ANSI_HIGH_INTENSITY_TEXT = "\033[0;9"
ANSI_BOLD_HIGH_INTENSITY_TEXT = "\033[1;9"
ANSI_BLACK = "0m"
ANSI_RED = "1m"
ANSI_GREEN = "2m"
ANSI_YELLOW = "3m"
ANSI_BLUE = "4m"
ANSI_MAGENTA = "5m"
ANSI_CYAN = "6m"
ANSI_WHITE = "7m"


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


class Database:
    """Structure for the database (holds tables)

    Members:
        name (str): The name of the database.
        tables (Dict[str, Table]): The dictionary of tables in the database.
    """

    def __init__(self, name):
        self.name = name
        self.tables = {}

    def add_table(self, table):
        """Add a new table to the database.

        Args:
            table (Table): The table object to add.
        """
        if table.name in self.tables:
            print(f"Table '{table.name}' already exists.")
        else:
            self.tables[table.name] = table
            print(f"Table '{table.name}' created successfully.")

    def get_table(self, table_name):
        """Get a table from the database.

        Args:
            table_name (str): The name of the table to retrieve.

        Returns:
            Table: The table object if found, None otherwise.
        """
        return self.tables.get(table_name, None)


class Table:
    """Structure for a table (holds rows)

    Members:
        name (str): The name of the table.
        columns (List[str]): The list of column names.
        rows (List[List[str]]): The list of rows in the table.
    """

    def __init__(self, name, columns):
        self.name = name
        self.columns = columns
        self.rows = []

    def insert_row(self, row_data):
        """Insert a new row into the table.

        Args:
            row_data (List[str]): The data for the new row.
        """
        if len(row_data) != len(self.columns):
            print("Error: Column count does not match.")
            return
        self.rows.append(row_data)
        print(f"Inserted row into '{self.name}': {row_data}")

    def select_all(self):
        """Select and print all rows in the table."""
        print(f"Table: {self.name}")
        for row in self.rows:
            print(row)


class CommandMapping:
    """Command Mapping Structure

    Members:
        command_str (str): The command string.
        command (enum.Enum): The command type.
        level (enum.Enum): The command level.
    """

    def __init__(self, command_str, command, level):
        self.command_str = command_str
        self.command = command
        self.level = level


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


db_environment = {"current_db": None}  # The active database


def to_lowercase(s):
    return s.lower()


def parse_command(command, core_map, sub_map=None):
    """Parse a command string and return the corresponding command type and level.

    Args:
        command (str): The command string to parse.
        core_map (List[CommandMapping]): The core command mapping list.
        sub_map (List[CommandMapping], optional): The subcommand mapping list. Defaults to None.

    Returns:
        Tuple[CoreCommandType, CommandLevel]: The command type and level.
    """
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


def create_database(database_name):
    """Create a new database and set it as the active database.

    Args:
        database_name (str): The name of the database to create.
    """
    if (
        db_environment["current_db"]
        and db_environment["current_db"].name == database_name
    ):
        print(f"Database '{database_name}' is already in use.")
    else:
        db_environment["current_db"] = Database(database_name)
        print(f"Database '{database_name}' created and set as active.")


def create_table(table_name, columns_def):
    """Create a new table in the active database.

    Args:
        table_name (str): The name of the table to create.
        columns_def (str): The definition of the table columns.
    """
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
    # switcher is a dictionary mapping command types to functions
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


def parse_and_execute_create_command(full_command):
    """Parse and execute subcommands for 'CREATE'

    #TODO: Add support for more 'CREATE' subcommands.
    #TODO      - create index
    #TODO      - create unique index
    #TODO      - create view
    #TODO      - create or replace view
    Args:
        full_command (str): The full command string.
    """
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


def parse_and_execute_insert_command(full_command):
    """Parse and execute the 'INSERT' command.

    Possible syntaxes:
    - insert into <table> values (...)
    - insert into <table> (...) values (...)
    - insert into <table> (...) select ...
    - insert into <table> select ...
    - insert into <table> (...) default values
    - insert into <table> default values
    Args:
        full_command (str): The full command string.
    """
    parts = full_command.split(" ", 2)
    if len(parts) < 3 or "into" not in parts[1]:
        print("Error: Invalid insert syntax. Use 'insert into <table> values (...)'")
        return

    table_name = parts[2].split(" ")[0]
    row_data = parts[2].split("(", 1)[1].strip(")")
    insert_row(table_name, row_data)


def parse_and_execute_select_command(full_command):
    """Parse and execute the 'SELECT' command.

    Args:
        full_command (str): The full command string.
    """
    parts = full_command.split(" ", 2)
    if len(parts) < 3 or "from" not in parts[1]:
        print("Error: Invalid select syntax. Use 'select * from <table>'")
        return

    table_name = parts[2]
    select_all_from_table(table_name)


def exit_program(exception=None):
    """Function to handle exiting the program
    #TODO: Add cleanup code (for files) here if needed.
    #TODO: Add a confirmation prompt before exiting.
    #TODO: Add a way to save the database state before exiting.
    #TODO: Add a way to handle errors and exceptions gracefully.
    """
    if KeyboardInterrupt:
        print(
            ANSI_REG_TEXT
            + ANSI_CYAN
            + "\nExiting due to keyboard interrupt..."
            + ANSI_RESET
        )
    elif exception:
        print(
            ANSI_HIGH_INTENSITY_TEXT
            + ANSI_RED
            + f"\nExiting due to exception: {exception}"
            + ANSI_RESET
        )
    else:
        print(ANSI_REG_TEXT + ANSI_MAGENTA + "\nExiting..." + ANSI_RESET)
    exit()


def start_repl():
    """Start the REPL(Read-Eval-Print Loop) system to read, evaluate, and print commands.
    #TODO: Add a way to handle multi-line commands.
    #TODO: Add a way to handle command history.
    #TODO: Add a way to handle command completion.
    #TODO: Add a way to handle command editing.
    #TODO: Add a way to handle command shortcuts.
    """
    while True:
        try:
            user_input = input(
                ANSI_REG_TEXT
                + ANSI_CYAN
                + "db "
                + ANSI_RESET
                + ANSI_REG_TEXT
                + ANSI_YELLOW
                + "> "
                + ANSI_RESET
            ).strip()
            if user_input:
                core_command, _ = parse_command(user_input, core_command_map)
                execute_core_command(core_command, user_input)
        except EOFError:
            print("\nEOFError: Exiting...")
        except KeyboardInterrupt:
            confirm = input("\nDo you want to exit? (Y/n): ")
            if confirm.lower() == "n":
                continue
            else:
                # TODO: Possibly do some cleanup before exiting
                exit_program(KeyboardInterrupt)
        except Exception as e:
            print(f"\nAn error occurred: {e}")
            exit_program(e)


if __name__ == "__main__":
    print(ANSI_BOLD_TEXT + ANSI_GREEN + "Mini SQL DB starting...\n" + ANSI_RESET)
    start_repl()

for tld in tlds:
    for chars in itertools.product(valid_chars, repeat=1):
        domain = prefix + ''.join(chars) + tld
        if is_domain_available(domain):
            print(f"Available: {domain}")
