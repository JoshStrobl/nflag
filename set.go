// This file contains functionality relating to the setting of flags in nflag

package nflag

import (
	"errors"
	"reflect"
)

// Set is for setting a flag.
func Set(flagName string, providedFlag Flag) error {
	var errorResponse error
	var typeFound bool // Default typeFound to false

	if providedFlag.Type == "" { // If no Type was provided
		providedFlag.Type = "bool" // Default to bool
	}

	for _, val := range []string{"bool", "float64", "int", "string"} {
		if providedFlag.Type == val { // If the Type is the same as what is allowed (this specific value)
			typeFound = true // Since this is a valid type

			if providedFlag.DefaultValue == nil { // If no default value was provided
				if providedFlag.Type == "bool" { // If the type is bool
					providedFlag.DefaultValue = false
				} else if providedFlag.Type == "float64" { // If the type is float64
					providedFlag.DefaultValue = 0
				} else if providedFlag.Type == "int" { // If the type is int
					providedFlag.DefaultValue = 0
				} else if providedFlag.Type == "string" {
					providedFlag.DefaultValue = ""
				}
			} else { // If a Default value was provided
				if providedFlag.Type == reflect.TypeOf(providedFlag.DefaultValue).String() { // If the Type and DefaultValue are the same type
					providedFlag.AllowNothing = true // Enforce AllowNothing since a default value is provided
				} else { // If the types are not the same
					errorResponse = errors.New("Mismatch Flag and DefaultValue types.")
				}
			}

			if errorResponse == nil { // If there was no error
				providedFlag.Value = "flag-not-provided" // Default Value to flag-not-provided so we can differ between when a flag isn't provided v.s. no value
				Flags[flagName] = providedFlag           // Set the flag name in Flags to the providedFlag struct

				flagNameLength := len(flagName) // Get the length of the flagName

				if LongestFlagLength < flagNameLength { // If the currently stored LongestFlagLength is less than the current flagNameLength
					LongestFlagLength = (flagNameLength + 4) // Change to this flagNameLength, adding 4 for spacing for longest flag
				}
			}

			break
		}
	}

	if !typeFound { // If this is a non-allowed type
		errorResponse = errors.New("Type is not allowed. Please use bool, float64, int, or string.")
	}

	return errorResponse
}