// nflag is an open source alternative to Golang's flag package.
// nflag allows OS-specific flag styling, type option simplification, and allowance of empty values.

package nflag

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Config is the config options of nflag
var Config ConfigOptions

// Flags is a map of flags registered with nflag
var Flags map[string]Flag

// LongestFlagLength is an length of the longest Flag, used to ensure we output help cleanly.
var LongestFlagLength int

// Package Init
func init() {
	Flags = make(map[string]Flag) // Make the Flags map
	SetOSFlagString()             // Default to using appropriate OS flag string
}

// Configure is for configuration of nflag prior to usage.
func Configure(providedConfig ConfigOptions) {
	Config = providedConfig // Set Config to providedConfig

	if !Config.OSSpecificFlag { // If we are overriding OSSpecificFlag
		if Config.FlagString == "" { // If no flag string was provided
			SetOSFlagString() // Set to appropriate OS flag string
		}
	}
}

// PrintFlags will print all the flags that are set and their defaults
func PrintFlags() {
	if Config.ProgramDescription != "" { // If we were provided a description of the program
		fmt.Println(Config.ProgramDescription + "\n")
	}

	fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " " + Config.FlagString + "novalueflag" + " " + Config.FlagString + "valueflag=value")
	fmt.Println("The following options are available:")

	// #region Sort Flags

	var flagNames []string

	for flagName := range Flags { // For each flagName in Flags
		flagNames = append(flagNames, flagName) // Append flagName
	}

	sort.Strings(flagNames) // Sort the flagNames

	// #endregion

	for _, flagName := range flagNames { // For each flagName and trimmedFlag struct in Flags
		flag := Flags[flagName] // Get the flag

		thisFlagNameLength := len(flagName)                        // Get the length of this flagName
		flagLengthDiff := (LongestFlagLength - thisFlagNameLength) // Get the difference in flag name length compared to the longest one

		fmt.Println(Config.FlagString + flagName + strings.Repeat(" ", flagLengthDiff) + flag.Descriptor) // Output the flag, creating enough spacing to along descriptor
		if flag.DefaultValue != nil {                                                                     // If DefaultValue is not nil

			if (flag.DefaultValue != "") && (flag.DefaultValue != false) { // If the default value is not an empty string and not false
				fmt.Println("Default Value: ", flag.DefaultValue)
			}
		}

		if flag.AllowNothing { // If we are allowing to pass only the flag
			fmt.Println("Allows Providing Only Flag")
		}

		fmt.Println("")
	}

	os.Exit(1)
}
