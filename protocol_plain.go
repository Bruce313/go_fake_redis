package main

import (
	"bufio"
	"io"
)

type protocol interface {
	Send(content string, w io.Writer) error
	Receive(r io.Reader) (string, error)
}

type protocolPlain struct {
}

func newProtocolPlain() *protocolPlain {
	return &protocolPlain{}
}

const DELIM = byte('\n')

func (self *protocolPlain) Send(content string, w io.Writer) error {
	en := append([]byte(content), DELIM)
	_, err := w.Write(en)
	return err
}

func (self *protocolPlain) Receive(r io.Reader) (string, error) {
	br := bufio.NewReader(r)
	line, err := br.ReadString(DELIM)
	return line, err
}
