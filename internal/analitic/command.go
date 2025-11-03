package analitic

type Command interface {
	Execute() (string, error)
}

type Invoker struct {
	commands []Command
}

func (i *Invoker) AddCommand(cmd Command) {
	i.commands = append(i.commands, cmd)
}

func (i *Invoker) Run() (string, error) {
	var lastResult string
	var lastError error

	for _, cmd := range i.commands {
		result, err := cmd.Execute()
		lastResult = result
		lastError = err

		if err != nil {
			return result, err
		}
	}

	return lastResult, lastError
}
func (i *Invoker) Clear() {
	i.commands = []Command{}
}
