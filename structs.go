package nflag

type ConfigOptions struct {
	OSSpecificFlags      bool
	OSSpecificFlagString string
}

// Flag is used in the Config string map
type Flag struct {
	Descriptor string
	// Type, DefaultValue, Value: Must be bool, float64, int, or string
	Type                string
	DefaultValue, Value interface{}
	// Allow passing of no value, mainly for triggering certain actions
	AllowNothing bool
}
