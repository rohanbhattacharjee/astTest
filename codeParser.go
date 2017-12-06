package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func myParseSchema(fileName string, functionName string) (schemaDefs []SchemaDef) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	//ast.Print(fset, f)

	// Uncomment the statement above to see the abstract syntax tree (AST) for the code.
	// Once you see the AST, the following code will make sense.
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			if function.Name.String() == functionName {
				var elements = function.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts

				for _, elem := range elements {
					currElem := elem.(*ast.KeyValueExpr)
					if currElem.Key.(*ast.Ident).Name == "Schema" {
						schemaDefs = make([]SchemaDef, len(currElem.Value.(*ast.CompositeLit).Elts))

						for i, schemaElem := range currElem.Value.(*ast.CompositeLit).Elts {
							var nameOfField = schemaElem.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value
							schemaDefs[i].FieldName = nameOfField

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
