package cmd

type Cmd struct {
    Prompt, Intro string
    client        interface{}
}
