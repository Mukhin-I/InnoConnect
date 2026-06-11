package logger

import "log"

// This package introduces unified logging package for InnoConnet project

// For logging info
func Info(s string) {
	log.Println("INFO: " + s)
}

// For logging warnings
func Warn(s string) {
	log.Println("WARN: " + s)
}

// For logging errors
func Error(s string) {
	log.Println("ERROR: " + s)
}