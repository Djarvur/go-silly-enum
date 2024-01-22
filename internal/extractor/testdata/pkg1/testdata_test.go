package pkg1_test

import (
	"testing"

	"github.com/Djarvur/go-silly-enum/internal/extractor/testdata/pkg1"
)

const (
	TestVal01 = iota
	TestVal02
	TestVal03
)

type Test1Enum uint8

const (
	TestVal11 Test1Enum = iota
	TestVal12 Test1Enum = iota
	TestVal13 Test1Enum = iota
)

type Test2Enum uint8

const (
	TestVal21 Test2Enum = 0
	TestVal22 Test2Enum = 1
	TestVal23 Test2Enum = 2
)

func TestTest1Enum(t *testing.T) {
	t.Logf("%T %v %#+v", pkg1.TestVal11, pkg1.TestVal11, pkg1.TestVal11)
	t.Logf("%T %v %#+v", pkg1.TestVal12, pkg1.TestVal12, pkg1.TestVal12)
	t.Logf("%T %v %#+v", pkg1.TestVal13, pkg1.TestVal13, pkg1.TestVal13)
	t.Logf("%T %v %#+v", pkg1.TestVal16, pkg1.TestVal16, pkg1.TestVal16)
	t.Logf("%T %v %#+v", pkg1.TestVal15, pkg1.TestVal15, pkg1.TestVal15)
}
