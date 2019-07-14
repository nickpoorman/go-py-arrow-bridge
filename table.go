package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
)

func PyTableToTable(pyTable *python3.PyObject) (array.Table, error) {
	schema, cols, err := PyTableToColumns(pyTable)
	if err != nil {
		return nil, err
	}

	// Build the table
	table := array.NewTable(schema, cols, -1) // -1 tells it to determine the numRows from the first column

	return table, nil
}

// PyTableToColumns returns the records in the pyarrow table.
func PyTableToColumns(pyTable *python3.PyObject) (*arrow.Schema, []array.Column, error) {
	// Get the PySchema from the PyTable
	pySchema, err := PySchemaFromPyTable(pyTable)
	if err != nil {
		return nil, nil, err
	}
	defer pySchema.DecRef()

	// Get the GoSchema
	schema, err := PySchemaToSchema(pySchema)
	if err != nil {
		return nil, nil, err
	}

	columns, err := PyTableToColumnsWithSchema(pyTable, schema)
	if err != nil {
		return nil, nil, err
	}

	return schema, columns, nil
}

// PyTableToColumns returns the columns in the pyarrow table.
func PyTableToColumnsWithSchema(pyTable *python3.PyObject, schema *arrow.Schema) ([]array.Column, error) {
	fields := schema.Fields()
	columns := make([]array.Column, 0, len(fields))

	for i := range fields {
		pyColumn, err := PyTableGetPyColumn(pyTable, fields[i].Name)
		if err != nil {
			return nil, err
		}
		defer pyColumn.DecRef()

		col, err := PyColumnToColumnWithField(pyColumn, fields[i])
		if err != nil {
			return nil, err
		}
		// columns[i] = *col
		columns = append(columns, *col)
	}

	return columns, nil
}

// PyTableGetPyColumn returns the PyColumn given the name from the PyTable
func PyTableGetPyColumn(pyTable *python3.PyObject, name string) (*python3.PyObject, error) {
	pyName := python3.PyUnicode_FromString(name)
	defer pyName.DecRef()

	pyColumn := CallPyFunc(pyTable, "column", pyName)
	if pyColumn == nil {
		return nil, errors.New("could not get pyColumn")
	}

	return pyColumn, nil
}
