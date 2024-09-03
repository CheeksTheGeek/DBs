# Compiler and flags
CC := gcc
CFLAGS := -Wall -Wextra -pedantic -std=c99 -g
LDFLAGS := 
INCLUDES := -I./include

# Directories
SRC_DIR := src
OBJ_DIR := obj
BIN_DIR := bin
INCLUDE_DIR := include

# Source, object, and executable names
SRCS := $(wildcard $(SRC_DIR)/*.c)
OBJS := $(patsubst $(SRC_DIR)/%.c, $(OBJ_DIR)/%.o, $(SRCS))
DEPS := $(OBJS:.o=.d)
TARGET := $(BIN_DIR)/mini_sql_db

# Default target
all: $(TARGET)

# Rule for compiling the target executable
$(TARGET): $(OBJS) | $(BIN_DIR)
	$(CC) $(CFLAGS) -o $@ $(OBJS) $(LDFLAGS)

# Rule for compiling object files and generating dependency files
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c | $(OBJ_DIR)
	$(CC) $(CFLAGS) $(INCLUDES) -MMD -MP -c $< -o $@

# Create necessary directories
$(OBJ_DIR):
	mkdir -p $(OBJ_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Include dependency files generated during compilation
-include $(DEPS)

# Clean target to remove compiled files and directories
clean:
	rm -rf $(OBJ_DIR) $(BIN_DIR)

# Run the program
run: $(TARGET)
	./$(TARGET)

# Format code with clang-format (optional target)
format:
	@if ls $(SRC_DIR)/*.c 1> /dev/null 2>&1; then clang-format -i $(SRC_DIR)/*.c; fi
	@if ls $(INCLUDE_DIR)/*.h 1> /dev/null 2>&1; then clang-format -i $(INCLUDE_DIR)/*.h; fi

# Phony targets (not actual files)
.PHONY: all clean format run