package main

import (
    "../"
    "fmt"
    "strconv"
)

type MyCmd struct{}

func (this MyCmd) Help() {
    fmt.Println("Available commands:")
    fmt.Println("add go")
}

func (this MyCmd) Help_go() {
    fmt.Println("go name")
}

func (this MyCmd) Do_go(name string) {
    for _, r := range name {
        fmt.Println(r)
    }
}

func (this MyCmd) Help_add() {
    fmt.Println("add a b")
}

func (this MyCmd) Do_add(a, b string) {
    ai, _ := strconv.Atoi(a)
    bi, _ := strconv.Atoi(b)
    fmt.Printf("%s + %s = %d\n", a, b, ai+bi)
}

func main() {
    cmd := cmd.New(new(MyCmd))
    cmd.Prompt = "(ExampleOfCmd) "
    cmd.Intro = "这是个cmd包的使用例子"
    cmd.Cmdloop()
}
