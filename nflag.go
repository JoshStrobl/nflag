// nflag is an open source alternative to Golang's flag package.
// nflag allows OS-specific flag styling, type option simplification, and allowance of empty values.

package nflag

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var Config ConfigOptions
var Flags map[string]Flag

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
	if providedConfig.OSSpecificFlags == false { // If we are overriding OSSpecificFlags
		Config.OSSpecificFlags = false
		Config.OSSpecificFlagString = "--" // Enforce --
	}

	if providedConfig.OSSpecificFlagString != "" { // If we are overriding OSSpecificFlagString
		Config.OSSpecificFlagString = providedConfig.OSSpecificFlagString // Set OSSpecificFlagString to whatever is provided on config
	}

	if providedConfig.ShowHelpIfNoArgs { // If ShowHelpIfNoArgs was set to true
		Config.ShowHelpIfNoArgs = true // Change to true
	}
}

// Get
// This function will get the flag value and returns it, or an error if the flag does not exist.
func Get(flagName string) (interface{}, error) {
	flag, exists := Flags[flagName] // Get the trimmedFlag and the exists bool of this flagName in Flags
	var val interface{}
	var err error

	if exists { // If the flag exists
		val = flag.Value
	} else {
		err = errors.New("Flag does not exist.")
	}

	return val, err
}

// #region GetAs functions

// GetAsBool
// This function will get the flag value and convert it to bool, or an error if the flag does not exist.
func GetAsBool(flagName string) (bool, error) {
	var flagValue bool
	flagValueInterface, flagValueError := Get(flagName)

	if flagValueError == nil { // If this flag does exist
		flagValue = flagValueInterface.(bool) // Type inference to bool
	}

	return flagValue, flagValueError
}

// GetAsInt
// This function will get the flag value and convert it to int, or an error if the flag does not exist.
func GetAsInt(flagName string) (int, error) {
	var flagValue int
	flagValueInterface, flagValueError := Get(flagName)

	if flagValueError == nil { // If this flag does exist
		flagValue = flagValueInterface.(int) // Type inference to int
	}

	return flagValue, flagValueError
}

// GetAsFloat64
// This function will get the flag value and convert it to float64, or an error if the flag does not exist.
func GetAsFloat64(flagName string) (float64, error) {
	var flagValue float64
	flagValueInterface, flagValueError := Get(flagName)

	if flagValueError == nil { // If this flag does exist
		flagValue = flagValueInterface.(float64) // Type inference to float64
	}

	return flagValue, flagValueError
}

// GetAsString
// This function will get the flag value and convert it to string, or an error if the flag does not exist.
func GetAsString(flagName string) (string, error) {
	var flagValue string
	flagValueInterface, flagValueError := Get(flagName)

	if flagValueError == nil { // If this flag does exist
		flagValue = flagValueInterface.(string) // Type inference to string
	}

	return flagValue, flagValueError
}

// #endregion

// Set
// This function is for setting a flag.
func Set(flagName string, providedFlag Flag) error {
	var errorResponse error

	if providedFlag.Type == "" { // If no Type was provided
		providedFlag.Type = "bool" // Default to bool
	}

	allowedTypes := []string{"bool", "float64", "int", "string"} // Define allowed types
	isAllowedType := false                                       // Default to providedFlag.Type not being the type

	for _, val := range allowedTypes {
		if providedFlag.Type == val { // If the Type is the same as what is allowed (this specific value)
			isAllowedType = true
		}
	}

	if isAllowedType { // If this is an allowed type
		if (providedFlag.DefaultValue == nil) && (providedFlag.AllowNothing) { // If no default value was provided and we are allowing nothing
			if providedFlag.Type == "bool" { // If the type is bool
				providedFlag.DefaultValue = false
			} else if providedFlag.Type == "float64" { // If the type is float64
				providedFlag.DefaultValue = 1.1
			} else if providedFlag.Type == "int" { // If the type is int
				providedFlag.DefaultValue = 1
			} else if providedFlag.Type == "string" {
				providedFlag.DefaultValue = ""
			}
		} else if (providedFlag.DefaultValue != nil) && (providedFlag.AllowNothing == false) { // If a Default value was provided
			providedFlag.AllowNothing = true                                             // Enforce AllowNothing since a default value is provided
			if providedFlag.Type != reflect.TypeOf(providedFlag.DefaultValue).String() { // If the Type and DefaultValue aren't the same type
				errorResponse = errors.New("Mismatch Flag and DefaultValue types.")
			}
		}

		if errorResponse == nil { // If there was no error
			Flags[flagName] = providedFlag // Set the flag name in Flags to the providedFlag struct
		}
	} else { // If this is a non-allowed type
		errorResponse = errors.New("Type is not allowed. Please use bool, float64, int, or string.")
	}

	return errorResponse
}

// Parse
// This function will parse input for flags
func Parse() {
	providedFlags := os.Args[1:] // ProvidedFlags start at index 1 (after binary name)

	if len(providedFlags) != 0 { // If flags were provided
		if providedFlags[0] != (Config.OSSpecificFlagString + "help") { // If the first providedFlag is not help
			for _, flag := range providedFlags { // For each flag provided
				flagNameValueSplit := strings.Split(flag, "=")                                          // Split the flag name and the value
				flagName := strings.Replace(flagNameValueSplit[0], Config.OSSpecificFlagString, "", -1) // Get the flagName and remove the OSSpecificFlagString
				flagValue := ""                                                                         // Default to nothing

				if len(flagNameValueSplit) == 2 { // If a value was provided
					flagValue = flagNameValueSplit[1] // Set to the value in flagNameValueSplit
				} else { // If the flag was passed but no value
					trimmedFlag, exists := Flags[flagName] // Get the trimmedFlag and the exists bool of this flagName in Flags

					if exists && (trimmedFlag.Type == "bool") { // If the flag exists and its type is bool
						flagValue = "true" // By passing just --flag, it implies --flag is true
					}
				}

				ParseVal(flagName, flagValue) // Parse the value
			}
		} else { // If we should output the help message
			PrintFlags() // Print the flags
			os.Exit(1)
		}

		for flagName, _ := range Flags { // For each flagName in Flags
			if Flags[flagName].Value == nil { // If no value is set for this Flag (which happens if it wasn't parsed via
				ParseVal(flagName, "flag-not-provided")
			}
		}
	} else if Config.ShowHelpIfNoArgs { // If no arguments are passed and ShowHelpIfNoArgs is true
		PrintFlags()
		os.Exit(1)
	}
}

// ParseVal
// This function will check Flags for the provided flag and parse the provided value
func ParseVal(flagName string, flagValue string) {
	trimmedFlag, exists := Flags[flagName] // Get the trimmedFlag and the exists bool of this flagName in Flags

	if exists { // If the flag exists
		if (flagValue != "") && (flagValue != "flag-not-provided") { // If a value was defined
			var conversionError error // Define conversionError as any error from not correctly converting the value to the type

			if trimmedFlag.Type == "bool" { // If the type is bool
				trimmedFlag.Value, conversionError = strconv.ParseBool(flagValue) // Convert to bool
			} else if trimmedFlag.Type == "float64" { // If the type is float64
				trimmedFlag.Value, conversionError = strconv.ParseFloat(flagValue, 64) // Convert to float64
			} else if trimmedFlag.Type == "int" { // If the type is int
				trimmedFlag.Value, conversionError = strconv.Atoi(flagValue) // Convert to int
			} else { // Remaining type, being int
				trimmedFlag.Value = flagValue
			}

			if conversionError != nil { // If there was a conversionError
				fmt.Println("An incorrect value was provided when using " + Config.OSSpecificFlagString + "flagName . Exiting.")
				os.Exit(1)
			}
		} else { // If no value was provided by the user
			returnRequiredValueMessage := false

			if flagValue == "" { // If the flag was provided by the user but no content
				if trimmedFlag.AllowNothing { // If providing no content is permitted
					if trimmedFlag.Type == "bool" { // If the flag is bool
						trimmedFlag.Value = true // Force true
					} else { // If the flag is not bool
						trimmedFlag.Value = trimmedFlag.DefaultValue // Set to defaultvalue
					}
				} else { // If providing no content is NOT permitted
					returnRequiredValueMessage = true
				}
			} else { // If this flag was not provided
				if !trimmedFlag.Required { // If the flag is NOT required
					trimmedFlag.Value = trimmedFlag.DefaultValue // Set the value to the default
				} else { // If the flag IS required
					returnRequiredValueMessage = true
				}
			}

			if returnRequiredValueMessage {
				fmt.Println("A required value for " + Config.OSSpecificFlagString + flagName + " was not provided. Exiting.")
				os.Exit(1)
			}
		}

		Flags[flagName] = trimmedFlag // Update the map
	} else { // If the flag does not exist.
		fmt.Println(Config.OSSpecificFlagString + flagName + " does not exist. Exiting.")
		os.Exit(1)
	}
}

// PrintFlags
// This function will print all the flags that are set and their defaults
func PrintFlags() {
	fmt.Println("Usage: " + Config.OSSpecificFlagString + "example=value")
	fmt.Println("The following options are available:")

	for flagName, flag := range Flags { // For each flagName and trimmedFlag struct in Flags
		fmt.Println(Config.OSSpecificFlagString + flagName + " " + flag.Type)
		fmt.Println("\t" + flag.Descriptor)

		if flag.DefaultValue != nil { // If DefaultValue is not nil
			fmt.Println("\t" + "Default Value: " + reflect.ValueOf(flag.DefaultValue).String())
		}

		fmt.Println("\t" + "Allows Providing Only Flag: " + strconv.FormatBool(flag.AllowNothing))
		fmt.Println("")
	}
}
