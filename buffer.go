package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow/memory"
)

func PyBuffersToBuffers(pyBuffers *python3.PyObject) ([]*memory.Buffer, error) {
	// First buffer is the null mask buffer, second is the values.
	// [<pyarrow.lib.Buffer object at 0x113d46a08>, <pyarrow.lib.Buffer object at 0x114761998>]
	if !python3.PyList_Check(pyBuffers) {
		return nil, errors.New("pyBuffers is not a list")
	}

	length := python3.PyList_Size(pyBuffers)
	buffers := make([]*memory.Buffer, 0, length)
	for i := 0; i < length; i++ {
		buffer, err := PyBuffersGetBuffer(pyBuffers, i)
		if err != nil {
			return nil, err
		}
		// buffers[i] = buffer
		buffers = append(buffers, buffer)
	}

	return buffers, nil
}

func PyBuffersGetBuffer(pyBuffers *python3.PyObject, i int) (*memory.Buffer, error) {
	// Get the buffer at index i
	pyBuffer := python3.PyList_GetItem(pyBuffers, i)
	if pyBuffer == nil {
		return nil, errors.New("could not get pyBuffer")
	}
	defer pyBuffer.DecRef()

	goBytes, err := PyBufferToBytes(pyBuffer)
	if err != nil {
		return nil, err
	}

	buffer := memory.NewBufferBytes(goBytes)
	return buffer, nil
}

func PyBufferToBytes(pyBuffer *python3.PyObject) ([]byte, error) {
	// <pyarrow.lib.Buffer object at 0x113d46a08>
	// Convert the buffer to our Py_buffer struct type
	pyBuf, err := python3.PyObject_GetBuffer(pyBuffer, python3.PyBUF_SIMPLE)
	if err {
		return nil, errors.New("could not get pyBuf")
	}

	goBytes := python3.PyObject_GetBufferBytes(pyBuf)
	if goBytes == nil {
		return nil, errors.New("could not get goBytes")
	}

	return goBytes, nil
}
