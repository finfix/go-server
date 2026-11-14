package mcp

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"
)

// schemaOptions — суммы decimal в JSON — числа полной точности (MarshalJSONWithoutQuotes), а не
// объекты, как вывела бы схема из устройства структуры
var schemaOptions = &jsonschema.ForOptions{ //nolint:exhaustruct
	TypeSchemas: map[reflect.Type]*jsonschema.Schema{
		reflect.TypeFor[decimal.Decimal](): {Type: "number"}, //nolint:exhaustruct
	},
}

// AddTool регистрирует инструмент со схемами аргументов и ответа, в которых суммы — числа
func AddTool[In, Out any](server *sdk.Server, tool *sdk.Tool, handler sdk.ToolHandlerFor[In, Out]) {
	inputSchema, err := jsonschema.For[In](schemaOptions)
	if err != nil {
		panic(err)
	}
	outputSchema, err := jsonschema.For[Out](schemaOptions)
	if err != nil {
		panic(err)
	}
	tool.InputSchema, tool.OutputSchema = inputSchema, outputSchema
	sdk.AddTool(server, tool, handler)
}
