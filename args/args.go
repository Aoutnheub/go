package args

import (
    "errors"
    "fmt"
    "os"
    "strings"
)

type ANSICode string

const (
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
)

type Results struct {
    /// Stores flag values after parsing
    Flag map[string]bool
    /// Stores option values after parsing
    Option map[string]string
    /// Stores positional arguments after parsing
    Positional []string
    /// Stores the command after parsing
    Command string
}

type option struct {
    Help string
    DefaultsTo string
    Allowed []string
}

type Parser struct {
    flags map[string]string
    flagsAbbr map[rune]string
    options map[string]option
    optionsAbbr map[rune]string
    commands map[string]string
    positional []string
    name string
    description string
    cachedHelp string

    /// Return an error if the first argument isn't a command. Ignored if no
    /// commands have been added
    CommandRequired bool
    /// Header displayed by the `Help` function before the command descriptions
    CommandsHelpMsg string
    /// Header displayed by the `Help` function before the flag descriptions
    FlagsHelpMsg string
    /// Header displayed by the `Help` function before the option descriptions
    OptionsHelpMsg string
    /// Color the output of the `Help` function
    Colors bool
    /// Color of the title outputed by the `Help` function
    TitleColor ANSICode
    /// Color of the description outputed by the `Help` function
    DescriptionColor ANSICode
    /// Color of the headers outputed by the `Help` function
    HeaderColor ANSICode
    /// Color of the command names outputed by the `Help` function
    CommandColor ANSICode
    /// Color of the command's description outputed by the `Help` function
    CommandDescriptionColor ANSICode
    /// Color of the flag names outputed by the `Help` function
    FlagColor ANSICode
    /// Color of the flag's description outputed by the `Help` function
    FlagDescriptionColor ANSICode
    /// Color of the option names outputed by the `Help` function
    OptionColor ANSICode
    /// Color of the option's description outputed by the `Help` function
    OptionDescriptionColor ANSICode
    /// Color of the option's allowed values outputed by the `Help` function
    OptionAllowedColor ANSICode
}

func (ap *Parser) getFlagsAbbr() map[string]rune {
    abbr := map[string]rune{}
    for k, v := range ap.flagsAbbr {
        abbr[v] = k
    }

    return abbr
}

func (ap *Parser) getOptionsAbbr() map[string]rune {
    abbr := map[string]rune{}
    for k, v := range ap.optionsAbbr {
        abbr[v] = k
    }

    return abbr
}

func (ap *Parser) isAllowedOptionValue(opt string, val string) bool {
    op := ap.options[opt]
    alLen := len(op.Allowed)
    allowed := false
    if alLen != 0 {
        for i := 0; i < alLen; i++ {
            if op.Allowed[i] == val {
                allowed = true
                break
            }
        }
    }else {
        allowed = true
    }

    return allowed
}

/// Initialize the struct
func (ap *Parser) Init(name string, description string) {
    if name != "" { ap.name = name }
    if description != "" { ap.description = description }
    ap.CommandRequired = false
    ap.CommandsHelpMsg = "COMMANDS"
    ap.FlagsHelpMsg = "FLAGS"
    ap.OptionsHelpMsg = "OPTIONS"
    ap.flags = map[string]string{}
    ap.flagsAbbr = map[rune]string{}
    ap.options = map[string]option{}
    ap.optionsAbbr = map[rune]string{}
    ap.commands = map[string]string{}
    ap.Colors = false
    ap.TitleColor = ANSIGreen
    ap.DescriptionColor = ANSIWhite
    ap.HeaderColor = ANSIRed
    ap.CommandColor = ANSIMagenta
    ap.CommandDescriptionColor = ANSIWhite
    ap.FlagColor = ANSIBlue
    ap.FlagDescriptionColor = ANSIWhite
    ap.OptionColor = ANSIBlue
    ap.OptionDescriptionColor = ANSIWhite
    ap.OptionAllowedColor = ANSIYellow
}

/// Add a flag
/// @param name flag's name
/// @param help flag's description
/// @param abbr flag's abbreviation
/// @return Error is the  flag already exists
func (ap *Parser) AddFlag(name string, help string, abbr rune) error {
    _, foundFl := ap.flags[name]
    _, foundOp := ap.options[name]
    if !foundFl && !foundOp {
        ap.flags[name] = help
        _, foundFlAbr := ap.flagsAbbr[abbr]
        _, foundOpAbr := ap.optionsAbbr[abbr]
        if !foundFlAbr && !foundOpAbr {
            ap.flagsAbbr[abbr] = name
        }else {
            return errors.New(fmt.Sprintf("duplicate argument: %s", string(abbr)))
        }
    }else {
        return errors.New(fmt.Sprintf("duplicate argument: %s", name))
    }

    return nil
}

/// Add an option
/// @param name option's name
/// @param help option's description
/// @param abbr option's abbreviation
/// @param defaultsTo option's default value
/// @param allowed option's allowed values. Doesn't necessarily need to contain the default value
/// @return An error if the option already exists
func (ap *Parser) AddOption(
    name string, help string, abbr rune, defaultsTo string, allowed []string,
) error {
    _, foundOp := ap.options[name]
    _, foundFl := ap.flags[name]
    if !foundOp && !foundFl {
        ap.options[name] = option{Help: help, DefaultsTo: defaultsTo, Allowed: allowed}
        if abbr != '\000' {
            _, foundOpAbr := ap.optionsAbbr[abbr]
            _, foundFlAbr := ap.flagsAbbr[abbr]
            if !foundOpAbr && !foundFlAbr {
                ap.optionsAbbr[abbr] = name
            }else {
                return errors.New(fmt.Sprintf("duplicate argument: %s", string(abbr)))
            }
        }
    }else {
        return errors.New(fmt.Sprintf("duplicate argument: %s", name))
    }

    return nil
}

/// Add a command
/// @param name command's name
/// @param help command's description
/// @return Error if the command already exists
func (ap *Parser) AddCommand(name string, help string) error {
    _, found := ap.commands[name]
    if !found {
        ap.commands[name] = help
    }else {
        return errors.New(fmt.Sprintf("duplicate argument: %s", name))
    }

    return nil
}

/// Display the help message
func (ap *Parser) Help() {
    if ap.name != "" {
        if ap.Colors { fmt.Print(ap.TitleColor) }
        fmt.Print(ap.name)
        if ap.Colors { fmt.Print("\033[0m") }
    }
    if ap.description != "" {
        if ap.Colors { fmt.Print(ap.DescriptionColor) }
        fmt.Print(" - ")
        var indent string
        nameLen := len(ap.name)
        for i := 0; i < nameLen + 3; i++ {
            indent += " "
        }
        token := strings.IndexRune(ap.description, '\n')
        if token != -1 {
            last := 0
            for token != -1 {
                if last == 0 {
                    fmt.Print(ap.description[last:token + 1])
                }else {
                    fmt.Printf("%s%s\n", indent, ap.description[last:token + 1])
                }
                last = token + 1
                token = strings.IndexRune(ap.description[last:], '\n')
            }
            if last < len(ap.description) - 1 {
                fmt.Printf("%s%s\n", indent, ap.description[last:])
            }
        }else {
            fmt.Println(ap.description)
        }
        if ap.Colors { fmt.Print("\033[0m") }
    }
    fmt.Println()

    if len(ap.commands) != 0 {
        if ap.CommandsHelpMsg != "" {
            if ap.Colors { fmt.Print(ap.HeaderColor) }
            fmt.Printf("%s\n", ap.CommandsHelpMsg)
            if ap.Colors { fmt.Print("\033[0m") }
        }
        for k, v := range ap.commands {
            if ap.Colors { fmt.Print(ap.CommandColor) }
            fmt.Printf("    %s\n", k)
            if ap.Colors { fmt.Print("\033[0m") }
            indent := "        "
            if v != "" {
                if ap.Colors { fmt.Print(ap.CommandDescriptionColor) }
                token := strings.IndexRune(v, '\n')
                if token != -1 {
                    last := 0
                    for token != -1 {
                        fmt.Printf("%s%s\n", indent, v[last:token + 1])
                        last = token + 1
                        token = strings.IndexRune(v[last:], '\n')
                    }
                    if last < len(v) - 1 {
                        fmt.Printf("%s%s\n", indent, v[last:])
                    }
                }else {
                    fmt.Printf("%s%s\n", indent, v)
                }
                if ap.Colors { fmt.Print("\033[0m") }
            }
            fmt.Println()
        }
    }

    if len(ap.flags) != 0 {
        abbr := ap.getFlagsAbbr()
        if ap.FlagsHelpMsg != "" {
            if ap.Colors { fmt.Print(ap.HeaderColor) }
            fmt.Printf("%s\n", ap.FlagsHelpMsg)
            if ap.Colors { fmt.Print("\033[0m") }
        }
        for k, v := range ap.flags {
            if ap.Colors { fmt.Print(ap.FlagColor) }
            fmt.Printf("    --%s", k)
            tmp, found := abbr[k]
            if found {
                fmt.Printf(", -%c", tmp)
            }
            if ap.Colors { fmt.Print("\033[0m") }
            fmt.Println()
            indent := "        "
            if v != "" {
                if ap.Colors { fmt.Print(ap.FlagDescriptionColor) }
                token := strings.IndexRune(v, '\n')
                if token != -1 {
                    last := 0
                    for token != -1 {
                        fmt.Printf("%s%s\n", indent, v[last:token + 1])
                        last = token + 1
                        token = strings.IndexRune(v[last:], '\n')
                    }
                    if last < len(v) - 1 {
                        fmt.Printf("%s%s\n", indent, v[last:])
                    }
                }else {
                    fmt.Printf("%s%s\n", indent, v)
                }
                if ap.Colors { fmt.Print("\033[0m") }
            }
            fmt.Println()
        }
    }

    if len(ap.options) != 0 {
        abbr := ap.getOptionsAbbr()
        if ap.OptionsHelpMsg != "" {
            if ap.Colors { fmt.Print(ap.HeaderColor) }
            fmt.Printf("%s\n", ap.OptionsHelpMsg)
            if ap.Colors { fmt.Print("\033[0m") }
        }
        for k, v := range ap.options {
            if ap.Colors { fmt.Print(ap.OptionColor) }
            fmt.Printf("    --%s", k)
            tmp, found := abbr[k]
            if found {
                fmt.Printf(", -%c", tmp)
            }
            if ap.Colors { fmt.Print("\033[0m") }
            allowedLen := len(v.Allowed)
            if allowedLen != 0 {
                if ap.Colors { fmt.Print(ap.OptionAllowedColor) }
                fmt.Print(" ")
                for ii := 0; ii < allowedLen; ii++ {
                    if ii != allowedLen - 1 {
                        fmt.Printf("%s|", v.Allowed[ii])
                    }else {
                        fmt.Print(v.Allowed[ii])
                    }
                }
                if ap.Colors { fmt.Print("\033[0m") }
            }
            fmt.Println()
            indent := "        "
            if v.Help != "" {
                if ap.Colors { fmt.Print(ap.OptionDescriptionColor) }
                token := strings.IndexRune(v.Help, '\n')
                if token != -1 {
                    last := 0
                    for {
                        fmt.Printf("%s%s", indent, v.Help[last:token + 1])
                        last = token + 1
                        token = strings.IndexRune(v.Help[last:], '\n')
                        if token != -1 {
                            token += last
                        }else {
                            break
                        }
                    }
                    if last < len(v.Help) - 1 {
                        fmt.Printf("%s%s\n", indent, v.Help[last:])
                    }
                }else {
                    fmt.Printf("%s%s\n", indent, v.Help)
                }
                if ap.Colors { fmt.Print("\033[0m") }
            }
            fmt.Println()
        }
    }
}

/// Parse the command line arguments
/// @return A "Results" struct with the argument values or error
func (ap *Parser) Parse() (*Results, error) {
    results := new(Results)
    results.Flag = map[string]bool{}
    results.Option = map[string]string{}

    for k := range ap.flags {
        results.Flag[k] = false
    }
    for k, v := range ap.options {
        results.Option[k] = v.DefaultsTo
    }

    i := 0
    skipCommandCheck := false
    args := os.Args[1:]
    argsLen := len(args)
    for i < argsLen {
        curArgLen := len(args[i])
        if !skipCommandCheck && i == 0 && len(ap.commands) != 0 {
            _, found := ap.commands[args[i]]
            if found {
                results.Command = args[i]
                i++
            }else {
                if ap.CommandRequired {
                    return nil, errors.New(
                        fmt.Sprintf("invalid argument: \"%s\" is not a command", args[i]),
                    )
                }
            }
            skipCommandCheck = true
        }else {
            if args[i] == "--" {
                for ii := i + 1; ii < argsLen; ii++ {
                    results.Positional = append(results.Positional, args[ii])
                }
                i = argsLen
            }else {
                if curArgLen > 2 {
                    if args[i][0] == '-' && args[i][1] != '-' {
                        equals := strings.IndexRune(args[i], '=')
                        if equals != -1 {
                            if equals == 2 { // option
                                op, found := ap.optionsAbbr[rune(args[i][1])]
                                if found {
                                    if curArgLen > 3 {
                                        tmp := args[i][3:]
                                        if ap.isAllowedOptionValue(op, tmp) {
                                            results.Option[op] = tmp
                                            i++
                                        }else {
                                            return nil, errors.New(
                                                fmt.Sprintf("invalid value: %s -> %s", string(args[i][1]), tmp),
                                            )
                                        }
                                    }else {
                                        return nil, errors.New(
                                            fmt.Sprintf("missing value: %s", string(args[i][1])),
                                        )
                                    }
                                }
                            }else { // multiple flags and one option
                                for ii := 1; ii < equals - 1; ii++ {
                                    fl, found := ap.flagsAbbr[rune(args[i][ii])]
                                    if found {
                                        results.Flag[fl] = true
                                    }else {
                                        return nil, errors.New(
                                            fmt.Sprintf("invalid argument: flag %s does not exist", string(args[i][ii])),
                                        )
                                    }
                                }
                                if equals + 1 < curArgLen {
                                    op, found := ap.optionsAbbr[rune(args[i][equals - 1])]
                                    if found {
                                        tmp := args[i][equals + 1:]
                                        if ap.isAllowedOptionValue(op, tmp) {
                                            results.Option[op] = tmp
                                        }else {
                                            return nil, errors.New(
                                                fmt.Sprintf("missing value: %s", string(args[i][1])),
                                            )
                                        }
                                    }
                                }else {
                                    return nil, errors.New(
                                        fmt.Sprintf("missing value: %s", string(args[i][1])),
                                    )
                                }
                                i++
                            }
                        }else { // option and value with no space
                            op, found := ap.optionsAbbr[rune(args[i][1])]
                            if found {
                                tmp := args[i][2:]
                                if ap.isAllowedOptionValue(op, tmp) {
                                    results.Option[op] = tmp
                                }else {
                                    return nil, errors.New(
                                        fmt.Sprintf("missing value: %s", string(args[i][1])),
                                    )
                                }
                            }else { // multiple flags
                                for ii := 1; ii < curArgLen; ii++ {
                                    fl, found := ap.flagsAbbr[rune(args[i][ii])]
                                    if found {
                                        results.Flag[fl] = true
                                    }else {
                                        return nil, errors.New(
                                            fmt.Sprintf("invalid argument: flag %s does not exist", string(args[i][ii])),
                                        )
                                    }
                                }
                            }
                            i++
                        }
                    }else if args[i][0] == '-' && args[i][1] == '-' {
                        equals := strings.IndexRune(args[i], '=')
                        if equals != -1 {
                            op := args[i][2:equals]
                            var val string
                            if equals + 1 < curArgLen {
                                val = args[i][equals + 1:]
                            }else {
                                return nil, errors.New(
                                    fmt.Sprintf("missing value: %s", op),
                                )
                            }
                            if ap.isAllowedOptionValue(op, val) {
                                results.Option[op] = val
                            }else {
                                return nil, errors.New(
                                    fmt.Sprintf("invalid value: %s -> %s", op, val),
                                )
                            }
                            i++
                        }else {
                            arg := args[i][2:]
                            _, found := ap.flags[arg]
                            if found {
                                results.Flag[arg] = true
                                i++
                            }else {
                                _, found := ap.options[arg]
                                if found {
                                    if i + 1 < argsLen {
                                        if len(args[i + 1]) == 0 {
                                            results.Option[arg] = args[i + 1]
                                        }else {
                                            if args[i + 1][0] != '-' {
                                                if ap.isAllowedOptionValue(arg, args[i + 1]) {
                                                    results.Option[arg] = args[i + 1]
                                                }else {
                                                    return nil, errors.New(
                                                        fmt.Sprintf("invalid value: %s -> %s", arg, args[i + 1]),
                                                    )
                                                }
                                            }else {
                                                return nil, errors.New(
                                                    fmt.Sprintf("missing value: %s", args[i][2:]),
                                                )
                                            }
                                        }
                                    }else {
                                        return nil, errors.New(
                                            fmt.Sprintf("missing value: %s", args[i][2:]),
                                        )
                                    }
                                    i += 2
                                }else {
                                    return nil, errors.New(
                                        fmt.Sprintf("invalid argument: %s", arg),
                                    )
                                }
                            }
                        }
                    }else {
                        results.Positional = append(results.Positional, args[i])
                        i++
                    }
                }else if curArgLen == 2 {
                    if args[i][0] == '-' {
                        fl, found := ap.flagsAbbr[rune(args[i][1])]
                        if found {
                            results.Flag[fl] = true
                            i++
                        }else {
                            op, found := ap.optionsAbbr[rune(args[i][1])]
                            if found {
                                if i + 1 < argsLen {
                                    if len(args[i + 1]) == 0 {
                                        results.Option[op] = args[i + 1]
                                    }else {
                                        if args[i + 1][0] != '-' {
                                            results.Option[op] = args[i + 1]
                                        }else {
                                            return nil, errors.New(
                                                fmt.Sprintf("missing value: %s", args[i]),
                                            )
                                        }
                                    }
                                }else {
                                    return nil, errors.New(
                                        fmt.Sprintf("missing value: %s", args[i]),
                                    )
                                }
                                i += 2
                            }else {
                                return nil, errors.New(
                                    fmt.Sprintf("invalid argument: %s", string(args[i][1])),
                                )
                            }
                        }
                    }else {
                        results.Positional = append(results.Positional, args[i])
                        i++
                    }
                }else {
                    results.Positional = append(results.Positional, args[i])
                    i++
                }
            }
        }
    }

    return results, nil
}
