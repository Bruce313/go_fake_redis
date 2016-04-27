package main

type command interface {
	exec(cli *client) (string, error)
	parse(...string) bool
}

var (
	commandGet = &getCommand{}
)

type getCommand struct {
	key string
}

func (gc *getCommand) exec(cli *client) (string, error) {
	v := cli.db.get(gc.key)
	if v == nil {
		return "no value", nil	
	}
	return string(v.result()), nil
}

func (gc *getCommand) parse(argvs ...string) bool {
	if len(argvs) == 1 {
		gc.key = argvs[0]	
		return true
	}
	return false
}

func lookupCommand(name string, argv ...string) (command, error) {
	return commandGet, nil
}
