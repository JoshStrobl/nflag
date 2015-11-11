package nflag

type ConfigOptions struct {
	OSSpecificFlags      bool
	OSSpecificFlagString string
	ShowHelpIfNoArgs     bool
}

// Flag is used in the Config string map
type Flag struct {
	Descriptor string
	// Type, DefaultValue, Value: Must be bool, float64, int, or string
	Type                string
	DefaultValue, Value interface{}
	// Required - If this flag is required to be passed
	Required bool
	// Allow passing of no value, mainly for triggering certain actions
	AllowNothing bool
}
