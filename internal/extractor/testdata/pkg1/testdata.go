package pkg1

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

const (
	TestVal24 Test2Enum = 3
)

type Test3Enum Test1Enum

const (
	TestVal31 Test3Enum = 0
	TestVal32 Test3Enum = 1
	TestVal33 Test3Enum = 2
)
