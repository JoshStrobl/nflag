## About ##

nflag is an open source alternative to Golang's flag package. nflag allows flag styling, type option simplification, and allowance of empty values.

nflag enables you to do things like:
1. Enforce required flags.
2. Enable optional flags with both default values and no values (we will then provide a default value).
3. Use OS-specific flag styling or whatever style you like.

### OS-Specific Flag Styling (Defaults) ###

**Non-Windows:** If you are using a non-Windows operating system, we default to using `--` in front of the flag name.

**Windows:** If you are using a Windows operating system, we default to using `/` in front of the flag name.

## Usage (User) ##

In the example below, we are:

1. Required `build-for` flag.
2. Having `enable-awesomeness` be optional `bool` flag. By passing this with no value, it'll default to the type specified (in this fake instance, `true`).

``` bash
./executable --build-for=linux_amd64 --enable-awesomeness
```

Calling `help`, prepended by the OSSpecificFlagString, will print the flags.

Example call:

``` bash
./executable --help
```

## Usage (Developer) ##

You should probably just use `godoc` but I mean, whatever works for you I guess. This is mainly for convenience.

### Structs ###

**ConfigOptions**:

``` go
type ConfigOptions struct {
    OSSpecificFlags      bool
    OSSpecificFlagString string
}
```

**Flag**:

``` go
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

```

### Variables ###

``` go
var Config ConfigOptions
var Flags map[string]Flag
```

### Methods ###

*The only ones that really matter.*

#### Configure ####

This function is for configuration of nflag prior to usage. This is an **optional** function, if you don't care to change any configuration options.

``` go
func Configure(providedConfig ConfigOptions)
```

#### Get ####

This function will get the flag value and returns it, or an error if the flag does not exist.

``` go
func Get(flagName string) interface{}, error
```

#### Set ####

This function is for setting a flag.

``` go
func Set(flagName string, providedFlag Flag) error
```

We will return an error:
1. If you provide a Type and DefaultValue but their types don't match.
2. If you provide a type other than: `bool`, `float64`, `int`, or `string`.

#### Parse ####

This is a function that will parse `os.Args` and do nflag magic. You are required to call this after all your `Set()` calls.

``` go
func Parse()
```

#### PrintFlags ####

This function will print all the flags that are set and their defaults.

``` go
func PrintFlags()
```

## Example Usage ##

### Configure ###

The below `Configure` call will change `--` to `-`:

``` go
    nflag.Configure(nflag.ConfigOptions{OSSpecificFlagString: "-"})
```

### Get ###

Let's pretend we set a flag via `Set` before, called `number-of-people`, which is an `int`. Our below function call will return the `Value`.

``` go
intInterface, _ := nflag.Get("number-of-people")
int := intInterface.(int)
```

### Set ###

The below `Set` call will create a required flag called `build-for`, where you must pass a string:

``` go
nflag.Set("build-for", nflag.Flag{
    Descriptor: "What are we building this for.",
    Type : "string",
    Required : true,
})
```

The below `Set` call will create an optional flag called `test` with a default string of `defaultval`. We are allowing nothing to be provided, since we have a default value:

``` go
nflag.Set("test", nflag.Flag{
    Descriptor: "This is a test flag.",
    Type: "string",
    DefaultValue: "defaultval",
    AllowNothing : true,
})
```

The below `Set` call will create an optional flag called `nothing`. We are allowing nothing to be provided and assume the value is true for this `bool` flag, since this flag was passed in the first place.

``` go
nflag.Set("nothin", nflag.Flag{
    Descriptor: "This is to test nothing.",
    Type : "bool",
    // No need to provide a defaultval for bool if AllowNothing is true
    AllowNothing: true,
})
```

### PrintFlags ###

Calling `PrintFlags` will output your flags in the following format:

```
Usage: --example=value
The following options are available:
{Config.OSSpecificFlagString}{flagName} {Flag Type}
    {Flag Descriptor}
    {If Default Value - "Default Value: {Flag Default Value}"}
    {"Allows Providing Only Flag: " {Flag AllowsNothing}}
```

**Example**:

``` bash
Usage: --example=value
The following options are available:
--required bool
	This is to test required flag.
	Allows Providing Only Flag: false

--test string
	This is a test flag.
	Default Value: defaultval
	Allows Providing Only Flag: true
```