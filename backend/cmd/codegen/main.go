// Package main generates TypeScript types and constants from Go code
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Parse the errors.go file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "../../pkg/errors/errors.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	var errorCodes []struct {
		Name  string
		Value string
	}

	// Extract error code constants
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if strings.HasPrefix(name.Name, "Err") && len(valueSpec.Values) > i {
							if basicLit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
								value := strings.Trim(basicLit.Value, `"`)
								errorCodes = append(errorCodes, struct {
									Name  string
									Value string
								}{
									Name:  name.Name,
									Value: value,
								})
							}
						}
					}
				}
			}
		}
		return true
	})

	// Generate TypeScript file
	generateTypeScript(errorCodes)

	fmt.Printf("✓ Generated %d error codes to frontend/src/lib/errors-generated.ts\n", len(errorCodes))
}

func generateTypeScript(codes []struct{ Name, Value string }) {
	var b strings.Builder

	b.WriteString("// AUTO-GENERATED - DO NOT EDIT\n")
	b.WriteString(fmt.Sprintf("// Generated at %s\n", time.Now().Format(time.RFC3339)))
	b.WriteString("// Source: backend/pkg/errors/errors.go\n\n")

	b.WriteString("/**\n")
	b.WriteString(" * Error codes matching backend error definitions\n")
	b.WriteString(" * These are auto-generated from Go error constants\n")
	b.WriteString(" */\n")
	b.WriteString("export const ErrorCodes = {\n")

	for _, code := range codes {
		// Convert Go constant name to TypeScript constant
		// ErrAuthInvalidCredentials -> AUTH_INVALID_CREDENTIALS
		tsName := strings.TrimPrefix(code.Name, "Err")
		b.WriteString(fmt.Sprintf("  %s: '%s',\n", tsName, code.Value))
	}

	b.WriteString("} as const;\n\n")

	b.WriteString("export type ErrorCode = typeof ErrorCodes[keyof typeof ErrorCodes];\n\n")

	b.WriteString("/**\n")
	b.WriteString(" * Check if a string is a valid error code\n")
	b.WriteString(" */\n")
	b.WriteString("export function isValidErrorCode(code: string): code is ErrorCode {\n")
	b.WriteString("  return Object.values(ErrorCodes).includes(code as ErrorCode);\n")
	b.WriteString("}\n")

	// Resolve output path and ensure it stays within the project (no path traversal)
	outputPath := filepath.Join("..", "..", "..", "frontend", "src", "lib", "errors-generated.ts")
	clean := filepath.Clean(outputPath)
	if !strings.HasSuffix(filepath.ToSlash(clean), "frontend/src/lib/errors-generated.ts") {
		fmt.Fprintf(os.Stderr, "Error: resolved output path %q is outside expected location\n", clean)
		os.Exit(1)
	}
	if err := os.WriteFile(clean, []byte(b.String()), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}
