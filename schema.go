package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
)

// PySchemaFromPyTable returns a pyarrow schema from a pyarrow Table.
func PySchemaFromPyTable(pyTable *python3.PyObject) (*python3.PyObject, error) {
	pySchema := pyTable.GetAttrString("schema")
	if pySchema == nil {
		return nil, errors.New("could not get pySchema")
	}
	return pySchema, nil
}

// PySchemaToSchema given a Python schema gets the Go Arrow schema.
func PySchemaToSchema(pySchema *python3.PyObject) (*arrow.Schema, error) {
	// start with the field names
	pyFieldNames, err := getPyFieldNames(pySchema)
	if err != nil {
		return nil, err
	}
	defer func() {
		for i := range pyFieldNames {
			pyFieldNames[i].DecRef()
		}
	}()

	// Get the fields
	fields, err := getFields(pySchema, pyFieldNames)
	if err != nil {
		return nil, err
	}

	return arrow.NewSchema(fields, nil), nil
}

func getPyFieldNames(pySchema *python3.PyObject) ([]*python3.PyObject, error) {
	pyFieldNames := pySchema.GetAttrString("names")
	if pyFieldNames == nil {
		return nil, errors.New("could not get pyFieldNames")
	}
	defer pyFieldNames.DecRef()

	// verify the result is a list
	if !python3.PyList_Check(pyFieldNames) {
		return nil, errors.New("not a list of field names")
	}

	length := python3.PyList_Size(pyFieldNames)
	pyNames := make([]*python3.PyObject, 0, length)
	for i := 0; i < length; i++ {
		pyName := python3.PyList_GetItem(pyFieldNames, i)
		if pyName == nil {
			return nil, errors.New("could not get name")
		}
		pyName.IncRef()
		// pyNames[i] = pyName
		pyNames = append(pyNames, pyName)
	}

	return pyNames, nil
}

func getFields(pySchema *python3.PyObject, pyFieldNames []*python3.PyObject) ([]arrow.Field, error) {
	fields := make([]arrow.Field, 0, len(pyFieldNames))
	for _, pyFieldName := range pyFieldNames {
		field, err := getField(pySchema, pyFieldName)
		if err != nil {
			return nil, err
		}
		// fields[i] = *field
		fields = append(fields, *field)
	}
	return fields, nil
}

func getField(schema *python3.PyObject, fieldName *python3.PyObject) (*arrow.Field, error) {
	pyField := CallPyFunc(schema, "field_by_name", fieldName)
	if pyField == nil {
		return nil, errors.New("could not get pyField")
	}
	defer pyField.DecRef()

	field, err := PyFieldToField(pyField)
	if err != nil {
		return nil, err
	}

	return field, nil
}
