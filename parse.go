// This file contains functionality relating to the parsing of flags and values in nflag

package nflag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parse
// This function will parse input for flags
func Parse() {
	providedFlags := os.Args[1:] // ProvidedFlags start at index 1 (after binary name)

	if (len(providedFlags) != 0) && (providedFlags[0] != (Config.OSSpecificFlagString + "help")) { // If there was no flags or the first providedFlag is not help
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
	} else { // If no flags were provided or help was
		if len(providedFlags) == 0 { // If no flags were provided
			if !Config.ShowHelpIfNoArgs { // If we should not show help if no args
				for flagName, _ := range Flags { // For each flagName in Flags
					if Flags[flagName].Value == nil { // If no value is set for this Flag (which happens if it wasn't parsed)
						ParseVal(flagName, "flag-not-provided")
					}
				}
			} else { // If we should show help if no args
				OutputHelp = true // Change OutputHelp to true
			}
		} else { // If help was provided
			OutputHelp = true // Change OutputHelp to true
		}
	}

	if OutputHelp { // If we are outputting help flags
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
				fmt.Println("An incorrect value was provided when using " + Config.OSSpecificFlagString + flagName + ". Exiting.")
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
				fmt.Println("A required value for " + Config.OSSpecificFlagString + flagName + " was not provided.")
				OutputHelp = true // Redefine OutputHelp to true
			}
		}

		Flags[flagName] = trimmedFlag // Update the map
	} else { // If the flag does not exist.
		fmt.Println(Config.OSSpecificFlagString + flagName + " does not exist. Exiting.")
		os.Exit(1)
	}
}
