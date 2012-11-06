package main

import "../cmd"

type MyCmd struct {
    cmd.Cmd
}

func main() {
    my := MyCmd{}
    cmd := cmd.New(my)
    cmd.Cmdloop()
}
