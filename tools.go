package main

import "github.com/ollama/ollama/api"

var tools = []api.Tool{
    {
        Type: "function",
        Function: api.ToolFunction{
            Name:        "get_time",
            Description: "Returns current time",
            Parameters:  api.ToolFunctionParameters{
				Type: "object",
				Properties: &api.ToolPropertiesMap{},
				Required:   []string{},
			},
        },
    },
	 {
        Type: "function",
        Function: api.ToolFunction{
            Name:        "get_date",
            Description: "Returns current date",
            Parameters:  api.ToolFunctionParameters{
				Type: "object",
				Properties: &api.ToolPropertiesMap{},
				Required:   []string{},
			},
        },
    },
}
