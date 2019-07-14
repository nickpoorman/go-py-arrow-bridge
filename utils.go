package bridge

import "github.com/DataDog/go-python3"

// A helper for first fetching the function and then calling it
func CallPyFunc(obj *python3.PyObject, name string, args ...*python3.PyObject) *python3.PyObject {
	fn := obj.GetAttrString(name)
	defer fn.DecRef()

	return fn.CallFunctionObjArgs(args...)
}

func GetIntAttr(obj *python3.PyObject, attr string) (int, bool) {
	v := obj.GetAttrString(attr)
	if v == nil {
		return 0, false
	}
	defer v.DecRef()
	return python3.PyLong_AsLong(v), true
}
