package test

import "caseGenerator/parser/expression/z3/test/testpackage"

// 这里都是ident的不同场景

var PackageGlobalString = "hello"
var PackageGlobalBool = true
var PackageGlobalNumber = 0
var PackageGlobalReferenceNumber = PackageGlobalNumber
var PackageGlobalReceiver = GlobalStruct{Name: "Name"}
var PackageGlobalReceiverName = PackageGlobalReceiver.Name
var PackageGlobalReceiverAge = PackageGlobalReceiver.Age
var PackageGlobalReceiverIsPass = PackageGlobalReceiver.IsPass
var PackageGlobalCallName = getGlobalStruct().Name // 有可能是call.属性
var PackageGlobalFunction = GetIdent()             // 有可能是call
var PackageGlobalSelectorNumber = testpackage.PACKAGE_INT
var PackageGlobalSelectorStruct = testpackage.GlobalStruct{GlobalName: "GlobalName"}

type GlobalStruct struct {
	Name   string
	Age    int
	IsPass bool
}

func getGlobalStruct() GlobalStruct {
	return GlobalStruct{}
}

func GetIdent() int {
	return 1
}
