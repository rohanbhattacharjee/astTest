package main

type SchemaDef struct {
	FieldName  string
	DataType   string
	IsRequired bool
	IsOptional bool
	IsComputed bool
	IsForceNew bool
}
