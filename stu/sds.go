package stu

import (
	"errors"
	"fmt"
)

type Sds struct {
	buf    []byte
	length int
	free   int
}

const show_LENGTH = 30

//print sds
func (self *Sds) String() string {
	m := show_LENGTH
	if self.length < show_LENGTH {
		m = self.length
	}
	return fmt.Sprintf("{\n\tlen:%d\n\tfree:%d\n\tcontent:%s\n}\n",
		self.length, self.free, self.buf[:m])
}

func NewSds(b []byte) *Sds {
	//todo: dive into slice impl of alloc
	l := len(b)
	nb := make([]byte, 2*l)
	copy(nb, b)
	return &Sds{
		buf:    nb,
		length: l,
		free:   0,
	}
}

func NewSdsString(s string) *Sds {
	return NewSds([]byte(s))
}

const default_EMPTY_LEN = 20

func NewSdsEmpty() *Sds {
	return &Sds{
		buf:    nil,
		length: 0,
		free:   default_EMPTY_LEN,
	}
}

func (self *Sds) Destory() {

}

func (self *Sds) Avail() int {
	return self.free
}

func (self *Sds) Len() int {
	return self.length
}

func (self *Sds) Dup() *Sds {
	b := make([]byte, self.length+self.free)
	copy(b, self.buf)
	return &Sds{
		buf:    b,
		length: self.length,
		free:   self.free,
	}
}

func (self *Sds) Clear() {
	self.free = self.free + self.length
	self.length = 0
}

//concat []byte to tail of sds
func (self *Sds) Cat(tail []byte) error {
	//if < 1M after
	err := self.growTo(self.Len() + len(tail))
	if err != nil {
		return err
	}
	copy(self.buf[self.Len():], tail)
	self.length = self.length + len(tail)
	return nil
}

const grow_BOUND = 1024 * 1024

func (self *Sds) growTo(l int) error {
	if l > MAX_LENGTH {
		return ErrLengthTooBig
	}
	var realLen = 2 * l
	if 2*l >= grow_BOUND {
		realLen = l + 1024*1024
	}
	return self.grow(realLen - self.length)
}

const MAX_LENGTH = 512 * 1024 * 1024

var ErrLengthTooBig = errors.New(fmt.Sprintf("max length of sds is:%d", MAX_LENGTH))

//expand count length
func (self *Sds) grow(count int) error {
	if self.length+count > MAX_LENGTH {
		return ErrLengthTooBig
	}
	self.buf = append(self.buf, make([]byte, count)...)
	self.free = self.free + count
	return nil
}

var (
	ErrOutofRange = errors.New("range out")
)

//clear data out of range
func (self *Sds) Range(start, end int) error {
	if start < 0 || end > self.length {
		return ErrOutofRange
	}
	self.buf = self.buf[start:end]
	self.free = self.free + (self.length - end + start)
	self.length = end - start
	return nil
}

//trim delim from begin and end
func (self *Sds) Trim(delim []byte) error {
	lenDe := len(delim)
	if lenDe > self.Len() {
		return nil
	}
	head := 0
	tail := self.length
	//find head
	trimHead := true
	for i := 0; i < lenDe; i++ {
		if delim[i] != self.buf[i] {
			trimHead = false
			break
		}
	}
	if trimHead {
		head = lenDe
	}
	if self.Len()-lenDe >= head {
		trimTail := true
		for i := lenDe - 1; i >= 0; i-- {
			if self.buf[self.Len()-lenDe+i] != delim[i] {
				trimTail = false
				break
			}
		}
		if trimTail {
			tail = self.Len() - lenDe
		}
	}
	return self.Range(head, tail)
}

//determine two sds equal or not
func (self *Sds) Compare(other *Sds) bool {
	if self.Len() != other.Len() {
		return false
	}
	l := other.Len()
	for i := 0; i < l; i++ {
		if self.buf[i] != other.buf[i] {
			return false
		}
	}
	return true
}

//
