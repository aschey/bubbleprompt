package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dop251/goja"
)

func (m completerModel) evaluateStatement(statement statement) goja.Value {
	parent := m.vm.GlobalObject()
	switch {
	case statement.Expression != nil:
		if statement.Expression.Token != nil {
			return parent
		}
		return m.evaluateExpression(parent, *statement.Expression)
	case statement.Assignment != nil:
		return m.evaluateAssignment(parent, *statement.Assignment)
	}
	return parent
}

func getString(value goja.Value) string {
	if value.ExportType().String() == "string" {
		return `"` + value.String() + `"`
	}
	return value.String()
}

func (m completerModel) evaluateAssignment(parent *goja.Object, assignment assignment) goja.Value {
	if assignment.Expression != nil {
		return m.evaluateExpression(parent, *assignment.Expression)
	}
	return goja.Null()
}

func (m completerModel) evaluateExpression(parent *goja.Object, expression expression) goja.Value {
	var value goja.Value = nil
	switch {
	case expression.PropAccessor != nil:
		value = m.evalutePropAccessor(parent, *expression.PropAccessor)
	case expression.Token != nil:
		value = m.evaluateToken(parent, *expression.Token)
	case expression.Object != nil:
		value = m.evaluteObject(parent, *expression.Object)
	}

	if expression.InfixOp != nil && expression.Expression != nil {
		rightSide := m.evaluateExpression(parent, *expression.Expression)
		val, err := m.vm.RunString(getString(value) + expression.InfixOp.Op + getString(rightSide))
		if err != nil {
			fmt.Println(err)
		}

		return val
	}
	return value
}

func (m completerModel) evaluteObject(parent *goja.Object, object object) goja.Value {
	if object.Properties != nil && len(*object.Properties) > 0 {
		props := *object.Properties
		last := props[len(props)-1]
		return m.evaluteKeyValuePair(parent, last)
	}
	return goja.Null()
}

func (m completerModel) evaluteKeyValuePair(parent *goja.Object, keyValuePair keyValuePair) goja.Value {
	if keyValuePair.Delim != nil {
		if keyValuePair.Value == nil || keyValuePair.Value.Token != nil {
			return parent
		}

		return m.evaluateExpression(parent, *keyValuePair.Value)
	}
	return goja.Null()
}

func (m completerModel) evaluateToken(parent *goja.Object, token token) goja.Value {
	switch {
	case token.Literal != nil:
		return m.evaluateLiteral(*token.Literal)
	case token.Variable != nil:
		return parent.Get(*token.Variable)
	}
	return goja.Null()
}

func (m completerModel) evaluateLiteral(literal literal) goja.Value {
	literalVal := ""
	switch {
	case literal.Str != nil:
		literalVal = *literal.Str
	case literal.Boolean != nil:
		literalVal = strconv.FormatBool(*literal.Boolean)
	case literal.Number != nil:
		num := *literal.Number
		// Check if number has decimal
		if num-math.Floor(num) > 0 {
			literalVal = strconv.FormatFloat(num, 'f', 4, 64)
		} else {
			literalVal = strconv.FormatInt(int64(*literal.Number), 10)
		}
	}
	val, _ := m.vm.RunString(literalVal)
	return val
}

func (m completerModel) evalutePropAccessor(parent *goja.Object, propAccessor propAccessor) goja.Value {
	curVal := parent.Get(propAccessor.Identifier)
	return m.evaluateAccessor(m.vm.ToObject(curVal), propAccessor.Accessor)
}

func (m completerModel) evaluateAccessor(parent *goja.Object, accessor accessor) goja.Value {
	var value goja.Value
	switch {
	case accessor.Indexer != nil:
		if accessor.Indexer.Expression == nil {
			return parent
		}
		value = m.evaluateIndexer(parent, *accessor.Indexer)
	case accessor.Delim != nil:
		if accessor.Prop == nil {
			// Invalid case: two delimiters with no value, don't return suggestions
			if accessor.Accessor != nil {
				return goja.Null()
			}
			return parent
		}
		value = parent.Get(*accessor.Prop)
		if accessor.Accessor == nil {
			return parent
		}
	}
	if accessor.Accessor != nil {
		return m.evaluateAccessor(m.vm.ToObject(value), *accessor.Accessor)
	}

	return value
}

func (m completerModel) evaluateIndexer(parent *goja.Object, indexer indexer) goja.Value {
	val := m.evaluateExpression(parent, *indexer.Expression)
	return parent.Get(val.String())
}
