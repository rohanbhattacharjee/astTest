package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func parseSchema(fileName string, funcNamesWithSchemaDefs map[string]bool) (schemaDefs []SchemaDef) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	//ast.Print(fset, f)

	// Uncomment the statement above to see the abstract syntax tree (AST) for the code.
	// Once you see the AST, the following code will make sense.
	// Logic: Looks for a function in the current file which will have a schema definition. Various correctness assumptions (in the input file) are assumed here.
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			if _, exists := funcNamesWithSchemaDefs[function.Name.String()]; exists {
				var elements = function.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts

				for _, elem := range elements {
					currElem := elem.(*ast.KeyValueExpr)
					if currElem.Key.(*ast.Ident).Name == "Schema" {
						schemaDefs = make([]SchemaDef, len(currElem.Value.(*ast.CompositeLit).Elts))

						for i, schemaElem := range currElem.Value.(*ast.CompositeLit).Elts {
							var nameOfField = schemaElem.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value
							schemaDefs[i].FieldName = nameOfField

							// E.g. DataSources often have Filters. Schema definitions seem to have <key=Filter, Value=FunctionName>. Handling that here.
							if _, isCompositLit := schemaElem.(*ast.KeyValueExpr).Value.(*ast.CompositeLit); !isCompositLit {
								schemaDefs[i].FieldName = fmt.Sprintf("** %s", schemaDefs[i].FieldName)
							} else {
								for _, typeInfo := range schemaElem.(*ast.KeyValueExpr).Value.(*ast.CompositeLit).Elts {
									switch typeInfo.(*ast.KeyValueExpr).Key.(*ast.Ident).Name {
									case "Type":
										dataTypePart1 := typeInfo.(*ast.KeyValueExpr).Value.(*ast.SelectorExpr).X.(*ast.Ident).Name
										dataTypePart2 := typeInfo.(*ast.KeyValueExpr).Value.(*ast.SelectorExpr).Sel.Name
										schemaDefs[i].DataType = fmt.Sprintf("%s.%s", dataTypePart1, dataTypePart2)

									case "Computed":
										schemaDefs[i].IsComputed, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

									case "Required":
										schemaDefs[i].IsRequired, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

									case "Optional":
										schemaDefs[i].IsOptional, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

									case "ForceNew":
										schemaDefs[i].IsForceNew, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

									default:
										fmt.Println("\t****Unknown TypeInfo****")
										fmt.Printf("\tNameOfSchemaField: %s, Attribute: %s\n", nameOfField, typeInfo.(*ast.KeyValueExpr).Key.(*ast.Ident).Name)
										fmt.Println("\t*********")
										fmt.Println()
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return
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
				funcNames[resourceFuncName] = false
			}
		}
	}

	return funcNames
}
