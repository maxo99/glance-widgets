package glance

import (
	"html/template"

	"github.com/tidwall/gjson"
)

// GetTemplateFunctions returns the combined template functions from Glance
func GetTemplateFunctions() template.FuncMap {
	funcs := make(template.FuncMap)

	// Add global template functions
	for key, value := range globalTemplateFunctions {
		funcs[key] = value
	}

	// Add custom API template functions
	customFuncs := customAPITemplateFuncs
	for key, value := range customFuncs {
		funcs[key] = value
	}

	return funcs
}

// DecoratedGJSONResult wraps gjson.Result for template use
type DecoratedGJSONResult = decoratedGJSONResult

// NewDecoratedGJSONResult creates a new decorated GJSON result
func NewDecoratedGJSONResult(result gjson.Result) *DecoratedGJSONResult {
	return &decoratedGJSONResult{Result: result}
}
