package decl

import "github.com/Knetic/govaluate"

// CalculateExpress 使用 govaluate 计算公式
func CalculateExpress(expression string, variables map[string]any) (any, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(variables)
	if err != nil {
		return 0, err
	}

	return result, nil
}
