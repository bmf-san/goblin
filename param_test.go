package goblin

import (
	"reflect"
	"testing"
)

func TestGetParamsAndPutParams(t *testing.T) {
	params := &[]Param{
		{
			key:   "id",
			value: "123",
		},
		{
			key:   "name",
			value: "john",
		},
	}

	tree := newTree()
	tree.paramsPool.New = func() interface{} {
		// NOTE: It is better to set the maximum value of paramters to capacity.
		return &[]Param{}
	}
	params = tree.getParams()
	tree.putParams(params)

	expectedParams := &[]Param{}
	actualParams := tree.getParams()
	if !reflect.DeepEqual(actualParams, expectedParams) {
		t.Errorf("actual:%v expected:%v", actualParams, expectedParams)
	}
}
