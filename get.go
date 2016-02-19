// This file contains the functionality relating to the getting of values from nflag

package nflag

import (
	"errors"
)

// Get
// This function will get the flag value and returns it, or an error if the flag does not exist.
func Get(flagName string) (interface{}, error) {
	flag, exists := Flags[flagName] // Get the trimmedFlag and the exists bool of this flagName in Flags
	var val interface{}
	var err error

	if exists { // If the flag exists
		if flag.Value != "flag-not-provided" { // If the value of the flag is not "flag-not-provided"
			val = flag.Value
		} else { // If the flag of the value is nil
			err = errors.New("Value of " + flagName + " is nil.")
		}
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

// IsDefaultValue
// This function will return a boolean as to whether or not the value of the flag is the same as DefaultValue
func IsDefaultValue(flagName string) (bool, error) {
	var isDefaultValue bool
	flagValue, flagValueError := Get(flagName)

	if flagValueError == nil { // If the flag exists
		if Flags[flagName].DefaultValue == flagValue { // If the DefaultValue is the same as the flagValue
			isDefaultValue = true // Change isDefaultValue to true
		}
	}

	return isDefaultValue, flagValueError
}
