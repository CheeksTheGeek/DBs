// internal/parser/parser.go
package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chaitanyasharma/DBs/go-db-lite/internal/types"
)

// ParseCommand parses the command from the input buffer
func ParseCommand(inputBuffer *types.InputBuffer) (types.CommandType, error) {
	tokens := strings.Fields(strings.ToLower(string(inputBuffer.Buffer)))

	if len(tokens) == 0 {
		return types.UnknownCommand{}, errors.New("no command provided")
	}

	coreCommand := findCommand(tokens[0], types.CoreCommandMap)
	if coreCommand != nil {
		if subCommandMap := findSubCommandMap(coreCommand); subCommandMap != nil {
			if len(tokens) < 2 {
				return nil, fmt.Errorf("The top level command \"%s\" requires a subcommand. Available options: %v", coreCommand.CommandName(), listSubCommands(subCommandMap))
			}
			subCommand := findCommand(tokens[1], subCommandMap)
			if subCommand != nil {
				return subCommand, nil
			}
			return types.UnknownCommand{}, fmt.Errorf("Unknown subcommand \"%s\" for command \"%s\"", tokens[1], coreCommand.CommandName())
		}
		return coreCommand, nil
	}

	return types.UnknownCommand{}, fmt.Errorf("Unknown command \"%s\"", tokens[0])
}

// findCommand finds the command in the command map
func findCommand(commandStr string, commandMap []types.CommandMapping) types.CommandType {
	for _, mapping := range commandMap {
		if mapping.CommandStr == commandStr {
			return mapping.Command
		}
	}
	return nil
}

// findSubCommandMap finds the subcommand map for the command
func findSubCommandMap(command types.CommandType) []types.CommandMapping {
	switch cmd := command.(type) {
	case types.CoreCommand:
		for _, mapping := range types.MultiLevelCommandMap {
			if mapping.CommandName == cmd.CommandName() {
				return mapping.CommandMap
			}
		}
	}
	return nil
}

// listSubCommands lists the subcommands for the command
func listSubCommands(commandMap []types.CommandMapping) []string {
	var subCommands []string
	for _, mapping := range commandMap {
		subCommands = append(subCommands, mapping.CommandStr)
	}
	return subCommands
}
