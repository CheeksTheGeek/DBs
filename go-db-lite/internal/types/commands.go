package types

// Enum Type Names
type CoreCommandType int
type CreateCommandType int
type DropCommandType int
type AlterCommandType int
type GrantCommandType int
type RevokeCommandType int
type LockCommandType int
type UnknownCommandType int

type CommandLevel int

// Core command enum
const (
	CmdExit    CoreCommandType = iota // single level
	CmdInsert                         // single level
	CmdSelect                         // single level
	CmdUpdate                         // single level
	CmdDelete                         // single level
	CmdCreate                         // multi level
	CmdDrop                           // multi level
	CmdAlter                          // multi level
	CmdGrant                          // multi level
	CmdRevoke                         // multi level
	CmdLock                           // multi level
	CmdHelp                           // single level
	CmdUnknown                        // single level
)

// Create-specific command enum
const (
	CmdCreateDatabase CreateCommandType = iota
	CmdCreateTable
	CmdCreateIndex
	CmdCreateUniqueIndex
	CmdCreateView
	CmdCreateOrReplaceView
	CmdCreateProcedure
	CmdCreateUnknown
)

// Drop-specific command enum
const (
	CmdDropTable DropCommandType = iota
	CmdDropIndex
	CmdDropView
	CmdDropProcedure
	CmdDropUnknown
)

// Alter-specific command enum
const (
	CmdAlterTable AlterCommandType = iota
	CmdAlterIndex
	CmdAlterView
	CmdAlterProcedure
	CmdAlterUnknown
)

// Grant-specific command enum
const (
	CmdGrantSelect GrantCommandType = iota
	CmdGrantInsert
	CmdGrantUpdate
	CmdGrantDelete
	CmdGrantCreate
	CmdGrantDrop
	CmdGrantAlter
	CmdGrantIndex
	CmdGrantView
	CmdGrantProcedure
	CmdGrantUnknown
)

// Revoke-specific command enum
const (
	CmdRevokeSelect RevokeCommandType = iota
	CmdRevokeInsert
	CmdRevokeUpdate
	CmdRevokeDelete
	CmdRevokeCreate
	CmdRevokeDrop
	CmdRevokeAlter
	CmdRevokeIndex
	CmdRevokeView
	CmdRevokeProcedure
	CmdRevokeUnknown
)

// Lock-specific command enum
const (
	CmdLockDatabase LockCommandType = iota
	CmdLockTable
	CmdLockIndex
	CmdLockView
	CmdLockProcedure
	CmdLockUnknown
)

const (
	CmdUnknownCommand UnknownCommandType = iota
)

const (
	CoreCommandLevel CommandLevel = iota
	CreateCommandLevel
	DropCommandLevel
	AlterCommandLevel
	GrantCommandLevel
	RevokeCommandLevel
	LockCommandLevel
	UnknownCommandLevel
)

// CommandType is a common interface for all command types.
type CommandType interface {
	CommandName() string
	Level() CommandLevel
	Command() interface{}
}

type ExitCommand struct{}

func (c ExitCommand) CommandName() string {
	return "exit"
}

func (c ExitCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c ExitCommand) Command() interface{} {
	return c
}

// CoreCommand represents commands in the core set.
type CoreCommand struct {
	command CoreCommandType
}

func (c CoreCommand) CommandName() string {
	return [...]string{
		"exit",
		"insert",
		"select",
		"update",
		"delete",
		"create",
		"drop",
		"alter",
		"grant",
		"revoke",
		"lock",
		"unknown"}[c.command]
}

func (c CoreCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c CoreCommand) Command() interface{} {
	return c.command
}

type SelectCommand struct{}

func (c SelectCommand) CommandName() string {
	return "select"
}

func (c SelectCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c SelectCommand) Command() interface{} {
	return c
}

type InsertCommand struct{}

func (c InsertCommand) CommandName() string {
	return "insert"
}

func (c InsertCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c InsertCommand) Command() interface{} {
	return c
}

type UpdateCommand struct{}

func (c UpdateCommand) CommandName() string {
	return "update"
}

func (c UpdateCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c UpdateCommand) Command() interface{} {
	return c
}

type DeleteCommand struct{}

func (c DeleteCommand) CommandName() string {
	return "delete"
}

func (c DeleteCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c DeleteCommand) Command() interface{} {
	return c
}

// CreateCommand represents commands specific to "create".
type CreateCommand struct {
	command CreateCommandType
}

func (c CreateCommand) CommandName() string {
	return [...]string{
		"database",
		"table",
		"index",
		"unique index",
		"view",
		"or replace view",
		"procedure"}[c.command]
}

func (c CreateCommand) Level() CommandLevel {
	return CreateCommandLevel
}

func (c CreateCommand) Command() interface{} {
	return c.command
}

// DropCommand represents commands specific to "drop".
type DropCommand struct {
	command DropCommandType
}

func (c DropCommand) CommandName() string {
	return [...]string{
		"table",
		"index",
		"view",
		"procedure"}[c.command]
}

func (c DropCommand) Level() CommandLevel {
	return DropCommandLevel
}

func (c DropCommand) Command() interface{} {
	return c.command
}

// AlterCommand represents commands specific to "alter".
type AlterCommand struct {
	command AlterCommandType
}

func (c AlterCommand) CommandName() string {
	return [...]string{
		"table",
		"index",
		"view",
		"procedure"}[c.command]
}

func (c AlterCommand) Level() CommandLevel {
	return AlterCommandLevel
}

func (c AlterCommand) Command() interface{} {
	return c.command
}

// GrantCommand represents commands specific to "grant".
type GrantCommand struct {
	command GrantCommandType
}

func (c GrantCommand) CommandName() string {
	return [...]string{
		"select",
		"insert",
		"update",
		"delete",
		"create",
		"drop",
		"alter",
		"index",
		"view",
		"procedure"}[c.command]
}

func (c GrantCommand) Level() CommandLevel {
	return GrantCommandLevel
}

func (c GrantCommand) Command() interface{} {
	return c.command
}

// RevokeCommand represents commands specific to "revoke".
type RevokeCommand struct {
	command RevokeCommandType
}

func (c RevokeCommand) CommandName() string {
	return [...]string{
		"select",
		"insert",
		"update",
		"delete",
		"create",
		"drop",
		"alter",
		"index",
		"view",
		"procedure"}[c.command]
}

func (c RevokeCommand) Level() CommandLevel {
	return RevokeCommandLevel
}

func (c RevokeCommand) Command() interface{} {
	return c.command
}

// LockCommand represents commands specific to "lock".
type LockCommand struct {
	command LockCommandType
}

func (c LockCommand) CommandName() string {
	return [...]string{
		"database",
		"table",
		"index",
		"view",
		"procedure"}[c.command]
}

func (c LockCommand) Level() CommandLevel {
	return LockCommandLevel
}

func (c LockCommand) Command() interface{} {
	return c.command
}

type HelpCommand struct{}

func (c HelpCommand) CommandName() string {
	return "help"
}

func (c HelpCommand) Level() CommandLevel {
	return CoreCommandLevel
}

func (c HelpCommand) Command() interface{} {
	return c
}

// UnknownCommand represents commands that are unknown.
type UnknownCommand struct {
	command UnknownCommandType
}

func (c UnknownCommand) CommandName() string {
	return "unknown"
}

func (c UnknownCommand) Level() CommandLevel {
	return UnknownCommandLevel
}

func (c UnknownCommand) Command() interface{} {
	return c.command
}

// CommandMapping is a mapping of command names to command types.
type CommandMapping struct {
	CommandStr string
	Command    CommandType
}

var CoreCommandMap = []CommandMapping{
	{"exit", CoreCommand{CmdExit}},
	{"insert", CoreCommand{CmdInsert}},
	{"select", CoreCommand{CmdSelect}},
	{"update", CoreCommand{CmdUpdate}},
	{"delete", CoreCommand{CmdDelete}},
	{"create", CoreCommand{CmdCreate}},
	{"drop", CoreCommand{CmdDrop}},
	{"alter", CoreCommand{CmdAlter}},
	{"grant", CoreCommand{CmdGrant}},
	{"revoke", CoreCommand{CmdRevoke}},
	{"lock", CoreCommand{CmdLock}},
	{"help", CoreCommand{CmdHelp}},
	{"unknown", CoreCommand{CmdUnknown}},
}

var CreateCommandMap = []CommandMapping{
	{"database", CreateCommand{CmdCreateDatabase}},
	{"table", CreateCommand{CmdCreateTable}},
	{"index", CreateCommand{CmdCreateIndex}},
	{"unique index", CreateCommand{CmdCreateUniqueIndex}},
	{"view", CreateCommand{CmdCreateView}},
	{"or replace view", CreateCommand{CmdCreateOrReplaceView}},
	{"procedure", CreateCommand{CmdCreateProcedure}},
	{"unknown", CreateCommand{CmdCreateUnknown}},
}

var DropCommandMap = []CommandMapping{
	{"table", DropCommand{CmdDropTable}},
	{"index", DropCommand{CmdDropIndex}},
	{"view", DropCommand{CmdDropView}},
	{"procedure", DropCommand{CmdDropProcedure}},
	{"unknown", DropCommand{CmdDropUnknown}},
}

var AlterCommandMap = []CommandMapping{
	{"table", AlterCommand{CmdAlterTable}},
	{"index", AlterCommand{CmdAlterIndex}},
	{"view", AlterCommand{CmdAlterView}},
	{"procedure", AlterCommand{CmdAlterProcedure}},
	{"unknown", AlterCommand{CmdAlterUnknown}},
}

var GrantCommandMap = []CommandMapping{
	{"select", GrantCommand{CmdGrantSelect}},
	{"insert", GrantCommand{CmdGrantInsert}},
	{"update", GrantCommand{CmdGrantUpdate}},
	{"delete", GrantCommand{CmdGrantDelete}},
	{"create", GrantCommand{CmdGrantCreate}},
	{"drop", GrantCommand{CmdGrantDrop}},
	{"alter", GrantCommand{CmdGrantAlter}},
	{"index", GrantCommand{CmdGrantIndex}},
	{"view", GrantCommand{CmdGrantView}},
	{"procedure", GrantCommand{CmdGrantProcedure}},
	{"unknown", GrantCommand{CmdGrantUnknown}},
}

var RevokeCommandMap = []CommandMapping{
	{"select", RevokeCommand{CmdRevokeSelect}},
	{"insert", RevokeCommand{CmdRevokeInsert}},
	{"update", RevokeCommand{CmdRevokeUpdate}},
	{"delete", RevokeCommand{CmdRevokeDelete}},
	{"create", RevokeCommand{CmdRevokeCreate}},
	{"drop", RevokeCommand{CmdRevokeDrop}},
	{"alter", RevokeCommand{CmdRevokeAlter}},
	{"index", RevokeCommand{CmdRevokeIndex}},
	{"view", RevokeCommand{CmdRevokeView}},
	{"procedure", RevokeCommand{CmdRevokeProcedure}},
	{"unknown", RevokeCommand{CmdRevokeUnknown}},
}

var LockCommandMap = []CommandMapping{
	{"table", LockCommand{CmdLockTable}},
	{"index", LockCommand{CmdLockIndex}},
	{"view", LockCommand{CmdLockView}},
	{"procedure", LockCommand{CmdLockProcedure}},
	{"unknown", LockCommand{CmdLockUnknown}},
}

var UnknownCommandMap = []CommandMapping{
	{"unknown", UnknownCommand{CmdUnknownCommand}},
}

var MultiLevelCommandMap = []struct {
	CommandName string
	CommandMap  []CommandMapping
}{
	{"create", CreateCommandMap},
	{"drop", DropCommandMap},
	{"alter", AlterCommandMap},
	{"grant", GrantCommandMap},
	{"revoke", RevokeCommandMap},
	{"lock", LockCommandMap},
}
