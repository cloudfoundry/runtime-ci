package configwriterfakes

import "github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"

type FakeEnvironment struct {
	GetBooleanStub    func(string) (bool, error)
	getBooleanStubMap map[string]getBoolReturnType
}

type getBoolReturnType struct {
	b bool
	e error
}

func (fake *FakeEnvironment) GetBoolean(arg1 string) (bool, error) {
	stub := fake.getBooleanStubMap[arg1]
	return stub.b, stub.e
}

func (fake *FakeEnvironment) GetBooleanReturnsFor(varName string, returnBool bool, returnError error) {
	if fake.getBooleanStubMap == nil {
		fake.getBooleanStubMap = map[string]getBoolReturnType{}
	}
	fake.getBooleanStubMap[varName] = getBoolReturnType{b: returnBool, e: returnError}
}

var _ configwriter.Environment = new(FakeEnvironment)
