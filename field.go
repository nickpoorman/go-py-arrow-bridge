package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
)

// PyFieldToField given a Python field gets the Go Arrow field.
func PyFieldToField(pyField *python3.PyObject) (*arrow.Field, error) {
	pyName := pyField.GetAttrString("name")
	if pyName == nil {
		return nil, errors.New("could not get pyName")
	}
	defer pyName.DecRef()

	pyDtype := pyField.GetAttrString("type")
	if pyDtype == nil {
		return nil, errors.New("could not get pyDtype")
	}
	defer pyDtype.DecRef()

	pyNullable := pyField.GetAttrString("nullable")
	if pyNullable == nil {
		return nil, errors.New("could not get pyNullable")
	}
	defer pyNullable.DecRef()

	// TODO: Implement
	// pyMetadata := CallPyFunc(pyField, "metadata")
	// if pyMetadata == nil {
	// 	return nil, errors.New("could not get pyMetadata")
	// }

	name := python3.PyUnicode_AsUTF8(pyName)
	dtype, err := PyDataTypeToDataType(pyDtype)
	if err != nil {
		return nil, err
	}
	nullable := python3.PyBool_Check(pyNullable)

	field := &arrow.Field{
		Name:     name,
		Type:     dtype,
		Nullable: nullable,
		// TODO: Implement
		// Metadata: metadata,
	}

	return field, nil
}
