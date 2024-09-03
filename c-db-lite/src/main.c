#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "main.h"
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>

char db_filename[256];

/**
 * @brief Create a new table
 * 
 * @param table_name Name of the table
 * @param columns Array of TableColumn structures
 * @param column_count Number of columns in the table
 */
void create_table(const char* table_name, TableColumn* columns, size_t column_count) {
    Table* table = (Table*) malloc(sizeof(Table));
    table->name = strdup(table_name);
    table->columns = (TableColumn*) malloc(column_count * sizeof(TableColumn));
    table->column_count = column_count;

    for (size_t i = 0; i < column_count; ++i) {
        table->columns[i] = columns[i];
    }

    printf("Table '%s' created with %zu columns.\n", table_name, column_count);

    // Store table metadata in the database file
    snprintf(db_filename, sizeof(db_filename), "%s.db", table_name); // Assuming the database file is named after the table

    int db_file = open(db_filename, O_WRONLY | O_APPEND);
    if (db_file < 0) {
        perror("Error opening database file to store table metadata");
        return;
    }
    // Write table metadata
    dprintf(db_file, "Table: %s\n", table_name);
    for (size_t i = 0; i < column_count; ++i)
        dprintf(db_file, "Column: %s, Type: %d, Size: %zu, Primary Key: %d, Nullable: %d\n", 
                        columns[i].name, columns[i].type, columns[i].size, columns[i].is_primary_key, columns[i].is_nullable);
    close(db_file);
}

TableColumn create_table_column(const char* name, DataType type, size_t size, int is_primary_key, int is_nullable) {
    TableColumn column;
    column.name = strdup(name);
    column.type = type;
    column.size = size;
    column.is_primary_key = is_primary_key;
    column.is_nullable = is_nullable;
    return column;
}


/**
 * @brief Create a database object
 * 
 * @param database_name 
 */
void create_database(const char* database_name) {
    if (strstr(database_name, ".db") != NULL) {
        snprintf(db_filename, sizeof(db_filename), "%s", database_name);
    } else {
        snprintf(db_filename, sizeof(db_filename), "%s.db", database_name); // Append ".db" extension
    }
    // Check if the file already exists
    if (!access(db_filename, F_OK)) { // checks if the file exists and can be accessed by current permissions
        printf("Database '%s' already exists.\n", database_name);
        return;
    }
    // Create the file with read/write permissions
    int db_file = open(db_filename, O_CREAT | O_WRONLY, 0644); // 0644 = rw-r--r--
    if (db_file < 0) {
        perror("Error creating database");
        return;
    }
    printf("Database '%s' created successfully as '%s'.\n", database_name, db_filename);
    // Optionally, write some metadata (e.g., a simple header)
    const char* header = "SQLite-like DB file\n";
    write(db_file, header, strlen(header));
    // Close the file
    close(db_file);
}

/**
 * @brief Helper function to parse column definitions from the input string
 * 
 * @param columns_string The string containing the column definitions
 * @param columns The array to store parsed TableColumn objects
 * @return size_t The number of parsed columns
 */
size_t parse_columns(const char* columns_string, TableColumn* columns) {
    size_t column_count = 0;
    char* column_def = strtok(strdup(columns_string), ",");

    while (column_def != NULL) {
        char data_type[128], column_name[128];
        int is_primary_key = 0, is_nullable = 1;

        // Parse column definition (assuming basic format: column_name data_type [PRIMARY KEY] [NOT NULL])
        sscanf(column_def, "%s %s", column_name, data_type);

        // Check if PRIMARY KEY or NOT NULL are present
        if (strstr(column_def, "PRIMARY KEY") != NULL) {
            is_primary_key = 1;
        }
        if (strstr(column_def, "NOT NULL") != NULL) {
            is_nullable = 0;
        }

        // Map the data type string to SqlDataType enum
        DataType type;
        if (strcasecmp(data_type, "INT") == 0) {
            type = SQL_TYPE_INT;
        } else if (strncasecmp(data_type, "VARCHAR", 7) == 0) {
            type = SQL_TYPE_VARCHAR;
        } else if (strcasecmp(data_type, "FLOAT") == 0) {
            type = SQL_TYPE_FLOAT;
        } else if (strcasecmp(data_type, "BOOL") == 0) {
            type = SQL_TYPE_BOOL;
        } else {
            printf("Unknown data type: %s\n", data_type);
            return 0;
        }

        // Add the parsed column to the columns array
        columns[column_count++] = create_table_column(column_name, type, 255, is_primary_key, is_nullable);

        // Move to the next column definition
        column_def = strtok(NULL, ",");
    }

    return column_count;
}

/**
 * @brief Initializes the input buffer
 * 
 * @return The input buffer
 */
In_Buf* input_buffer_init(void) {
    In_Buf* in_buf = (In_Buf*) malloc(sizeof(In_Buf));
    in_buf->buffer = NULL;
    in_buf->buffer_length = 0;
    in_buf->input_length = 0;
    return in_buf;
}

/**
 * @brief Frees the input buffer
 * 
 * @param in_buf (In_Buf*) The input buffer
 */
void input_buffer_free(In_Buf* in_buf) {
    free(in_buf->buffer);
    free(in_buf);
}

/**
 * @brief Reads input from stdin and stores it in the input buffer
 * 
 * @param in_buf (In_Buf*) The input buffer
 */
void read_input(In_Buf* in_buf) {
    ssize_t bytes_read = getline(&(in_buf->buffer), &(in_buf->buffer_length), stdin);
    if (bytes_read <= 0) {
        printf("Error reading input\n");
        exit(EXIT_FAILURE);
    }
    in_buf->input_length = bytes_read - 1;
    in_buf->buffer[bytes_read - 1] = 0; // Remove trailing newline
}

/**
 * @brief Converts a string to lowercase
 * 
 * @param str (char*) The string to convert
 */
void to_lowercase(char* str) {
    for (; *str; ++str) *str = tolower(*str);
}

/**
 * @brief Parses the command and subcommand from the input buffer
 * 
 * @param command (char*) The full command string to parse
 * @return The appropriate command union and level
 */
CommandUnion parse_command(char* command, CommandMapping core_map[], CommandMapping sub_map[], CommandLevel* level) {

    char command_lower[strlen(command) + 1];
    strcpy(command_lower, command);
    to_lowercase(command_lower);
    char* token = strtok(command_lower, " ");
    
    if (token == NULL) {
        *level = UNKNOWN_COMMAND_LEVEL;
        CommandUnion unknown_command;
        unknown_command.core_command = CMD_UNKNOWN;
        return unknown_command;
    }

    for (CommandMapping* mapping = core_map; mapping->command_str != NULL; ++mapping) {
        if (!strcmp(token, mapping->command_str)) {
            *level = mapping->level;

            token = strtok(NULL, " ");
            if (token != NULL && sub_map != NULL) {

                for (CommandMapping* sub_mapping = sub_map; sub_mapping->command_str != NULL; ++sub_mapping) {
                    if (!strcmp(token, sub_mapping->command_str)) {
                        *level = sub_mapping->level;  // Update to subcommand level
                        return sub_mapping->command;
                    }
                }
                // If no matching subcommand was found
                printf("Unknown subcommand: '%s'\n", token);
            } else if (token == NULL && sub_map != NULL) {
                // Handle the case where no subcommand is provided
                printf("No subcommand provided for the command '%s'.\n", mapping->command_str);
            }
            return mapping->command;  // Return the core command if no subcommand
        }
    }

    *level = UNKNOWN_COMMAND_LEVEL;
    CommandUnion unknown_command;
    unknown_command.core_command = CMD_UNKNOWN;
    printf("Unknown command: '%s'\n", command);
    return unknown_command;
}
/**
 * @brief Executes a command with subcommands from the input buffer
 * 
 * @param command (char*) The command string to execute
 */
void execute_command(char* command) {
    CommandLevel level;
    CommandUnion cmd_union;
    
    char command_copy[strlen(command) + 1];
    strcpy(command_copy, command);
    
    // Parse the core command
    cmd_union = parse_command(command_copy, core_command_map, NULL, &level);

    if (level == CORE_COMMAND_LEVEL) {
        switch (cmd_union.core_command) {
            case CMD_EXIT:
                printf("Exiting...\n");
                exit(EXIT_SUCCESS);
                break;
            case CMD_INSERT:
                printf("Insert command detected which is currently not supported.\n");
                break;
            case CMD_SELECT:
                printf("Select command detected which is currently not supported.\n");
                break;
            case CMD_UPDATE:
                printf("Update command detected which is currently not supported.\n");
                break;
            case CMD_DELETE:
                printf("Delete command detected which is currently not supported.\n");
                break;
            case CMD_CREATE:
                cmd_union = parse_command(command + strlen("create") + 1, create_command_map, NULL, &level);
                if (level == CREATE_COMMAND_LEVEL) {
                    switch (cmd_union.create_command) {
                        case CMD_CREATE_DATABASE:{
                            printf("Creating database...\n");
                            char* db_name = command + strlen("create database") + 1;  // Extract database name
                            create_database(db_name);
                            break;
                        }
                        case CMD_CREATE_TABLE:
                            printf("Creating table...\n");
                            char* table_name = strtok(command + strlen("create table") + 1, "(");
                            char* columns_def = strtok(NULL, ")");
                            if (table_name && columns_def) {
                                TableColumn columns[256]; // Assuming max 10 columns for simplicity
                                size_t column_count = parse_columns(columns_def, columns);
                                if (column_count > 0)
                                    create_table(table_name, columns, column_count);
                                else
                                    printf("Error parsing columns for table '%s'.\n", table_name);
                            } else
                                printf("Error parsing table name or columns.\n");
                            break;
                        case CMD_CREATE_INDEX:
                            printf("Creating index is currently not supported.\n");
                            break;
                        case CMD_CREATE_UNIQUE_INDEX:
                            printf("Creating unique index is currently not supported.\n");
                            break;
                        case CMD_CREATE_VIEW:
                            printf("Creating view is currently not supported.\n");
                            break;
                        case CMD_CREATE_OR_REPLACE_VIEW:
                            printf("Creating or replacing view is currently not supported.\n");
                            break;
                        case CMD_CREATE_PROCEDURE:
                            printf("Creating procedure is currently not supported.\n");
                            break;
                        default:
                            printf("Unknown create subcommand.\n");
                            break;
                    }
                }
                break;
            case CMD_DROP:
                printf("Drop command detected. Parsing subcommand...\n");
                cmd_union = parse_command(command + strlen("drop") + 1, drop_command_map, NULL, &level);
                if (level == DROP_COMMAND_LEVEL) {
                    switch (cmd_union.drop_command) {
                        case CMD_DROP_TABLE:
                            printf("Dropping table currently not supported.\n");
                            break;
                        case CMD_DROP_INDEX:
                            printf("Dropping index currently not supported.\n");
                            break;
                        case CMD_DROP_VIEW:
                            printf("Dropping view currently not supported.\n");
                            break;
                        case CMD_DROP_PROCEDURE:
                            printf("Dropping procedure currently not supported.\n");
                            break;
                        default:
                            printf("Unknown drop subcommand.\n");
                            break;
                    }
                }
                break;
            case CMD_ALTER:
                printf("Alter command detected. Parsing subcommand...\n");
                cmd_union = parse_command(command + strlen("alter") + 1, alter_command_map, NULL, &level);
                if (level == ALTER_COMMAND_LEVEL) {
                    switch (cmd_union.alter_command) {
                        case CMD_ALTER_TABLE:
                            printf("Altering table is currently not supported.\n");
                            break;
                        case CMD_ALTER_INDEX:
                            printf("Altering index is currently not supported.\n");
                            break;
                        case CMD_ALTER_VIEW:
                            printf("Altering view is currently not supported.\n");
                            break;
                        case CMD_ALTER_PROCEDURE:
                            printf("Altering procedure is currently not supported.\n");
                            break;
                        default:
                            printf("Unknown alter subcommand.\n");
                            break;
                    }
                }
                break;
            case CMD_GRANT:
                printf("Grant command detected. Parsing subcommand...\n");
                cmd_union = parse_command(command + strlen("grant") + 1, grant_command_map, NULL, &level);
                if (level == GRANT_COMMAND_LEVEL) {
                    switch (cmd_union.grant_command) {
                        case CMD_GRANT_SELECT:
                            printf("Granting select permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_INSERT:
                            printf("Granting insert permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_UPDATE:
                            printf("Granting update permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_DELETE:
                            printf("Granting delete permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_CREATE:
                            printf("Granting create permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_DROP:
                            printf("Granting drop permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_ALTER:
                            printf("Granting alter permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_INDEX:
                            printf("Granting index permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_VIEW:
                            printf("Granting view permission is currently not supported.\n");
                            break;
                        case CMD_GRANT_PROCEDURE:
                            printf("Granting procedure permission is currently not supported.\n");
                            break;
                        default:
                            printf("Unknown grant subcommand.\n");
                            break;
                    }
                }
                break;
            case CMD_REVOKE:
                printf("Revoke command detected. Parsing subcommand...\n");
                cmd_union = parse_command(command + strlen("revoke") + 1, revoke_command_map, NULL, &level);
                if (level == REVOKE_COMMAND_LEVEL) {
                    switch (cmd_union.revoke_command) {
                        case CMD_REVOKE_SELECT:
                            printf("Revoking select permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_INSERT:
                            printf("Revoking insert permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_UPDATE:
                            printf("Revoking update permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_DELETE:
                            printf("Revoking delete permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_CREATE:
                            printf("Revoking create permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_DROP:
                            printf("Revoking drop permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_ALTER:
                            printf("Revoking alter permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_INDEX:
                            printf("Revoking index permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_VIEW:
                            printf("Revoking view permission is currently not supported.\n");
                            break;
                        case CMD_REVOKE_PROCEDURE:
                            printf("Revoking procedure permission is currently not supported.\n");
                            break;
                        default:
                            printf("Unknown revoke subcommand.\n");
                            break;
                    }
                }
                break;
            case CMD_LOCK:
                printf("Lock command detected. Parsing subcommand...\n");
                cmd_union = parse_command(command + strlen("lock") + 1, lock_command_map, NULL, &level);
                if (level == LOCK_COMMAND_LEVEL) {
                    switch (cmd_union.lock_command) {
                        case CMD_LOCK_TABLE:
                            printf("Locking table currently not supported.\n");
                            break;
                        case CMD_LOCK_INDEX:
                            printf("Locking index currently not supported.\n");
                            break;
                        case CMD_LOCK_VIEW:
                            printf("Locking view currently not supported.\n");
                            break;
                        case CMD_LOCK_PROCEDURE:
                            printf("Locking procedure currently not supported.\n");
                            break;
                        default:
                            printf(ANSI_BOLD_TEXT ANSI_RED "Unknown Lock Subcommand.\n" ANSI_RESET);
                            break;
                    }
                }
                break;
            default:
                printf(ANSI_BOLD_TEXT ANSI_RED "Unrecognized Core Command: '%s'\n" ANSI_RESET, command);
                break;
        }
    } else {
        printf(ANSI_BOLD_TEXT ANSI_RED "Unrecognized Command: '%s'\n" ANSI_RESET, command);
    }
}

/**
 * @brief Starts the REPL (Read-Eval-Print Loop)
 * 
 */
void start_repl(In_Buf* in_buf) {
    while (1) {
        printf("\x1b[36mdb \x1b[0m\x1b[33m > \x1b[0m");
        read_input(in_buf);
        execute_command(in_buf->buffer);
    }
}

/**
 * @brief Main function for a simple SQL database
 * 
 * @param argc The number of command-line arguments
 * @param argv The command-line arguments
 * @return The exit status
 */
int main(int argc, char* argv[]) {
    printf(ANSI_BOLD_HIGH_INTENSITY_TEXT ANSI_GREEN "Mini SQL DB starting...\n" ANSI_RESET);
    In_Buf* in_buf = input_buffer_init();
    start_repl(in_buf);
    input_buffer_free(in_buf);
    return 0;
}
