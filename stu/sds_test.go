package stu

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//test behavior, not the implement

const one_M = 1024 * 1024

func TestCompare(t *testing.T) {
	Convey("compare true", t, func() {
		one := "adfa4tadf"
		this := NewSdsString(one)
		that := NewSdsString(one)
		So(this.Compare(that), ShouldBeTrue)
	})
	Convey("compare length not equal", t, func() {
		one := "adsf"
		two := "adsffdsf"
		So(NewSdsString(one).Compare(NewSdsString(two)), ShouldBeFalse)
	})
}

func TestNewSds(t *testing.T) {
	Convey("new sds", t, func() {
		content := []byte("12")
		sd := NewSds([]byte("12"))
		So(sd.Len(), ShouldEqual, len(content))
	})
}

func TestNewSdsString(t *testing.T) {
	Convey("new sds string", t, func() {
		content := "中文"
		sd := NewSdsString(content)
		So(sd.Len(), ShouldEqual, len([]byte(content)))
	})
}

func TestDup(t *testing.T) {
	Convey("sds dup", t, func() {
		content := "dup"
		one := NewSdsString(content)
		two := one.Dup()
		So(one.Compare(two), ShouldBeTrue)
	})
}

func TestClear(t *testing.T) {
	Convey("sds clear", t, func() {
		content := "clear"
		sd := NewSdsString(content)
		sd.Clear()
		So(sd.Len(), ShouldEqual, 0)
	})
}

func TestCat(t *testing.T) {
	Convey("cat", t, func() {
		origin := "origin"
		tail := "I am a tail of cat"
		sd := NewSdsString(origin)
		err := sd.Cat([]byte(tail))
		So(err, ShouldBeNil)
		So(sd.Len(), ShouldEqual, (len([]byte(origin)) + len([]byte(tail))))
	})
	Convey("too big to cat", t, func() {
		origin := "origin"
		tail := bytes.Repeat(make([]byte, 1), MAX_LENGTH-len([]byte(origin))+1)
		sd := NewSdsString(origin)
		err := sd.Cat(tail)
		So(err, ShouldEqual, ErrLengthTooBig)
	})
}

func TestRange(t *testing.T) {
	Convey("range", t, func() {
		l := 100
		origin := bytes.Repeat(make([]byte, 1), l)
		start := int(math.Floor(float64(l) * rand.Float64()))
		end := start + int(math.Floor(float64(l-start)*rand.Float64()))
		sd := NewSds(origin)
		err := sd.Range(start, end)
		So(err, ShouldBeNil)
		So(sd.Len(), ShouldEqual, end-start)
	})
}

func TestTrim(t *testing.T) {
	Convey("Trim", t, func() {
		toTrim := "abc"
		left := "aadfa4erttwhuyhknlk5fad57fads"
		one := fmt.Sprintf("%s%s%s", toTrim, left, toTrim)
		sdLeft := NewSdsString(left)
		sd := NewSdsString(one)
		err := sd.Trim([]byte(toTrim))
		So(err, ShouldBeNil)
		So(sd.Compare(sdLeft), ShouldBeTrue)
	})
	Convey("trim head first", t, func() {
		toTrim := "fef"
		content := "fefef"
		sd := NewSdsString(content)
		err := sd.Trim([]byte(toTrim))
		So(err, ShouldBeNil)
		So(sd.Compare(NewSdsString("ef")), ShouldBeTrue)
	})
}
