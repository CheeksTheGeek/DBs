import os


# Data types for columns
SQL_TYPE_INT = 0
SQL_TYPE_VARCHAR = 1
SQL_TYPE_FLOAT = 2
SQL_TYPE_BOOL = 3

class InBuf:
    def __init__(self):
        self.buffer = None
        self.buffer_length = 0
        self.input_length = 0

class TableColumn:
    def __init__(self, name, col_type, size, is_primary_key, is_nullable):
        self.name = name
        self.type = col_type
        self.size = size
        self.is_primary_key = is_primary_key
        self.is_nullable = is_nullable

class Table:
    def __init__(self, name, columns):
        self.name = name
        self.columns = columns
        self.column_count = len(columns)

def create_table_column(name, col_type, size, is_primary_key, is_nullable):
    return TableColumn(name, col_type, size, is_primary_key, is_nullable)

db_filename = ""

def create_table(table_name, columns):
    global db_filename
    table = Table(table_name, columns)

    print(f"Table '{table_name}' created with {len(columns)} columns.")

    db_filename = f"{table_name}.db"

    if os.path.exists(db_filename):
        with open(db_filename, "a") as db_file:
            db_file.write(f"Table: {table_name}\n")
            for column in columns:
                db_file.write(f"Column: {column.name}, Type: {column.type}, Size: {column.size}, Primary Key: {column.is_primary_key}, Nullable: {column.is_nullable}\n")
    else:
        print("Error opening database file to store table metadata")

def create_database(database_name):
    global db_filename
    db_filename = f"{database_name}.db"

    if os.path.exists(db_filename):
        print(f"Database '{database_name}' already exists.")
        return

    with open(db_filename, "w") as db_file:
        print(f"Database '{database_name}' created successfully as '{db_filename}'.")
        header = "SQLite-like DB file\n"
        db_file.write(header)

def parse_columns(columns_string):
    columns = []
    column_defs = columns_string.split(",")
    
    for column_def in column_defs:
        parts = column_def.split()
        column_name = parts[0]
        data_type = parts[1]
        is_primary_key = "PRIMARY KEY" in column_def
        is_nullable = "NOT NULL" not in column_def

        col_type = {
            "INT": SQL_TYPE_INT,
            "VARCHAR": SQL_TYPE_VARCHAR,
            "FLOAT": SQL_TYPE_FLOAT,
            "BOOL": SQL_TYPE_BOOL
        }.get(data_type.upper(), None)

        if col_type is None:
            print(f"Unknown data type: {data_type}")
            return []

        columns.append(create_table_column(column_name, col_type, 255, is_primary_key, is_nullable))

    return columns

def input_buffer_init():
    return InBuf()

def input_buffer_free(in_buf):
    pass

def read_input(in_buf):
    in_buf.buffer = input("> ")
    in_buf.input_length = len(in_buf.buffer)

def to_lowercase(command):
    return command.lower()

def execute_command(command):
    if "create database" in command:
        db_name = command.replace("create database", "").strip()
        create_database(db_name)
    elif "create table" in command:
        parts = command.split("(", 1)
        table_name = parts[0].replace("create table", "").strip()
        columns_def = parts[1].strip(")")
        columns = parse_columns(columns_def)
        if columns:
            create_table(table_name, columns)
        else:
            print(f"Error parsing columns for table '{table_name}'.")
    else:
        print(f"Unknown command: {command}")

def start_repl():
    in_buf = input_buffer_init()

    while True:
        try:
            read_input(in_buf)
            execute_command(in_buf.buffer)
        except KeyboardInterrupt:
            print("Exiting...")
            break