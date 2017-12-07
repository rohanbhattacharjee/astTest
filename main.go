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

	var filesProcessed, filesSkipped int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !isResourceFile(file.Name()) && !isDataSourceFile(file.Name()) {
			continue
		}

		if doesFileHaveProblems(file.Name()) {
			filesSkipped += 1
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
		filesProcessed += 1
	}

	fmt.Printf("Processed file count = %v, Skipped file count = %v\n", filesProcessed, filesSkipped)
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

func doesFileHaveProblems(fileName string) bool {
	problemCode := map[string]bool{
		"identity_compartment_resource.go":         true,
		"identity_group_resource.go":               true,
		"identity_policy_resource.go":              true,
		"identity_user_resource.go":                true,
		"load_balancer_backendset_resource.go":     true,
		"load_balancer_listener_resource.go":       true,
		"objectstorage_bucket_resource.go":         true,
		"objectstorage_object_resource.go":         true,
		"objectstorage_preauthrequest_resource.go": true,
	}

	if _, ok := problemCode[fileName]; ok {
		return true
	}

	return false
}
