#pragma once

#include <stddef.h>  // for size_t
#include <sys/types.h>  // for ssize_t
#include "data_types.h"

#define ANSI_RESET "\x1b[0m"
#define ANSI_REG_TEXT "\x1b[0;3"
#define ANSI_BOLD_TEXT "\x1b[1;3"
#define ANSI_UNDERLINE_TEXT "\x1b[4;3"
#define ANSI_REG_BG "\x1b[4"
#define ANSI_HIGH_INTENSITY_BG "\x1b[0;10"
#define ANSI_HIGH_INTENSITY_TEXT "\x1b[0;9"
#define ANSI_BOLD_HIGH_INTENSITY_TEXT "\x1b[1;9"
#define ANSI_BLACK "0m"
#define ANSI_RED "1m"
#define ANSI_GREEN "2m"
#define ANSI_YELLOW "3m"
#define ANSI_BLUE "4m"
#define ANSI_MAGENTA "5m"
#define ANSI_CYAN "6m"
#define ANSI_WHITE "7m"
// Structure for input buffer
typedef struct {
    char* buffer;
    size_t buffer_length;
    ssize_t input_length;
} In_Buf;



// Structure for a table column
typedef struct {
    char* name;              // Name of the column
    DataType type;        // Data type of the column
    size_t size;             // Size of the column (for VARCHAR types)
    int is_primary_key;      // Is this column a primary key?
    int is_nullable;         // Is this column nullable?
} TableColumn;


// Structure for a table
typedef struct {
    char* name;              // Name of the table
    TableColumn* columns;    // Array of columns
    size_t column_count;     // Number of columns in the table
} Table;

// Core command enum
typedef enum {
    CMD_EXIT,
    CMD_INSERT,
    CMD_SELECT,
    CMD_UPDATE,
    CMD_DELETE,
    CMD_CREATE,
    CMD_DROP,
    CMD_ALTER,
    CMD_GRANT,
    CMD_REVOKE,
    CMD_LOCK,
    CMD_UNKNOWN
} CoreCommandType;

// Create-specific command enum
typedef enum {
    CMD_CREATE_DATABASE,
    CMD_CREATE_TABLE,
    CMD_CREATE_INDEX,
    CMD_CREATE_UNIQUE_INDEX,
    CMD_CREATE_VIEW,
    CMD_CREATE_OR_REPLACE_VIEW,
    CMD_CREATE_PROCEDURE
} CreateCommandType;

typedef enum {
    CMD_DROP_TABLE,
    CMD_DROP_INDEX,
    CMD_DROP_VIEW,
    CMD_DROP_PROCEDURE
} DropCommandType;

typedef enum {
    CMD_ALTER_TABLE,
    CMD_ALTER_INDEX,
    CMD_ALTER_VIEW,
    CMD_ALTER_PROCEDURE
} AlterCommandType;

typedef enum {
    CMD_GRANT_SELECT,
    CMD_GRANT_INSERT,
    CMD_GRANT_UPDATE,
    CMD_GRANT_DELETE,
    CMD_GRANT_CREATE,
    CMD_GRANT_DROP,
    CMD_GRANT_ALTER,
    CMD_GRANT_INDEX,
    CMD_GRANT_VIEW,
    CMD_GRANT_PROCEDURE
} GrantCommandType;

typedef enum {
    CMD_REVOKE_SELECT,
    CMD_REVOKE_INSERT,
    CMD_REVOKE_UPDATE,
    CMD_REVOKE_DELETE,
    CMD_REVOKE_CREATE,
    CMD_REVOKE_DROP,
    CMD_REVOKE_ALTER,
    CMD_REVOKE_INDEX,
    CMD_REVOKE_VIEW,
    CMD_REVOKE_PROCEDURE
} RevokeCommandType;

typedef enum {
    CMD_LOCK_TABLE,
    CMD_LOCK_INDEX,
    CMD_LOCK_VIEW,
    CMD_LOCK_PROCEDURE
} LockCommandType;

// Union to hold either type of command
typedef union {
    CoreCommandType core_command;
    CreateCommandType create_command;
    DropCommandType drop_command;
    AlterCommandType alter_command;
    GrantCommandType grant_command;
    RevokeCommandType revoke_command;
    LockCommandType lock_command;
} CommandUnion;

typedef enum { // core and all multi level commands like CREATE, DROP, ALTER, 
    CORE_COMMAND_LEVEL,
    CREATE_COMMAND_LEVEL,
    DROP_COMMAND_LEVEL,
    ALTER_COMMAND_LEVEL,
    GRANT_COMMAND_LEVEL,
    REVOKE_COMMAND_LEVEL,
    LOCK_COMMAND_LEVEL,
    UNKNOWN_COMMAND_LEVEL
} CommandLevel;

// Command mapping structure
typedef struct {
    const char* command_str;
    CommandUnion command;
    CommandLevel level;
} CommandMapping;

// Core commands mapping
CommandMapping core_command_map[] = {
    {"exit", { .core_command = CMD_EXIT }, CORE_COMMAND_LEVEL},
    {"insert", { .core_command = CMD_INSERT }, CORE_COMMAND_LEVEL},
    {"select", { .core_command = CMD_SELECT }, CORE_COMMAND_LEVEL},
    {"update", { .core_command = CMD_UPDATE }, CORE_COMMAND_LEVEL},
    {"delete", { .core_command = CMD_DELETE }, CORE_COMMAND_LEVEL},
    {"create", { .core_command = CMD_CREATE }, CORE_COMMAND_LEVEL},
    {"drop", { .core_command = CMD_DROP }, CORE_COMMAND_LEVEL},
    {"alter", { .core_command = CMD_ALTER }, CORE_COMMAND_LEVEL},
    {"grant", { .core_command = CMD_GRANT }, CORE_COMMAND_LEVEL},
    {"revoke", { .core_command = CMD_REVOKE }, CORE_COMMAND_LEVEL},
    {"lock", { .core_command = CMD_LOCK }, CORE_COMMAND_LEVEL},
    {NULL, { .core_command = CMD_UNKNOWN }, CORE_COMMAND_LEVEL},
};

// Create commands mapping
CommandMapping create_command_map[] = {
    {"database", { .create_command = CMD_CREATE_DATABASE }, CREATE_COMMAND_LEVEL},
    {"table", { .create_command = CMD_CREATE_TABLE }, CREATE_COMMAND_LEVEL},
    {"index", { .create_command = CMD_CREATE_INDEX },  CREATE_COMMAND_LEVEL},
    {"unique index", { .create_command = CMD_CREATE_UNIQUE_INDEX },  CREATE_COMMAND_LEVEL},
    {"view", { .create_command = CMD_CREATE_VIEW },  CREATE_COMMAND_LEVEL},
    {"or replace view", { .create_command = CMD_CREATE_OR_REPLACE_VIEW },  CREATE_COMMAND_LEVEL},
    {"procedure", { .create_command = CMD_CREATE_PROCEDURE },  CREATE_COMMAND_LEVEL},
    {NULL, { .create_command = CMD_CREATE_TABLE },  CREATE_COMMAND_LEVEL}
};

// Drop commands mapping
CommandMapping drop_command_map[] = {
    {"table", { .drop_command = CMD_DROP_TABLE }, DROP_COMMAND_LEVEL},
    {"index", { .drop_command = CMD_DROP_INDEX }, DROP_COMMAND_LEVEL},
    {"view", { .drop_command = CMD_DROP_VIEW }, DROP_COMMAND_LEVEL},
    {"procedure", { .drop_command = CMD_DROP_PROCEDURE }, DROP_COMMAND_LEVEL},
    {NULL, { .drop_command = CMD_DROP_TABLE }, DROP_COMMAND_LEVEL}
};

// Alter commands mapping
CommandMapping alter_command_map[] = {
    {"table", { .alter_command = CMD_ALTER_TABLE }, ALTER_COMMAND_LEVEL},
    {"index", { .alter_command = CMD_ALTER_INDEX }, ALTER_COMMAND_LEVEL},
    {"view", { .alter_command = CMD_ALTER_VIEW }, ALTER_COMMAND_LEVEL},
    {"procedure", { .alter_command = CMD_ALTER_PROCEDURE }, ALTER_COMMAND_LEVEL},
    {NULL, { .alter_command = CMD_ALTER_TABLE }, ALTER_COMMAND_LEVEL}
};

// Grant commands mapping
CommandMapping grant_command_map[] = {
    {"select", { .grant_command = CMD_GRANT_SELECT }, GRANT_COMMAND_LEVEL},
    {"insert", { .grant_command = CMD_GRANT_INSERT }, GRANT_COMMAND_LEVEL},
    {"update", { .grant_command = CMD_GRANT_UPDATE }, GRANT_COMMAND_LEVEL},
    {"delete", { .grant_command = CMD_GRANT_DELETE }, GRANT_COMMAND_LEVEL},
    {"create", { .grant_command = CMD_GRANT_CREATE }, GRANT_COMMAND_LEVEL},
    {"drop", { .grant_command = CMD_GRANT_DROP }, GRANT_COMMAND_LEVEL},
    {"alter", { .grant_command = CMD_GRANT_ALTER }, GRANT_COMMAND_LEVEL},
    {"index", { .grant_command = CMD_GRANT_INDEX }, GRANT_COMMAND_LEVEL},
    {"view", { .grant_command = CMD_GRANT_VIEW }, GRANT_COMMAND_LEVEL},
    {"procedure", { .grant_command = CMD_GRANT_PROCEDURE }, GRANT_COMMAND_LEVEL},
    {NULL, { .grant_command = CMD_GRANT_SELECT }, GRANT_COMMAND_LEVEL}
};

// Revoke commands mapping
CommandMapping revoke_command_map[] = {
    {"select", { .revoke_command = CMD_REVOKE_SELECT }, REVOKE_COMMAND_LEVEL},
    {"insert", { .revoke_command = CMD_REVOKE_INSERT }, REVOKE_COMMAND_LEVEL},
    {"update", { .revoke_command = CMD_REVOKE_UPDATE }, REVOKE_COMMAND_LEVEL},
    {"delete", { .revoke_command = CMD_REVOKE_DELETE }, REVOKE_COMMAND_LEVEL},
    {"create", { .revoke_command = CMD_REVOKE_CREATE }, REVOKE_COMMAND_LEVEL},
    {"drop", { .revoke_command = CMD_REVOKE_DROP }, REVOKE_COMMAND_LEVEL},
    {"alter", { .revoke_command = CMD_REVOKE_ALTER }, REVOKE_COMMAND_LEVEL},
    {"index", { .revoke_command = CMD_REVOKE_INDEX }, REVOKE_COMMAND_LEVEL},
    {"view", { .revoke_command = CMD_REVOKE_VIEW }, REVOKE_COMMAND_LEVEL},
    {"procedure", { .revoke_command = CMD_REVOKE_PROCEDURE }, REVOKE_COMMAND_LEVEL},
    {NULL, { .revoke_command = CMD_REVOKE_SELECT }, REVOKE_COMMAND_LEVEL}
};

// Lock commands mapping
CommandMapping lock_command_map[] = {
    {"table", { .lock_command = CMD_LOCK_TABLE }, LOCK_COMMAND_LEVEL},
    {"index", { .lock_command = CMD_LOCK_INDEX }, LOCK_COMMAND_LEVEL},
    {"view", { .lock_command = CMD_LOCK_VIEW }, LOCK_COMMAND_LEVEL},
    {"procedure", { .lock_command = CMD_LOCK_PROCEDURE }, LOCK_COMMAND_LEVEL},
    {NULL, { .lock_command = CMD_LOCK_TABLE }, LOCK_COMMAND_LEVEL}
};




// Function declarations
In_Buf* input_buffer_init(void);
void input_buffer_free(In_Buf* in_buf);
void read_input(In_Buf* in_buf);
void execute_command(char* command);
void start_repl(In_Buf* in_buf);

