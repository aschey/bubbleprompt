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
	}
	return parent
}

func getString(value goja.Value) string {
	if value.ExportType().String() == "string" {
		return `"` + value.String() + `"`
	}
	return value.String()
}

func (m completerModel) evaluateExpression(parent *goja.Object, expression expression) goja.Value {
	var value goja.Value = nil
	switch {
	case expression.PropAccessor != nil:
		value = m.evalutePropAccessor(parent, *expression.PropAccessor)
	case expression.Token != nil:
		value = m.evaluateToken(parent, *expression.Token)
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

func (m completerModel) evaluateToken(parent *goja.Object, token token) goja.Value {
	switch {
	case token.Literal != nil:
		return m.evaluateLiteral(*token.Literal)
	case token.Variable != nil:
		return parent.Get(*token.Variable)
	}
	return nil
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
	return m.evaluteAccessor(curVal.ToObject(m.vm), propAccessor.Accessor)
}

func (m completerModel) evaluteAccessor(parent *goja.Object, accessor accessor) goja.Value {
	var value goja.Value
	switch {
	case accessor.Indexer != nil:
		if accessor.Indexer.Expression == nil {
			return parent
		}
		value = m.evaluateIndexer(parent, *accessor.Indexer)
	case accessor.Delim != nil:
		if accessor.Prop == nil {
			return parent
		}
		value = parent.Get(*accessor.Prop)
		if accessor.Accessor == nil {
			return parent
		}
	}
	if accessor.Accessor != nil {
		return m.evaluteAccessor(value.ToObject(m.vm), *accessor.Accessor)
	}

	return value
}

func (m completerModel) evaluateIndexer(parent *goja.Object, indexer indexer) goja.Value {
	val := m.evaluateExpression(parent, *indexer.Expression)
	return parent.Get(val.String())
}
