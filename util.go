package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// This function basically looks at two functions defined in the provider.go, parses a map to get the list of functions that have schemas defined.
func readProviderForFuncNamesWithSchemaDefs(sourceDir string) map[string]bool {
	providerFileName := fmt.Sprintf("%s/%s", sourceDir, "provider.go")

	dataSourceFuncList := getDataSourceFuncNames(providerFileName, "dataSourcesMap")
	resourceFuncList := getResourceFuncNames(providerFileName, "resourcesMap")

	funcNames := make(map[string]bool)

	for k, v := range dataSourceFuncList {
		funcNames[k] = v
	}

	for k, v := range resourceFuncList {
		funcNames[k] = v
	}

	return funcNames
}

func getResourceFuncNames(fileName string, functionName string) map[string]bool {
	funcNames := make(map[string]bool)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok && function.Name.String() == functionName {
			funcMap := function.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts

			for _, mapEntry := range funcMap {
				resourceFuncName := mapEntry.(*ast.KeyValueExpr).Value.(*ast.CallExpr).Fun.(*ast.Ident).Name
				funcNames[resourceFuncName] = true
			}
		}
	}

	return funcNames
}

func getDataSourceFuncNames(fileName string, functionName string) map[string]bool {
	funcNames := make(map[string]bool)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok && function.Name.String() == functionName {
			funcMap := function.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts

			for _, mapEntry := range funcMap {
				dataSourceFuncName := mapEntry.(*ast.KeyValueExpr).Value.(*ast.CallExpr).Fun.(*ast.Ident).Name
				funcNames[dataSourceFuncName] = false
			}
		}
	}

	return funcNames
}

func findFuncNameWithSchemaDefsInAutoGenCode(fileName string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			if strings.HasSuffix(function.Name.String(), "Resource") || strings.HasSuffix(function.Name.String(), "Datasource") {
				return function.Name.String()
			}
		}
	}

	return ""
}

func isResourceFile(fileName string) bool {
	return strings.HasSuffix(fileName, "_resource.go")
}

func isDataSourceFile(fileName string) bool {
	return strings.HasSuffix(fileName, "_data_source.go")
}
