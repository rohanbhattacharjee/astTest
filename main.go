package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Println("Analyzing code ...")

	var sourceDir = "/Users/rohabhat/Documents/work/code/go/src/github.com/oracle/terraform-provider-oci/provider"

	var funcNamesWithSchemaDefs = readProviderForFuncNamesWithSchemaDefs(sourceDir)

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		panic("Could not read source directory")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !isResourceFile(file.Name()) && !isDataSourceFile(file.Name()) {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", sourceDir, file.Name())

		fmt.Printf("File: %s:\n", file.Name())

		// functionName := "ConsoleHistoryDataDatasource" //"ConsoleHistoryResource"
		schemaDefs := parseSchema(filePath, funcNamesWithSchemaDefs)

		for _, schemaDef := range schemaDefs {
			fmt.Printf("\tName: %s Type: %s isRequired: %v isOptional: %v isComputed: %v isForceNew: %v\n", schemaDef.FieldName, schemaDef.DataType, schemaDef.IsRequired, schemaDef.IsOptional, schemaDef.IsComputed, schemaDef.IsForceNew)
		}

		fmt.Println()
	}
}

// This function basically looks that two functions defined in the provider.go, parses a map to get the list of functions that have schemas defined.
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

func isResourceFile(fileName string) bool {
	return strings.HasSuffix(fileName, "_resource.go")
}

func isDataSourceFile(fileName string) bool {
	return strings.HasSuffix(fileName, "_data_source.go")
}
