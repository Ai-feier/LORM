package reflect

import (
	"reflect"
)

func ItrateFunc(entity any) (res map[string]FuncInfo, err error) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	res = make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func
		
		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn)
		inputValues := make([]reflect.Value, 0, numIn)

		// 第一个接受参数为方法调用类型
		input = append(input, reflect.TypeOf(entity))
		inputValues = append(inputValues, reflect.ValueOf(entity))
		
		for j := 1; j < numIn; j++ {
			fnIntype := fn.Type().In(j)
			input = append(input, fnIntype)
			inputValues = append(inputValues, reflect.Zero(fnIntype))
		}
		
		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))
		}

		resValues := fn.Call(inputValues)
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
			
	}
	return 
}

type FuncInfo struct {
	Name string
	InputTypes []reflect.Type
	OutputTypes []reflect.Type
	Result []any
}
