package stu

import (
	"testing"
)

func TestNewSds(t *testing.T) {
	sd := NewSds([]byte("12"))
	t.Logf("sds:%s", sd)
}
