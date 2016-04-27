package main

import (
	"fmt"
)

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
		return "(nil)", nil	
	}
	t := v.getType()
	if t != vtString {
		return fmt.Sprintf("can`t GET for type:%s", t), nil
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
