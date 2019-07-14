package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
)

// PyColumnToColumnWithField turns a PyColumn into a GoColumn
func PyColumnToColumnWithField(pyColumn *python3.PyObject, field arrow.Field) (*array.Column, error) {
	chunks, err := PyColumnToChunkedWithField(pyColumn, field)
	if err != nil {
		return nil, err
	}

	col := array.NewColumn(field, chunks)
	return col, nil
}

func PyColumnToChunkedWithField(pyColumn *python3.PyObject, field arrow.Field) (*array.Chunked, error) {
	pyChunked, err := PyColumnGetPyChunked(pyColumn)
	if err != nil {
		return nil, err
	}
	defer pyChunked.DecRef()

	return PyChunkedToChunked(pyChunked, field.Type)
}

func PyColumnGetPyChunked(pyColumn *python3.PyObject) (*python3.PyObject, error) {
	pyChunked := pyColumn.GetAttrString("data")
	if pyChunked == nil {
		return nil, errors.New("could not get pyChunked")
	}
	return pyChunked, nil
}
