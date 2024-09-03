package ansi

/*
 * ANSI Color Codes Usage in Print Or Other Functions
 *
 * Printf(AnsiRegText + AnsiGreen + "Hello, World!" + AnsiReset)
 * OR
 * PrintLn(AnsiRegText + AnsiGreen + "Hello, World!" + AnsiReset)
 * OR
 * Print(AnsiRegText + AnsiGreen + "Hello, World!" + AnsiReset)
 * OR
 * Println(AnsiRegText + AnsiGreen + "Hello, World!" + AnsiReset)
 */

const (
	Reset                 = "\x1b[0m"
	RegText               = "\x1b[0;3"
	BoldText              = "\x1b[1;3"
	UnderlineText         = "\x1b[4;3"
	RegBg                 = "\x1b[4"
	HighIntensityBg       = "\x1b[0;10"
	HighIntensityText     = "\x1b[0;9"
	BoldHighIntensityText = "\x1b[1;9"
	Black                 = "0m"
	Red                   = "1m"
	Green                 = "2m"
	Yellow                = "3m"
	Blue                  = "4m"
	Magenta               = "5m"
	Cyan                  = "6m"
	White                 = "7m"
)
