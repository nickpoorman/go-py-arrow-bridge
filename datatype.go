package bridge

import (
	"errors"
	"fmt"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
)

// PyDataTypeToDataType returns the Go arrow DataType given the Python type.
func PyDataTypeToDataType(pyDtype *python3.PyObject) (arrow.DataType, error) {
	// get the id
	id, err := PyDataTypeGetID(pyDtype)
	if err != nil {
		return nil, err
	}

	t := arrow.Type(id)
	return GetFromType(t)
}

func PyDataTypeGetID(pyDtype *python3.PyObject) (int, error) {
	v, ok := GetIntAttr(pyDtype, "id")
	if !ok {
		return 0, errors.New("could not get pyDtype.id")
	}
	return v, nil
}

var (
	dataTypeForType [32]arrow.DataType
)

// GetFromType returns a arrow.DataType for a given arrow.Type
func GetFromType(t arrow.Type) (arrow.DataType, error) {
	dtype := dataTypeForType[byte(t&0x1f)]
	if dtype == nil {
		return nil, fmt.Errorf("DataType for id=%v is not yet implemented", t)
	}
	return dtype, nil
}

// DataTypeFromType returns an arrow.DataType given the Type
func init() {
	dataTypeForType = [...]arrow.DataType{
		arrow.NULL:              arrow.Null,
		arrow.BOOL:              arrow.FixedWidthTypes.Boolean,
		arrow.UINT8:             arrow.PrimitiveTypes.Uint8,
		arrow.INT8:              arrow.PrimitiveTypes.Int8,
		arrow.UINT16:            arrow.PrimitiveTypes.Uint16,
		arrow.INT16:             arrow.PrimitiveTypes.Int16,
		arrow.UINT32:            arrow.PrimitiveTypes.Uint32,
		arrow.INT32:             arrow.PrimitiveTypes.Int32,
		arrow.UINT64:            arrow.PrimitiveTypes.Uint64,
		arrow.INT64:             arrow.PrimitiveTypes.Int64,
		arrow.FLOAT16:           arrow.FixedWidthTypes.Float16,
		arrow.FLOAT32:           arrow.PrimitiveTypes.Float32,
		arrow.FLOAT64:           arrow.PrimitiveTypes.Float64,
		arrow.STRING:            arrow.BinaryTypes.String,
		arrow.BINARY:            arrow.BinaryTypes.Binary,
		arrow.FIXED_SIZE_BINARY: nil, // arrow.FixedSizeBinaryType,
		arrow.DATE32:            arrow.PrimitiveTypes.Date32,
		arrow.DATE64:            arrow.PrimitiveTypes.Date64,
		arrow.TIMESTAMP:         nil, // arrow.FixedWidthTypes.Timestamp_s, // TODO
		arrow.TIME32:            nil, // arrow.FixedWidthTypes.Time32s, // TODO
		arrow.TIME64:            nil, // arrow.FixedWidthTypes.Time64us, // TODO
		arrow.INTERVAL:          nil, // arrow.FixedWidthTypes.MonthInterval, // TODO
		arrow.DECIMAL:           nil,
		arrow.LIST:              nil,
		arrow.STRUCT:            nil,
		arrow.UNION:             nil,
		arrow.DICTIONARY:        nil,
		arrow.MAP:               nil,
		arrow.EXTENSION:         nil,
		arrow.FIXED_SIZE_LIST:   nil,
		arrow.DURATION:          nil, // arrow.FixedWidthTypes.Duration_s, // TODO

		// invalid data types to fill out array size 2⁵-1
		31: nil,
	}
}
