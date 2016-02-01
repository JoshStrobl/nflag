// nflag is an open source alternative to Golang's flag package.
// nflag allows OS-specific flag styling, type option simplification, and allowance of empty values.

package nflag

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var Config ConfigOptions
var Flags map[string]Flag
var LongestFlagLength int

// Package Init
func init() {
	Flags = make(map[string]Flag)      // Make the Flags map
	Config.OSSpecificFlags = true      // Default to being true
	Config.OSSpecificFlagString = "--" // Default to using a double dash for declaring flags

	if runtime.GOOS == "windows" { // If we are on Windows
		Config.OSSpecificFlagString = "/" // Use a slash, since Windows likes to be "special"
	}
}

// Configure
// This function is for configuration of nflag prior to usage.
func Configure(providedConfig ConfigOptions) {
	if !providedConfig.OSSpecificFlags { // If we are overriding OSSpecificFlags
		Config.OSSpecificFlags = false
		Config.OSSpecificFlagString = "--" // Enforce --
	}

	if providedConfig.OSSpecificFlagString != "" { // If we are overriding OSSpecificFlagString
		Config.OSSpecificFlagString = providedConfig.OSSpecificFlagString // Set OSSpecificFlagString to whatever is provided on config
	}

	Config.ProgramDescription = providedConfig.ProgramDescription // Assign providedConfig program description (if any) to Config

	if providedConfig.ShowHelpIfNoArgs { // If ShowHelpIfNoArgs was set to true
		Config.ShowHelpIfNoArgs = true // Change to true
	}
}

// PrintFlags
// This function will print all the flags that are set and their defaults
func PrintFlags() {
	if Config.ProgramDescription != "" { // If we were provided a description of the program
		fmt.Println(Config.ProgramDescription + "\n")
	}

	fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " " + Config.OSSpecificFlagString + "novalueflag" + " " + Config.OSSpecificFlagString + "valueflag=value")
	fmt.Println("The following options are available:")

	// #region Sort Flags

	var flagNames []string

	for flagName, _ := range Flags { // For each flagName in Flags
		flagNames = append(flagNames, flagName) // Append flagName
	}

	sort.Strings(flagNames) // Sort the flagNames

	// #endregion

	for _, flagName := range flagNames { // For each flagName and trimmedFlag struct in Flags
		flag := Flags[flagName] // Get the flag

		thisFlagNameLength := len(flagName)                        // Get the length of this flagName
		flagLengthDiff := (LongestFlagLength - thisFlagNameLength) // Get the difference in flag name length compared to the longest one

		fmt.Println(Config.OSSpecificFlagString + flagName + strings.Repeat(" ", flagLengthDiff) + flag.Descriptor) // Output the flag, creating enough spacing to along descriptor
		if flag.DefaultValue != nil {                                                                               // If DefaultValue is not nil

			if (flag.DefaultValue != "") && (flag.DefaultValue != false) { // If the default value is not an empty string and not false
				fmt.Println("Default Value: ", flag.DefaultValue)
			}
		}

		if flag.AllowNothing { // If we are allowing to pass only the flag
			fmt.Println("Allows Providing Only Flag")
		}

		fmt.Println("")
	}
}
