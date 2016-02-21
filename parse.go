// This file contains functionality relating to the parsing of flags and values in nflag

package nflag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parse will parse input for flags
func Parse() {
	flagsToParse := Flags                  // Parse all flags by default
	providedFlags := os.Args[1:]           // ProvidedFlags start at index 1 (after binary name)
	lenProvidedFlags := len(providedFlags) // Get the length of the providedFlags

	if (lenProvidedFlags == 0 && Config.ShowHelpIfNoArgs) || ((lenProvidedFlags != 0) && (providedFlags[0] == Config.FlagString+"help")) { // If no args were provided or the first was help, and we should show help if no args
		PrintFlags() // Output help / print flags
	} else if lenProvidedFlags != 0 { // If flags were provided
		for _, flag := range providedFlags { // For each flag provided
			flagNameValueSplit := strings.Split(flag, "=")                                // Split the flag name and the value
			flagName := strings.Replace(flagNameValueSplit[0], Config.FlagString, "", -1) // Get the flagName and remove the FlagString

			existingFlag, exists := flagsToParse[flagName] // Get the existing flag and the exists bool

			if exists { // If this flag exists
				if len(flagNameValueSplit) == 2 { // If a value was provided
					existingFlag.Value = flagNameValueSplit[1] // Set to the value in flagNameValueSplit
				} else { // If no value was provided
					existingFlag.Value = "" // Set to nothing so we can differ between a flag passed and not passed
				}

				flagsToParse[flagName] = existingFlag // Add this flag name and struct to flagsToParse
			} else { // If the Flag does not exist
				fmt.Println(Config.FlagString + flagName + " does not exist.")
				PrintFlags()
			}
		}
	}

	for flagName, flag := range flagsToParse { // For each flagName and flag in flagsToParse
		if flag.Value != "flag-not-provided" { // If the Value was provided / this flag was passed
			if flag.Value != "" { // If a value was provided when using this flag
				var conversionError error // Define conversionError as any error from not correctly converting the value to the type
				flagValueAsString := flag.Value.(string)

				if flag.Type == "bool" { // If the type is bool
					flag.Value, conversionError = strconv.ParseBool(flagValueAsString) // Convert to bool
				} else if flag.Type == "float64" { // If the type is float64
					flag.Value, conversionError = strconv.ParseFloat(flagValueAsString, 64) // Convert to float64
				} else if flag.Type == "int" { // If the type is int
					flag.Value, conversionError = strconv.Atoi(flagValueAsString) // Convert to int
				}

				if conversionError != nil { // If there was a conversionError
					fmt.Println("An incorrect value was provided when using " + Config.FlagString + flagName + ".")
					PrintFlags()
				}
			} else { // If a value was not provided when using this flag
				if flag.AllowNothing { // If we are allowing nothing to be defined
					if flag.Type == "bool" { // If the flag.Type is bool
						flag.Value = true // Set to true
					} else { // If the flag.Type is not bool
						flag.Value = flag.DefaultValue // Change to DefaultValue
					}
				} else { // If we do not allow anything
					fmt.Println("A value must be provided when using " + Config.FlagString + flagName)
					PrintFlags()
				}
			}
		} else { // If the flag was not provided
			if !flag.Required { // If the flag is NOT required
				flag.Value = flag.DefaultValue // Change to DefaultValue
			} else { // If the flag was not passed but is required
				fmt.Println("A required value for " + Config.FlagString + flagName + " was not provided.")
				PrintFlags()
			}
		}

		Flags[flagName] = flag // Update the Flag if we get to this point
	}
}
