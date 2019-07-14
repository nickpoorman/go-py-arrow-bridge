package bridge

import (
	"errors"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

func PyChunkedToChunked(pyChunked *python3.PyObject, dtype arrow.DataType) (*array.Chunked, error) {
	// Convert pyChunks to []Interface
	chunks, err := PyChunkedToChunks(pyChunked, dtype)
	if err != nil {
		return nil, err
	}

	chunked := array.NewChunked(dtype, chunks)
	return chunked, nil
}

func PyChunkedToChunks(pyChunked *python3.PyObject, dtype arrow.DataType) ([]array.Interface, error) {
	pyChunks, err := PyChunkedGetPyChunks(pyChunked)
	if err != nil {
		return nil, err
	}
	defer pyChunks.DecRef()

	if !python3.PyList_Check(pyChunks) {
		return nil, errors.New("pyChunks is not a list")
	}

	length := python3.PyList_Size(pyChunks)
	chunks := make([]array.Interface, length)
	for i := 0; i < length; i++ {
		chunk, err := PyChunksGetChunk(pyChunks, i, dtype)
		if err != nil {
			return nil, err
		}
		chunks[i] = chunk
	}

	return chunks, nil
}

func PyChunkedGetPyChunks(pyChunked *python3.PyObject) (*python3.PyObject, error) {
	pyChunks := pyChunked.GetAttrString("chunks")
	if pyChunks == nil {
		return nil, errors.New("could not get pyChunks")
	}
	return pyChunks, nil
}

func PyChunksGetChunk(pyChunks *python3.PyObject, i int, dtype arrow.DataType) (array.Interface, error) {
	pyChunk, err := PyChunksGetPyChunk(pyChunks, i)
	if err != nil {
		return nil, err
	}
	defer pyChunk.DecRef()

	chunk, err := PyChunkToChunk(pyChunk, dtype)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

func PyChunksGetPyChunk(pyChunks *python3.PyObject, i int) (*python3.PyObject, error) {
	pyChunk := python3.PyList_GetItem(pyChunks, i)
	if pyChunk == nil {
		return nil, errors.New("could not get pyChunk from list")
	}
	return pyChunk, nil
}

func PyChunkToChunk(pyChunk *python3.PyObject, dtype arrow.DataType) (array.Interface, error) {
	data, err := PyChunkToData(pyChunk, dtype)
	if err != nil {
		return nil, err
	}
	defer data.Release()
	chunk := array.MakeFromData(data)
	return chunk, nil
}

func PyChunkToData(pyChunk *python3.PyObject, dtype arrow.DataType) (*array.Data, error) {
	buffers, err := PyChunkGetBuffers(pyChunk)
	if err != nil {
		return nil, err
	}

	nullCount, err := PyChunkGetNullCount(pyChunk)
	if err != nil {
		return nil, err
	}

	offset, err := PyChunkGetOffset(pyChunk)
	if err != nil {
		return nil, err
	}

	chunkLen, err := PyChunkGetLength(pyChunk)
	if err != nil {
		return nil, err
	}

	var childData []*array.Data // TODO: Implement
	data := array.NewData(dtype, chunkLen, buffers, childData, nullCount, offset)
	return data, nil
}

func PyChunkGetBuffers(pyChunk *python3.PyObject) ([]*memory.Buffer, error) {
	pyBuffers, err := PyChunkGetPyBuffers(pyChunk)
	if err != nil {
		return nil, err
	}
	defer pyBuffers.DecRef()

	return PyBuffersToBuffers(pyBuffers)
}

func PyChunkGetPyBuffers(pyChunk *python3.PyObject) (*python3.PyObject, error) {
	pyBuffersFunc := pyChunk.GetAttrString("buffers")
	if pyBuffersFunc == nil {
		return nil, errors.New("could not get pyBuffersFunc")
	}
	defer pyBuffersFunc.DecRef()

	pyBuffers := pyBuffersFunc.CallFunctionObjArgs()
	if pyBuffers == nil {
		return nil, errors.New("could not get pyBuffers")
	}

	return pyBuffers, nil
}

func PyChunkGetNullCount(pyChunk *python3.PyObject) (int, error) {
	v, ok := GetIntAttr(pyChunk, "null_count")
	if !ok {
		return 0, errors.New("could not get null_count")
	}
	return v, nil
}

func PyChunkGetOffset(pyChunk *python3.PyObject) (int, error) {
	v, ok := GetIntAttr(pyChunk, "offset")
	if !ok {
		return 0, errors.New("could not get offset")
	}
	return v, nil
}

func PyChunkGetLength(pyChunk *python3.PyObject) (int, error) {
	pyLength := CallPyFunc(pyChunk, "__len__")
	if pyLength == nil {
		return 0, errors.New("could not get pyChunk.__len__()")
	}
	defer pyLength.DecRef()
	length := python3.PyLong_AsLong(pyLength)
	return length, nil
}
