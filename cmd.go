package cmd

import (
    "os"
    "fmt"
)

func New(child interface{}) Cmd {
    this := Cmd{}
    this.child = child
    return this
}

func (this Cmd) Cmdloop() {
    for {
        os.Stdout.WriteString(this.Prompt)
        var in []byte
        os.Stdin.Read(in)
        fmt.Println(string(in))
    }
}
