package test

// RequestZeroInt 测试方法请求的ident
func RequestZeroInt(ZeroInt string) int {
	// 同时有request和全局变量时， 取得是request的内容
	// 入参是string，全局变量是int类型
	localVariable := ZeroInt
	if localVariable == "" {
		return 0
	}
	return -1
}

// GlobalBool 局部布尔变量
func GlobalBool() int {
	localVariable := PackageGlobalBool
	if localVariable {
		return 0
	}
	return 1
}

// GlobalNumber 局部int变量
func GlobalNumber() int {
	localVariable := PackageGlobalNumber
	return localVariable
}

// GlobalReferenceNumber 局部变量 是引用的其他局部变量
func GlobalReferenceNumber() int {
	localVariable := PackageGlobalReferenceNumber
	return localVariable
}

// GlobalString 局部string变量
func GlobalString() int {
	localVariable := PackageGlobalString
	if localVariable == "" {
		return 0
	}
	return 1
}

// GlobalReceiver 局部struct变量
func GlobalReceiver() int {
	localVariable := PackageGlobalReceiver
	if localVariable.Name != "" {
		return 1
	}
	return -1
}

// GlobalReceiverName 局部struct变量
func GlobalReceiverName() int {
	localVariable := PackageGlobalReceiverName
	if localVariable != "" {
		return 1
	}
	return -1
}

// GlobalReceiverAge 局部struct.int变量
func GlobalReceiverAge() int {
	localVariable := PackageGlobalReceiverAge
	return localVariable
}

// GlobalReceiverIsPass 局部struct.bool变量
func GlobalReceiverIsPass() int {
	localVariable := PackageGlobalReceiverIsPass
	if localVariable {
		return 0
	}
	return -1
}

// GlobalCallReceiverString 局部func()->struct.string`变量
func GlobalCallReceiverString() int {
	localVariable := PackageGlobalCallName
	if localVariable == "" {
		return 0
	}
	return -1
}

// GlobalCallInt 局部func()->int`变量
func GlobalCallInt() int {
	localVariable := PackageGlobalFunction
	return localVariable
}

// GlobalSelectorInt 局部xx.int`变量
func GlobalSelectorInt() int {
	localVariable := PackageGlobalSelectorNumber
	return localVariable
}

// GlobalSelectorStruct 局部xx.int`变量
func GlobalSelectorStruct() int {
	localVariable := PackageGlobalSelectorStruct
	if localVariable.GlobalName == "" {
		return 1
	}
	return -1
}
