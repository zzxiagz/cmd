package cmd

type Cmd struct {
    Prompt, Intro string
    child interface{}
}
