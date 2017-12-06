package main

import (
	"fmt"
)

func main() {
	fmt.Println("Analyzing code ...")

	fileName := "./otherSrc/someCode.go"
	functionName := "ConsoleHistoryResource"

	schemaDefs := myParseSchema(fileName, functionName)

	for _, schemaDef := range schemaDefs {
		fmt.Printf("Name: %s Type: %s isRequired: %v isOptional: %v isComputed: %v isForceNew: %v\n", schemaDef.FieldName, schemaDef.DataType, schemaDef.IsRequired, schemaDef.IsOptional, schemaDef.IsComputed, schemaDef.IsForceNew)
	}
}
