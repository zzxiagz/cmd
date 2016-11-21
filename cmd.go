package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/flynn/go-shlex"
	"github.com/peterh/liner"
)

type Cmd struct {
	liner  *liner.State
	Prompt string
	Intro  string
	client interface{}
}

var (
	CMD_EXITS []string
	CMD_HELPS []string
	CMD_LISTS []string
)

const (
	EOL       = '\n'
	EMPTY_ERR = "empty"
	QUIT_ERR  = "quit"

	SPACE_STR   = " "
	EMPTY_STR   = ""
	EXEC_PREFIX = "Do_"
	HELP_PREFIX = "Help_"

	DEFAULT_PROMPT = ">> "
)

func init() {
	CMD_EXITS = []string{"quit", "q"}
	CMD_HELPS = []string{"help", "h"}
	CMD_LISTS = []string{"list", "ls", "l", "?"}
}

func (cmd *Cmd) initLiner() {
	cmd.liner = liner.NewLiner()
}

func (cmd *Cmd) initCompleter() {
	t := reflect.TypeOf(cmd.client)
	cmdList := make([]string, 0)
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if strings.HasPrefix(methodName, EXEC_PREFIX) {
			cmdList = append(cmdList, methodName[len(EXEC_PREFIX):])
		}
	}

	cmd.liner.SetCompleter(func(line string) (c []string) {
		for _, n := range cmdList {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
}

func New(client interface{}) *Cmd {
	cmd := &Cmd{Prompt: DEFAULT_PROMPT, client: client}
	cmd.initLiner()
	cmd.initCompleter()

	return cmd
}

func (this *Cmd) intro() {
	if this.Intro != "" {
		fmt.Println(this.Intro)
	}
}

func (this *Cmd) isCommand(cmd string, cmdList []string) bool {
	for _, c := range cmdList {
		if c == cmd {
			return true
		}
	}
	return false
}

func (this *Cmd) isExit(cmd string) bool {
	return this.isCommand(cmd, CMD_EXITS)
}

func (this *Cmd) isHelp(cmd string) bool {
	return this.isCommand(cmd, CMD_HELPS)
}

func (this *Cmd) isList(cmd string) bool {
	return this.isCommand(cmd, CMD_LISTS)
}

func (this *Cmd) parseLine() (cmd string, args []string, err error) {
	rawInput, err := this.liner.Prompt(this.Prompt)
	if err != nil {
		return
	}

	input := rawInput
	if input == EMPTY_STR || strings.TrimSpace(input) == EMPTY_STR {
		err = errors.New(EMPTY_ERR)
		return
	}

	inputs, err := shlex.Split(input)
	if len(inputs) > 1 {
		args = make([]string, 0)
		for _, in := range inputs[1:] {
			x := strings.TrimSpace(in)
			if x != "" {
				args = append(args, x)
			}
		}
	}
	cmd = inputs[0]
	this.liner.AppendHistory(rawInput)
	return
}

func (this *Cmd) Cmdloop() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Cmd loop quit:", r)
		}
		this.liner.Close()
	}()

	this.intro()

	for {
		cmd, args, err := this.parseLine()
		if err != nil {
			switch err.Error() {
			case EMPTY_ERR:
				continue
			default:
				this.liner.Close()
				return
			}
		}

		if this.isExit(cmd) {
			return
		}

		method := ""
		if this.isHelp(cmd) {
			if len(args) >= 1 {
				method = HELP_PREFIX + args[0]
				args = args[1:]
			} else {
				continue
			}
		} else if this.isList(cmd) {
			this.listCommands(this.client)
			continue
		} else {
			method = EXEC_PREFIX + cmd
		}

		this.tryInvoke(this.client, method, args)
	}
}

func (this *Cmd) notFound(method string) {
	fmt.Printf("Invalid command: %s\n", method)
}

func (this *Cmd) listCommands(i interface{}) {
	fmt.Println("Available commands:")

	t := reflect.TypeOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		methodName := t.Method(i).Name
		if strings.HasPrefix(methodName, EXEC_PREFIX) {
			fmt.Printf("%s ", methodName[len(EXEC_PREFIX):])
		}
	}
	println()
}

func (this *Cmd) tryInvoke(i interface{}, methodName string, args []string) {
	method := reflect.ValueOf(i).MethodByName(methodName)
	if !method.IsValid() {
		this.notFound(methodName)
		return
	}

	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	method.Call(params)
}
