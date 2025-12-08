// Package rubyast provides utilities for working with Ruby AST nodes.
package rubyast

// Tree-sitter Ruby node types.
const (
	NodeCall        = "call"
	NodeMethodCall  = "method_call"
	NodeIdentifier  = "identifier"
	NodeString      = "string"
	NodeSymbol      = "symbol"
	NodeSimpleSymbol = "simple_symbol"
	NodeBlock       = "block"
	NodeDoBlock     = "do_block"
	NodeArguments   = "argument_list"
	NodeProgram     = "program"
)
