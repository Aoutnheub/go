package args

import (
    "os"
    "testing"
)

// NOTE: first argument is ignored

func TestParseFlags(t *testing.T) {
    os.Args = []string{"app.exe", "--flag01", "-fF"}
    var parser Parser
    parser.Init("Test", "")
    parser.AddFlag("flag01", "", '\000')
    parser.AddFlag("flag02", "", 'f')
    parser.AddFlag("flag03", "", 'F')
    parser.AddFlag("flag04", "", '\000')
    results, err := parser.Parse()
    if err != nil { t.Error(err) }
    if !results.Flag["flag01"] { t.Error() }
    if !results.Flag["flag02"] { t.Error() }
    if !results.Flag["flag03"] { t.Error() }
    if results.Flag["flag04"] { t.Error() }
}

func TestParseOptions(t *testing.T) {
    os.Args = []string{"app.exe", "--op01", "a", "-O", "b", "--op03=10"}
    var parser Parser
    parser.Init("Test", "")
    parser.AddOption("op01", "", '\000', "", []string{})
    parser.AddOption("op02", "", 'O', "", []string{})
    parser.AddOption("op03", "", '\000', "", []string{"10", "20", "30"})
    parser.AddOption("op04", "", 'o', "hello", []string{})
    results, err := parser.Parse()
    if err != nil { t.Error(err) }
    if results.Option["op01"] != "a" { t.Error() }
    if results.Option["op02"] != "b" { t.Error() }
    if results.Option["op03"] != "10" { t.Error() }
    if results.Option["op04"] != "hello" { t.Error() }
}

func TestParseCommand(t *testing.T) {
    os.Args = []string{"app.exe", "cmd01"}
    var parser Parser
    parser.Init("Test", "")
    parser.AddCommand("cmd01", "")
    parser.AddCommand("cmd02", "")
    results, err := parser.Parse()
    if err != nil { t.Error(err) }
    if results.Command != "cmd01" { t.Error() }
}

func TestParsePositional(t *testing.T) {
    os.Args = []string{"app.exe", "hello", "--flag", "world"}
    var parser Parser
    parser.Init("Test", "")
    parser.AddFlag("flag", "", '\000')
    results, err := parser.Parse()
    if err != nil { t.Error(err) }
    if results.Positional[0] != "hello" { t.Error() }
    if results.Positional[1] != "world" { t.Error() }
}

func TestParseAll(t *testing.T) {
    os.Args = []string{
        "app.exe", "cmd02", "--flag01", "-f", "uwu", "-Otest",
        "--op01", "a", "owo", "-Fo=hi",
    }
    var parser Parser
    parser.Init("Test", "")
    parser.AddFlag("flag01", "", '\000')
    parser.AddFlag("flag02", "", 'f')
    parser.AddFlag("flag03", "", 'F')
    parser.AddFlag("flag04", "", '\000')
    parser.AddOption("op01", "", '\000', "", []string{})
    parser.AddOption("op02", "", 'O', "", []string{})
    parser.AddOption("op03", "", '\000', "", []string{"10", "20", "30"})
    parser.AddOption("op04", "", 'o', "hello", []string{})
    parser.AddCommand("cmd01", "")
    parser.AddCommand("cmd02", "")
    results, err := parser.Parse()
    if err != nil { t.Error(err) }
    if !results.Flag["flag01"] { t.Error() }
    if !results.Flag["flag02"] { t.Error() }
    if !results.Flag["flag03"] { t.Error() }
    if results.Flag["flag04"] { t.Error() }
    if results.Option["op01"] != "a" { t.Error() }
    if results.Option["op02"] != "test" { t.Error() }
    if results.Option["op03"] != "" { t.Error() }
    if results.Option["op04"] != "hi" { t.Error() }
    if results.Command != "cmd02" { t.Error() }
    if results.Positional[0] != "uwu" { t.Error() }
    if results.Positional[1] != "owo" { t.Error() }
}