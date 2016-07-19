package fake

type FakeEnvironment struct {
	getBooleanStubMap map[string]getBoolReturnType
	getBooleanCallMap map[string]int

	getBooleanDefaultToTrueStubMap map[string]getBoolReturnType
	getBooleanDefaultToTrueCallMap map[string]int

	getStringStubMap map[string]string
	getStringCallMap map[string]int

	getIntegerStubMap map[string]getIntegerReturnType
	getIntegerCallMap map[string]int

	getBackendReturns getBackendReturnType
}

type getBoolReturnType struct {
	b bool
	e error
}

type getIntegerReturnType struct {
	i int
	e error
}

type getBackendReturnType struct {
	s string
	e error
}

func (fake *FakeEnvironment) GetBoolean(arg1 string) (bool, error) {
	if fake.getBooleanCallMap == nil {
		fake.getBooleanCallMap = map[string]int{}
	}
	fake.getBooleanCallMap[arg1] += 1

	stub := fake.getBooleanStubMap[arg1]
	return stub.b, stub.e
}

func (fake *FakeEnvironment) GetBooleanReturnsFor(varName string, returnBool bool, returnError error) {
	if fake.getBooleanStubMap == nil {
		fake.getBooleanStubMap = map[string]getBoolReturnType{}
	}
	fake.getBooleanStubMap[varName] = getBoolReturnType{b: returnBool, e: returnError}
}

func (fake *FakeEnvironment) GetBooleanCallCountFor(varName string) int {
	return fake.getBooleanCallMap[varName]
}

func (fake *FakeEnvironment) GetBooleanDefaultToTrue(arg1 string) (bool, error) {
	if fake.getBooleanDefaultToTrueCallMap == nil {
		fake.getBooleanDefaultToTrueCallMap = map[string]int{}
	}
	fake.getBooleanDefaultToTrueCallMap[arg1] += 1

	stub := fake.getBooleanDefaultToTrueStubMap[arg1]
	return stub.b, stub.e
}

func (fake *FakeEnvironment) GetBooleanDefaultToTrueReturnsFor(varName string, returnBool bool, returnError error) {
	if fake.getBooleanDefaultToTrueStubMap == nil {
		fake.getBooleanDefaultToTrueStubMap = map[string]getBoolReturnType{}
	}
	fake.getBooleanDefaultToTrueStubMap[varName] = getBoolReturnType{b: returnBool, e: returnError}
}

func (fake *FakeEnvironment) GetBooleanDefaultToTrueCallCountFor(varName string) int {
	return fake.getBooleanDefaultToTrueCallMap[varName]
}

func (fake *FakeEnvironment) GetString(arg1 string) string {
	if fake.getStringCallMap == nil {
		fake.getStringCallMap = map[string]int{}
	}
	fake.getStringCallMap[arg1] += 1

	stub := fake.getStringStubMap[arg1]
	return stub
}

func (fake *FakeEnvironment) GetStringReturnsFor(varName string, returnString string) {
	if fake.getStringStubMap == nil {
		fake.getStringStubMap = map[string]string{}
	}
	fake.getStringStubMap[varName] = returnString
}

func (fake *FakeEnvironment) GetStringCallCountFor(varName string) int {
	return fake.getStringCallMap[varName]
}

func (fake *FakeEnvironment) GetInteger(arg1 string) (int, error) {
	if fake.getIntegerCallMap == nil {
		fake.getIntegerCallMap = map[string]int{}
	}
	fake.getIntegerCallMap[arg1] += 1

	stub := fake.getIntegerStubMap[arg1]
	return stub.i, stub.e
}

func (fake *FakeEnvironment) GetIntegerReturnsFor(varName string, returnInt int, returnError error) {
	if fake.getIntegerStubMap == nil {
		fake.getIntegerStubMap = map[string]getIntegerReturnType{}
	}
	fake.getIntegerStubMap[varName] = getIntegerReturnType{i: returnInt, e: returnError}
}

func (fake *FakeEnvironment) GetIntegerCallCountFor(varName string) int {
	return fake.getIntegerCallMap[varName]
}

func (fake *FakeEnvironment) GetBackend() (string, error) {
	stub := fake.getBackendReturns
	return stub.s, stub.e
}

func (fake *FakeEnvironment) GetBackendReturns(s string, e error) {
	fake.getBackendReturns = getBackendReturnType{s: s, e: e}
}
