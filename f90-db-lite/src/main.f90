module db_constants
    use iso_c_binding, only: c_int, c_size_t, c_char
    implicit none

    ! ANSI escape codes
    character(len=6), parameter :: ANSI_RESET = char(27) // '[0m'
    character(len=6), parameter :: ANSI_REG_TEXT = char(27) // '[0;3'
    character(len=6), parameter :: ANSI_BOLD_TEXT = char(27) // '[1;3'
    character(len=6), parameter :: ANSI_UNDERLINE_TEXT = char(27) // '[4;3'
    character(len=6), parameter :: ANSI_REG_BG = char(27) // '[4'
    character(len=6), parameter :: ANSI_HIGH_INTENSITY_BG = char(27) // '[0;10'
    character(len=6), parameter :: ANSI_HIGH_INTENSITY_TEXT = char(27) // '[0;9'
    character(len=6), parameter :: ANSI_BOLD_HIGH_INTENSITY_TEXT = char(27) // '[1;9'
    character(len=3), parameter :: ANSI_BLACK = '0m'
    character(len=3), parameter :: ANSI_RED = '1m'
    character(len=3), parameter :: ANSI_GREEN = '2m'
    character(len=3), parameter :: ANSI_YELLOW = '3m'
    character(len=3), parameter :: ANSI_BLUE = '4m'
    character(len=3), parameter :: ANSI_MAGENTA = '5m'
    character(len=3), parameter :: ANSI_CYAN = '6m'
    character(len=3), parameter :: ANSI_WHITE = '7m'

    ! Core commands (enum)
    integer(c_int), parameter :: CMD_EXIT = 1
    integer(c_int), parameter :: CMD_INSERT = 2
    integer(c_int), parameter :: CMD_SELECT = 3
    integer(c_int), parameter :: CMD_UPDATE = 4
    integer(c_int), parameter :: CMD_DELETE = 5
    integer(c_int), parameter :: CMD_CREATE = 6
    integer(c_int), parameter :: CMD_DROP = 7
    integer(c_int), parameter :: CMD_ALTER = 8
    integer(c_int), parameter :: CMD_GRANT = 9
    integer(c_int), parameter :: CMD_REVOKE = 10
    integer(c_int), parameter :: CMD_LOCK = 11
    integer(c_int), parameter :: CMD_UNKNOWN = -1

    ! Create-specific command enum
    integer(c_int), parameter :: CMD_CREATE_DATABASE = 1
    integer(c_int), parameter :: CMD_CREATE_TABLE = 2
    integer(c_int), parameter :: CMD_CREATE_INDEX = 3
    integer(c_int), parameter :: CMD_CREATE_UNIQUE_INDEX = 4
    integer(c_int), parameter :: CMD_CREATE_VIEW = 5
    integer(c_int), parameter :: CMD_CREATE_OR_REPLACE_VIEW = 6
    integer(c_int), parameter :: CMD_CREATE_PROCEDURE = 7

    ! Drop-specific command enum
    integer(c_int), parameter :: CMD_DROP_TABLE = 1
    integer(c_int), parameter :: CMD_DROP_INDEX = 2
    integer(c_int), parameter :: CMD_DROP_VIEW = 3
    integer(c_int), parameter :: CMD_DROP_PROCEDURE = 4

    ! Alter-specific command enum
    integer(c_int), parameter :: CMD_ALTER_TABLE = 1
    integer(c_int), parameter :: CMD_ALTER_INDEX = 2
    integer(c_int), parameter :: CMD_ALTER_VIEW = 3
    integer(c_int), parameter :: CMD_ALTER_PROCEDURE = 4

    ! Grant-specific command enum
    integer(c_int), parameter :: CMD_GRANT_SELECT = 1
    integer(c_int), parameter :: CMD_GRANT_INSERT = 2
    integer(c_int), parameter :: CMD_GRANT_UPDATE = 3
    integer(c_int), parameter :: CMD_GRANT_DELETE = 4
    integer(c_int), parameter :: CMD_GRANT_CREATE = 5
    integer(c_int), parameter :: CMD_GRANT_DROP = 6
    integer(c_int), parameter :: CMD_GRANT_ALTER = 7
    integer(c_int), parameter :: CMD_GRANT_INDEX = 8
    integer(c_int), parameter :: CMD_GRANT_VIEW = 9
    integer(c_int), parameter :: CMD_GRANT_PROCEDURE = 10

    ! Revoke-specific command enum
    integer(c_int), parameter :: CMD_REVOKE_SELECT = 1
    integer(c_int), parameter :: CMD_REVOKE_INSERT = 2
    integer(c_int), parameter :: CMD_REVOKE_UPDATE = 3
    integer(c_int), parameter :: CMD_REVOKE_DELETE = 4
    integer(c_int), parameter :: CMD_REVOKE_CREATE = 5
    integer(c_int), parameter :: CMD_REVOKE_DROP = 6
    integer(c_int), parameter :: CMD_REVOKE_ALTER = 7
    integer(c_int), parameter :: CMD_REVOKE_INDEX = 8
    integer(c_int), parameter :: CMD_REVOKE_VIEW = 9
    integer(c_int), parameter :: CMD_REVOKE_PROCEDURE = 10

    ! Lock-specific command enum
    integer(c_int), parameter :: CMD_LOCK_TABLE = 1
    integer(c_int), parameter :: CMD_LOCK_INDEX = 2
    integer(c_int), parameter :: CMD_LOCK_VIEW = 3
    integer(c_int), parameter :: CMD_LOCK_PROCEDURE = 4

    ! Command level enum
    integer(c_int), parameter :: CORE_COMMAND_LEVEL = 1
    integer(c_int), parameter :: CREATE_COMMAND_LEVEL = 2
    integer(c_int), parameter :: DROP_COMMAND_LEVEL = 3
    integer(c_int), parameter :: ALTER_COMMAND_LEVEL = 4
    integer(c_int), parameter :: GRANT_COMMAND_LEVEL = 5
    integer(c_int), parameter :: REVOKE_COMMAND_LEVEL = 6
    integer(c_int), parameter :: LOCK_COMMAND_LEVEL = 7
    integer(c_int), parameter :: UNKNOWN_COMMAND_LEVEL = -1

end module db_constants

module db_structs
    use iso_c_binding, only: c_int, c_size_t, c_char
    use db_constants
    implicit none
    
    ! Structure for input buffer (In_Buf)
    type In_Buf
        character(len=:), allocatable :: buffer
        integer(c_size_t) :: buffer_length
        integer(c_int) :: input_length
    end type In_Buf

    ! Structure for a table column (TableColumn)
    type TableColumn
        character(len=:), allocatable :: name
        integer(c_int) :: type    ! Replace with your DataType equivalent in Fortran
        integer(c_size_t) :: size
        logical :: is_primary_key
        logical :: is_nullable
    end type TableColumn

    ! Structure for a table (Table)
    type Table
        character(len=:), allocatable :: name
        type(TableColumn), allocatable :: columns(:)
        integer(c_size_t) :: column_count
    end type Table

    ! Command Union (simulated with select type)
    type CommandUnion
        class(*), allocatable :: command
    end type CommandUnion

    ! Command Mapping Structure
    type CommandMapping
        character(len=:), allocatable :: command_str
        type(CommandUnion) :: command
        integer(c_int) :: level
    end type CommandMapping

    ! Command mappings arrays (example)
    type(CommandMapping), allocatable :: core_command_map(:)
    type(CommandMapping), allocatable :: create_command_map(:)
    type(CommandMapping), allocatable :: drop_command_map(:)
    type(CommandMapping), allocatable :: alter_command_map(:)
    type(CommandMapping), allocatable :: grant_command_map(:)
    type(CommandMapping), allocatable :: revoke_command_map(:)
    type(CommandMapping), allocatable :: lock_command_map(:)

    contains
    ! Here you would implement logic to initialize these mappings in Fortran
    ! for core_command_map, create_command_map, etc.
    
end module db_structs



module db_repl
    use iso_c_binding, only: c_int, c_size_t, c_char
    use db_structs
    use db_constants
    implicit none
    contains

    ! --------------------------------------------------------------------
    ! Function to initialize the input buffer
    function input_buffer_init() result(in_buf)
        type(In_Buf), allocatable :: in_buf
        allocate(in_buf)
        in_buf%buffer_length = 0
        in_buf%input_length = 0
    end function input_buffer_init

    ! --------------------------------------------------------------------
    ! Subroutine to free the input buffer
    subroutine input_buffer_free(in_buf)
        type(In_Buf), intent(inout) :: in_buf
        if (allocated(in_buf%buffer)) deallocate(in_buf%buffer)
        deallocate(in_buf)
    end subroutine input_buffer_free

    ! --------------------------------------------------------------------
    ! Subroutine to read input from stdin
    subroutine read_input(in_buf)
        type(In_Buf), intent(inout) :: in_buf
        character(len=255) :: input_line
        write(*,'(A)', advance="no") "\x1b[36mdb \x1b[0m\x1b[33m > \x1b[0m"
        read(*,'(A)', iostat=in_buf%input_length) input_line
        if (in_buf%input_length < 0) then
            print *, "Error reading input"
            stop
        endif
        in_buf%input_length = len_trim(input_line)
        allocate(character(len=in_buf%input_length) :: in_buf%buffer)
        in_buf%buffer = trim(input_line)
    end subroutine read_input

    ! --------------------------------------------------------------------
    ! Subroutine to convert a string to lowercase
    subroutine to_lowercase(str)
        character(len=*), intent(inout) :: str
        integer :: i
        do i = 1, len(str)
            str(i:i) = achar(iachar(str(i:i)) + 32 * (iachar(str(i:i)) >= iachar('A') .and. iachar(str(i:i)) <= iachar('Z')))
        end do
    end subroutine to_lowercase

    ! --------------------------------------------------------------------
    ! Function to parse a command from the input buffer
    function parse_command(command, core_map, sub_map, level) result(cmd_union)
        character(len=*), intent(inout) :: command
        type(CommandMapping), intent(in) :: core_map(:)
        type(CommandMapping), intent(in), optional :: sub_map(:)
        integer(c_int), intent(out) :: level
        type(CommandUnion) :: cmd_union
        character(len=255) :: token
        integer :: i

        cmd_union%command => NULL() ! Initialize the union as NULL

        ! Convert the command to lowercase for comparison
        call to_lowercase(command)
        call parse_token(command, token)

        ! Parse the core command
        do i = 1, size(core_map)
            if (trim(token) == trim(core_map(i)%command_str)) then
                level = core_map(i)%level
                cmd_union%command => core_map(i)%command
                return
            endif
        end do

        ! If no match, set unknown command
        level = UNKNOWN_COMMAND_LEVEL
        cmd_union%core_command = CMD_UNKNOWN
    end function parse_command

    ! --------------------------------------------------------------------
    ! Subroutine to execute the parsed command
    subroutine execute_command(command)
        character(len=*), intent(in) :: command
        type(CommandUnion) :: cmd_union
        integer(c_int) :: level

        ! Call the parse command function
        cmd_union = parse_command(command, core_command_map, create_command_map, level)

        if (level == CORE_COMMAND_LEVEL) then
            select case (cmd_union%core_command)
                case (CMD_EXIT)
                    print *, "Exiting..."
                    stop
                case (CMD_INSERT)
                    print *, "Insert command detected."
                case (CMD_SELECT)
                    print *, "Select command detected."
                case (CMD_UPDATE)
                    print *, "Update command detected."
                case (CMD_DELETE)
                    print *, "Delete command detected."
                case (CMD_CREATE)
                    print *, "Create command detected."
                case default
                    print *, "Unknown command."
            end select
        else
            print *, "Unknown command."
        endif
    end subroutine execute_command

    ! --------------------------------------------------------------------
    ! Subroutine to start the REPL (Read-Eval-Print Loop)
    subroutine start_repl()
        type(In_Buf), allocatable :: in_buf
        allocate(in_buf)
        do
            call read_input(in_buf)
            call execute_command(in_buf%buffer)
        end do
        call input_buffer_free(in_buf)
    end subroutine start_repl
end module db_repl