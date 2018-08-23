package golog

import (
	"runtime"
)

type Color int

const (
	Color_None Color = iota
	Color_Black
	Color_Red
	Color_Green
	Color_Yellow
	Color_Blue
	Color_Purple
	Color_DarkGreen
	Color_White
)

var colorByName = map[string]Color{
	"none":      Color_None,
	"black":     Color_Black,
	"red":       Color_Red,
	"green":     Color_Green,
	"yellow":    Color_Yellow,
	"blue":      Color_Blue,
	"purple":    Color_Purple,
	"darkgreen": Color_DarkGreen,
	"white":     Color_White,
}

var logColorPrefix []string

var logColorSuffix string

func init() {
	if runtime.GOOS == "windows" {
		logColorPrefix = []string{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		}
		logColorSuffix = ""
	} else {
		logColorPrefix = []string{
			"",          // None
			"\x1b[030m", // Black
			"\x1b[031m", // Red
			"\x1b[032m", // Green
			"\x1b[033m", // Yellow
			"\x1b[034m", // Blue
			"\x1b[035m", // Purple
			"\x1b[036m", // Darkgreen
			"\x1b[037m", // White
		}
		logColorSuffix = "\x1b[0m"
	}
}

func ColorFromLevel(l Level) Color {
	switch l {
	case Level_Debug:
		return Color_White
	case Level_Info:
		return Color_Green
	case Level_Warn:
		return Color_Yellow
	case Level_Error:
		return Color_Red
	case Level_Fatal:
		return Color_Purple
	}

	return Color_None
}
