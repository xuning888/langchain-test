package tool

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"go.starlark.net/lib/math"
	"go.starlark.net/starlark"
	"log"
)

type MCalculator struct {
	CallbacksHandler callbacks.Handler
}

func (m *MCalculator) Name() string {
	return "mCalculator"
}
func (m *MCalculator) Description() string {
	return `Useful for getting the result of a math expression. 
	The input to this tool should be a valid mathematical expression that could be executed by a starlark evaluator.`
}
func (m *MCalculator) Call(ctx context.Context, input string) (string, error) {

	log.Println(fmt.Sprintf("llm call MCalculator input: %s", input))

	if m.CallbacksHandler != nil {
		m.CallbacksHandler.HandleToolStart(ctx, input)
	}

	v, err := starlark.Eval(&starlark.Thread{Name: "main"}, "input", input, math.Module.Members)

	if err != nil {
		log.Printf("starlark Eval error: %v\n", err)
		return fmt.Sprintf("error from evaluator: %s", err.Error()), nil //nolint:nilerr
	}

	result := v.String()
	log.Println(fmt.Sprintf("starlark Eval result: %v", result))

	if m.CallbacksHandler != nil {
		m.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}
