package analitic

type Command interface {
	Execute() error
}

type Invoker struct {
	commands []Command
}

func (i *Invoker) AddCommand(cmd Command) {
	i.commands = append(i.commands, cmd)
}

func (i *Invoker) Run() {
	for _, cmd := range i.commands {
		cmd.Execute()
	}
}
func (i *Invoker) Clear() {
	i.commands = []Command{}
}
