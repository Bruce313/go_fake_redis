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

//print sds
func (self *Sds) String() string {
	m := 20
	if self.length < 20 {
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
		free:   l,
	}
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
func (self *Sds) Cat(tail []byte) {
	//if < 1M after
}

//if < 1M after grow, grow to 2 * len, otherwise grow 1M
func (self *Sds) Grow() error {
	size := self.length
	if size*2 > (10 ^ 6) {
		size = 10 ^ 6
	}
	return self.grow(size)
}
func (self *Sds) GrowTo(l int) {

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
	nb := self.buf[start:end]
	copy(self.buf, nb)
	self.length = end - start
	return nil
}

//trim delim from begin and end
func (self *Sds) Trim(delim []byte) {
	//	head := true
	//	tail := true
	//	for i, v := range delim {
	//		head = (v == self.buf[i])
	//		//tail = (v == )
	//	}
}

//determine two sds equal or not
func (self *Sds) Compare(other *Sds) bool {
	l := self.Len()
	if self.Len() > other.Len() {
		l = other.Len()
	}
	for i := 0; i < l; i++ {
		if self.buf[i] != other.buf[i] {
			return false
		}
	}
	return true
}
