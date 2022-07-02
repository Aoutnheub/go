## Description

Command line argument parser.

## Documentation

### Types

```go
ANSICode string
Results struct
Parser struct
```

### Constants

```go
ANSIDefault ANSICode = "\033[39m"
ANSIBGDefault ANSICode = "\033[49m"
ANSIBlack ANSICode = "\033[30m"
ANSIRed ANSICode = "\033[31m"
ANSIGreen ANSICode = "\033[32m"
ANSIYellow ANSICode = "\033[33m"
ANSIBlue ANSICode = "\033[34m"
ANSIMagenta ANSICode = "\033[35m"
ANSICyan ANSICode = "\033[36m"
ANSIWhite ANSICode = "\033[37m"
ANSIBGBlack ANSICode = "\033[40m"
ANSIBGRed ANSICode = "\033[41m"
ANSIBGGreen ANSICode = "\033[42m"
ANSIBGYellow ANSICode = "\033[43m"
ANSIBGBlue ANSICode = "\033[44m"
ANSIBGMagenta ANSICode = "\033[45m"
ANSIBGCyan ANSICode = "\033[46m"
ANSIBGWhite ANSICode = "\033[47m"
```

### Struct fields

#### Results

- `Flag map[string]bool`

    Stores flag values after parsing

- `Option map[string]string`

    Stores option values after parsing

- `Positional []string`

    Stores positional arguments after parsing

- `Command string`

    Stores the command after parsing

#### Parser

- `CommandRequired bool` default: `false`

    Return an error if the first argument isn't a command. Ignored if no commands have been added

- `CommandsHelpMsg string` default: `"COMMANDS"`

    Header displayed by the `Help` function before the command descriptions

- `FlagsHelpMsg string` default: `"FLAGS"`

    Header displayed by the `Help` function before the flag descriptions

- `OptionsHelpMsg string` default: `"OPTIONS"`

    Header displayed by the `Help` function before the option descriptions

- `Colors bool` default: `false`

    Color the output of the `Help` function

- `TitleColor ANSICode` default: `ANSIGreen`

    Color of the title outputed by the `Help` function

- `DescriptionColor ANSICode` default: `ANSIWhite`

    Color of the description outputed by the `Help` function

- `HeaderColor ANSICode` default: `ANSIRed`

    Color of the headers outputed by the `Help` function

- `CommandColor ANSICode` default: `ANSIMagenta`

    Color of the command names outputed by the `Help` function

- `CommandDescriptionColor ANSICode` default: `ANSIWhite`

    Color of the command's description outputed by the `Help` function

- `FlagColor ANSICode` default: `ANSIBlue`

    Color of the flag names outputed by the `Help` function

- `FlagDescriptionColor ANSICode` default: `ANSIWhite`

    Color of the flag's description outputed by the `Help` function

- `OptionColor ANSICode` default: `ANSIBlue`

    Color of the option names outputed by the `Help` function

- `OptionDescriptionColor ANSICode` default: `ANSIWhite`

    Color of the option's description outputed by the `Help` function

- `OptionAllowedColor ANSICode` default: `ANSIYellow`

    Color of the option's allowed values outputed by the `Help` function

### Struct methods

#### Parser

- `Init(name string, description string)`

    Initialize the struct

- `AddFlag(name string, help string, abbr rune) error`

    Add a flag

    - `name` flag's name
    - `help` flag's description
    - `abbr` flag's abbreviation

    **Returns**: An error if the flag already exists

- `AddOption(name string, help string, abbr rune, defaultsTo string, allowed []string) error`

    Add an option

    - `name` option's name
    - `help` option's description
    - `abbr` option's abbreviation
    - `defaultsTo` option's default value
    - `allowed` option's allowed values. Doesn't necessarily need to contain the default value

    **Returns**: An error if the option already exists

- `AddCommand(name string, help string) error`

    Add a command

    - `name` command's name
    - `help` command's description

    **Returns**: An error if the command already exists

- `Help()`

    Display the help message

- `Parse() (*Results, error)`

    Parse the command line arguments

    **Returns**: A `Results` struct with the argument values or error

## Example

```go
var parser args.Parser
parser.Init("Test", "This is a test program")
parser.AddCommand("run", "This is a test command")
parser.AddFlag("flag", "This is a test flag", 'f')
parser.AddOption("option", "This is a test option", 'o', "v1", []string{"v1", "v2"})

// Display the help message and exit if there are no arguments
if len(os.Args) < 2 {
    parser.Help()
    os.Exit(0)
}

results, _ := parser.Parse()
results.Command // Acess command
results.Flag["flag"] // Acess flags
results.Option["option"] // Acess options
results.Positional[0] // Acess positional arguments
```
