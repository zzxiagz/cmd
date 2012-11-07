package cmd

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "reflect"
    "strings"
)

type Cmd struct {
    Prompt, Intro string
    client        interface{}
}

var (
    CMD_EXITS []string
    CMD_HELPS []string
    CMD_LISTS []string
)

const (
    ERR_HINT = "ERR:"
    EOL       = '\n'
    SPACE     = " "
    EMPTY_STR = ""
    METHOD_DO = "Do_"

    DEFAULT_PROMPT = "(Cmd) "
)

func init() {
    CMD_EXITS = []string{"exit", "bye", "quit", "q"}
    CMD_HELPS = []string{"help", "h"}
    CMD_LISTS = []string{"list", "ls", "l"}
}

func New(client interface{}) Cmd {
    return Cmd{Prompt: DEFAULT_PROMPT, client: client}
}

func (this Cmd) bye() {
    fmt.Println("bye")
    os.Exit(0)
}

func (this Cmd) intro() {
    if this.Intro != "" {
        fmt.Println(this.Intro)
    }
}

func (this Cmd) isCommand(cmd string, knownCmds []string) bool {
    for _, c := range knownCmds {
        if c == cmd {
            return true
        }
    }
    return false
}

func (this Cmd) isExit(cmd string) bool {
    return this.isCommand(cmd, CMD_EXITS)
}

func (this Cmd) isHelp(cmd string) bool {
    return this.isCommand(cmd, CMD_HELPS)
}

func (this Cmd) isList(cmd string) bool {
    return this.isCommand(cmd, CMD_LISTS)
}

func (this Cmd) readInputs() (cmd string, args []string, err error) {
    // echo prompt before read input
    fmt.Print(this.Prompt)

    reader := bufio.NewReader(os.Stdin)
    in, e := reader.ReadBytes(EOL)
    if e != nil {
        if e != io.EOF {
            panic(e)
        }

        // EOF
        this.bye()
    }

    input := string(in[:len(in)-1]) // discard the EOL
    if strings.TrimSpace(input) == EMPTY_STR {
        err = errors.New("empty input")
        return
    }

    inputs := strings.Split(input, SPACE)
    if len(inputs) > 1 {
        args = make([]string, 0)
        for _, in := range inputs[1:] {
            x := strings.TrimSpace(in)
            x = strings.Trim(x, "'")
            if x != "" {
                args = append(args, x)
            }
        }
    }
    cmd = inputs[0]
    return
}

func (this Cmd) Cmdloop() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(ERR_HINT, r)
            this.Cmdloop()
        }
    }()

    this.intro()

    for {
        cmd, args, err := this.readInputs()
        if err != nil {
            continue
        }

        if this.isExit(cmd) {
            this.bye()
        }

        var method string
        if this.isHelp(cmd) {
            if len(args) >= 1 {
                method = "Help_" + args[0]
                args = args[1:]
            } else {
                method = "Help"
            }
        } else if this.isList(cmd) {
            this.listCommands(this.client)
            continue
        } else {
            method = METHOD_DO + cmd
        }

        this.tryInvoke(this.client, method, args)
    }
}

func (this Cmd) notFound(method string) {
    fmt.Printf("Invalid command: %s\n", method)
}

func (this Cmd) listCommands(i interface{}) {
    fmt.Println("Available commands:")

    t := reflect.TypeOf(i)
    for i:=0; i<t.NumMethod(); i++ {
        methodName := t.Method(i).Name
        if strings.HasPrefix(methodName, METHOD_DO) {
            fmt.Printf("%s ", methodName[len(METHOD_DO):])
        }
    }
    println()
}

func (this Cmd) tryInvoke(i interface{}, methodName string, args []string) {
    var method reflect.Value = reflect.ValueOf(i).MethodByName(methodName)
    if !method.IsValid() {
        this.notFound(methodName)
        return
    }

    params := make([]reflect.Value, len(args))
    for i, arg := range args {
        params[i] = reflect.ValueOf(arg)
    }

    // we don't care about the return value
    method.Call(params)
}
