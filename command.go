package main

type command interface {
	exec(cli *client) (string, error)
}

var (
	commandGet = &getCommand{}
)

type getCommand struct {
	key string
}

func (self *getCommand) exec(cli *client) (string, error) {
	return "value from get command", nil
}

func (self *getCommand) setKey(k string) {
	self.key = k
}

func lookupCommand(name string, argv ...string) (command, error) {
	return commandGet, nil
}
