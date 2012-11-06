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

var CMD_EXITS []string

const (
    EOL       = '\n'
    SPACE     = " "
    EMPTY_STR = ""
    METHOD_DO = "Do_"

    DEFAULT_PROMPT = "(Cmd) "

    CMD_HELP = "help"
)

func init() {
    CMD_EXITS = []string{"exit", "bye", "quit"}
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

func (this Cmd) isExit(cmd string) bool {
    for _, c := range CMD_EXITS {
        if c == cmd {
            return true
        }
    }
    return false
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
        args = inputs[1:]
    }
    cmd = inputs[0]
    return
}

func (this Cmd) Cmdloop() {
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
        if cmd == CMD_HELP {
            if len(args) >= 1 {
                method = "Help_" + args[0]
                args = args[1:]
            } else {
                method = "Help"
            }
        } else {
            method = METHOD_DO + cmd
        }

        this.tryInvoke(this.client, method, args)
    }
}

func (this Cmd) notFound(method string) {
    fmt.Printf("Invalid method: %s\n", method)
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
