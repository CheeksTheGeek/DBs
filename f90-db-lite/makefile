# Fortran Compiler and flags
FC := gfortran
FFLAGS := -Wall -Wextra -pedantic -std=f2008 -g
LDFLAGS := 
INCLUDES := -I./include

# Directories
SRC_DIR := src
OBJ_DIR := obj
BIN_DIR := bin

# Source, object, and executable names
SRCS := $(wildcard $(SRC_DIR)/*.f90)
OBJS := $(patsubst $(SRC_DIR)/%.f90, $(OBJ_DIR)/%.o, $(SRCS))
TARGET := $(BIN_DIR)/mini_sql_db

# Default target
all: $(TARGET)

# Rule for compiling the target executable
$(TARGET): $(OBJS) | $(BIN_DIR)
	$(FC) $(FFLAGS) -o $@ $(OBJS) $(LDFLAGS)

# Rule for compiling object files
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.f90 | $(OBJ_DIR)
	$(FC) $(FFLAGS) -c $< -o $@

# Create necessary directories
$(OBJ_DIR):
	mkdir -p $(OBJ_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Clean target to remove compiled files and directories
clean:
	rm -rf $(OBJ_DIR) $(BIN_DIR)

# Run the program
run: $(TARGET)
	./$(TARGET)

# Phony targets
.PHONY: all clean run