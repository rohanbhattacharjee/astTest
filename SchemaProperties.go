package main

type SchemaProperties struct {
	FieldName  string
	DataType   string
	IsRequired bool
	IsOptional bool
	IsComputed bool
	IsForceNew bool
}
