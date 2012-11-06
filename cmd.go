package cmd

import (
    "bufio"
    "os"
    "fmt"
    "reflect"
    "strings"
)

const (
    CMD_HELP = "help"
    CMD_EXIT = "exit"
    CMD_BYE = "bye"
    CMD_QUIT = "quit"
)

func New(child interface{}) Cmd {
    this := Cmd{}
    this.child = child
    fmt.Println(child)
    this.Prompt = "(Cmd) "
    return this
}

func (this Cmd) bye() {
    fmt.Println("bye")
    os.Exit(0)
}

func (this Cmd) Cmdloop() {
    for {
        fmt.Print(this.Prompt)
        reader := bufio.NewReader(os.Stdin)
        in, e := reader.ReadBytes('\n')
        if e != nil {
            this.bye()
        }

        inputs := strings.Split(string(in[:len(in)-1]), " ")
        cmd := inputs[0]
        if cmd == CMD_BYE || cmd == CMD_EXIT || cmd == CMD_QUIT {
            this.bye()
        }

        if cmd == CMD_HELP {
            var method string
            if len(inputs) == 2 {
                method = "Help_" + inputs[1]
            } else {
                method = "Help"
            }
            m := reflect.ValueOf(this.child).MethodByName(method)
            if m == nil {
                println("shit")
            }
            m.Call([]reflect.Value{})
        }

    }
}
