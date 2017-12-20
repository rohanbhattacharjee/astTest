package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func printSchema(fileName string, funcNamesWithSchemaDefs map[string]bool) {
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
				fmt.Printf("\tFunction name: %s()\n\n", function.Name.String())

				var schemaProperties = parseSchemaFromFunction(function)

				for _, property := range schemaProperties {
					fmt.Printf("\t\tName: %s Type: %s isRequired: %v isOptional: %v isComputed: %v isForceNew: %v\n", property.FieldName, property.DataType, property.IsRequired, property.IsOptional, property.IsComputed, property.IsForceNew)
				}

				fmt.Println()
			}
		}
	}

	return
}

func parseSchemaFromFunction(function *ast.FuncDecl) (schemaProperties []SchemaProperties) {
	var elements = function.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts
	for _, elem := range elements {
		currElem := elem.(*ast.KeyValueExpr)
		if currElem.Key.(*ast.Ident).Name == "Schema" {

			schemaProperties = make([]SchemaProperties, len(currElem.Value.(*ast.CompositeLit).Elts))

			for i, schemaElem := range currElem.Value.(*ast.CompositeLit).Elts {
				var nameOfField = schemaElem.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value
				schemaProperties[i].FieldName = nameOfField

				// E.g. DataSources often have Filters. Schema definitions seem to have <key=Filter, Value=FunctionName>. Handling that here.
				if _, isCompositLit := schemaElem.(*ast.KeyValueExpr).Value.(*ast.CompositeLit); !isCompositLit {
					fmt.Println("\t\t**** Look at this field ****")
					fmt.Printf("\t\t%s\n", schemaProperties[i].FieldName)
					fmt.Println("\t\t*********")
					fmt.Println()
				} else {
					for _, typeInfo := range schemaElem.(*ast.KeyValueExpr).Value.(*ast.CompositeLit).Elts {
						switch typeInfo.(*ast.KeyValueExpr).Key.(*ast.Ident).Name {
						case "Type":
							dataTypePart1 := typeInfo.(*ast.KeyValueExpr).Value.(*ast.SelectorExpr).X.(*ast.Ident).Name
							dataTypePart2 := typeInfo.(*ast.KeyValueExpr).Value.(*ast.SelectorExpr).Sel.Name
							schemaProperties[i].DataType = fmt.Sprintf("%s.%s", dataTypePart1, dataTypePart2)

						case "Computed":
							schemaProperties[i].IsComputed, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

						case "Required":
							schemaProperties[i].IsRequired, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

						case "Optional":
							schemaProperties[i].IsOptional, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

						case "ForceNew":
							schemaProperties[i].IsForceNew, _ = strconv.ParseBool(typeInfo.(*ast.KeyValueExpr).Value.(*ast.Ident).Name)

						default:
							fmt.Println("\t\t**** Unknown attribute ****")
							fmt.Printf("\t\tNameOfSchemaField: %s, Attribute: %s\n", nameOfField, typeInfo.(*ast.KeyValueExpr).Key.(*ast.Ident).Name)
							fmt.Println("\t\t*********")
							fmt.Println()
						}
					}
				}
			}
		}
	}

	return
}
