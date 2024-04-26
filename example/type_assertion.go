package example

import (
	"caseGenerator/example/dict"
	"context"
	"errors"
	"time"
)

// TypeAssertionAnyTest 为了方便测试，这里罗列多条类型断言
func TypeAssertionAnyTest(req any) error {
	_, ok := req.(string)
	if ok {
		return errors.New("类型断言是string")
	}
	_, ok = req.(bool)
	if ok {
		return errors.New("类型断言是bool")
	}
	_, ok = req.(map[string]string)
	if ok {
		return errors.New("类型断言是map[string]string")
	}
	_, ok = req.(map[Example]dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[Example]dict.ExampleDict")
	}
	_, ok = req.(map[*Example]*dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[*Example]*dict.ExampleDict")
	}
	_, ok = req.(map[context.Context][]string)
	if ok {
		return errors.New("类型断言是map[context.Context][]string")
	}
	_, ok = req.(map[string][][][][][]*Example)
	if ok {
		return errors.New("类型断言是 map[string][][][][][]*Example")
	}
	_, ok = req.(map[string]map[*Example]map[context.Context]map[time.Time]bool)
	if ok {
		return errors.New("类型断言是 map[string]map[*Example]map[context.Context]map[time.Time]bool")
	}
	return nil
}

// TypeMultipleAssertionAnyTest 为了方便测试，这里罗列多条类型断言， 排列组合
func TypeMultipleAssertionAnyTest(req1 any, req2 any) error {
	_, ok := req1.(string)
	if ok {
		return errors.New("类型断言是string")
	}
	_, ok = req1.(bool)
	if ok {
		return errors.New("类型断言是bool")
	}
	_, ok = req1.(map[string]string)
	if ok {
		return errors.New("类型断言是map[string]string")
	}
	_, ok = req1.(map[Example]dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[Example]dict.ExampleDict")
	}
	_, ok = req1.(map[*Example]*dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[*Example]*dict.ExampleDict")
	}
	_, ok = req1.(map[context.Context][]string)
	if ok {
		return errors.New("类型断言是map[context.Context][]string")
	}
	_, ok = req1.(map[string][][][][][]*Example)
	if ok {
		return errors.New("类型断言是 map[string][][][][][]*Example")
	}
	_, ok = req1.(map[string]map[*Example]map[context.Context]map[time.Time]bool)
	if ok {
		return errors.New("类型断言是 map[string]map[*Example]map[context.Context]map[time.Time]bool")
	}

	_, ok = req2.(string)
	if ok {
		return errors.New("类型断言是string")
	}
	_, ok = req2.(bool)
	if ok {
		return errors.New("类型断言是bool")
	}
	_, ok = req2.(map[string]string)
	if ok {
		return errors.New("类型断言是map[string]string")
	}
	_, ok = req2.(map[Example]dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[Example]dict.ExampleDict")
	}
	_, ok = req2.(map[*Example]*dict.ExampleDict)
	if ok {
		return errors.New("类型断言是map[*Example]*dict.ExampleDict")
	}
	_, ok = req2.(map[context.Context][]string)
	if ok {
		return errors.New("类型断言是map[context.Context][]string")
	}
	_, ok = req2.(map[string][][][][][]*Example)
	if ok {
		return errors.New("类型断言是 map[string][][][][][]*Example")
	}
	_, ok = req2.(map[string]map[*Example]map[context.Context]map[time.Time]bool)
	if ok {
		return errors.New("类型断言是 map[string]map[*Example]map[context.Context]map[time.Time]bool")
	}
	return nil
}
