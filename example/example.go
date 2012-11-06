package main

import (
    "../"
)

type MyCmd struct {
}

func (this MyCmd) Help_go() {
    println("haha")
}

func (this MyCmd) Help() {
    println("any help here")
}

func main() {
    my := new(MyCmd)
    cmd := cmd.New(my)
    cmd.Cmdloop()
}
