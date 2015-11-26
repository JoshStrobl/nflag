## About ##

nflag is an open source alternative to Golang's flag package. nflag allows flag styling, type option simplification, and allowance of empty values.

nflag enables you to do things like:
1. Enforce required flags.
2. Enable optional flags with both default values and no values (we will then provide a default value).
3. Use OS-specific flag styling or whatever style you like.
4. Print help message if no arguments are provided.

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

If the developer has specified during nflag configuration to show help if no arguments are passed, then the flags will be printed simply when doing `./executable`.

## Usage (Developer) ##

You should probably just use `godoc` but I mean, whatever works for you I guess. This is mainly for convenience.

### Structs ###

**ConfigOptions**:

``` go
type ConfigOptions struct {
    OSSpecificFlags      bool
    OSSpecificFlagString string
    ShowHelpIfNoArgs     bool
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

We also have the following helper functions, all of which take the `flagName`:

- `GetAsBool`: returns `bool, error`
- `GetAsFloat64`: returns `float64, error`
- `GetAsInt`: returns `int, error`
- `GetAsString`: returns `string, error`

#### IsDefaultValue ####

This function will return a boolean as to whether or not the value of the flag is the same as DefaultValue.

``` go
func IsDefaultValue(flagName string) (bool, error)
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

Alternatively, you could do:

``` go
int, _ := nflag.GetAsInt("number-of-people")
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
{Config.OSSpecificFlagString}{flagName} {Flag Descriptor}
{If Default Value - "Default Value: {Flag Default Value}"}
{If Allows Nothing - "Allows Providing Only Flag"}
```

**Example**:

``` bash
Usage: --example=value
The following options are available:
--required    This is to test required flag.
Allows Providing Only Flag

--test        This is a test flag.
```