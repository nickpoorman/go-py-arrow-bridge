package main

import (
	"fmt"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/go-bullseye/bullseye/dataframe"
	bridge "github.com/nickpoorman/go-py-arrow-bridge"
	"github.com/nickpoorman/pytasks"
)

func main() {
	py := pytasks.GetPythonSingleton()

	fooModule, err := py.ImportModule("foo")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := pytasks.GetPythonSingleton().NewTaskSync(func() {
			fooModule.DecRef()
		})
		if err != nil {
			panic(err)
		}
	}()

	var table array.Table
	taskErr := py.NewTaskSync(func() {
		pyTable := genPyTable(fooModule)
		table, err = bridge.PyTableToTable(pyTable)
		pyTable.DecRef()
	})
	if taskErr != nil {
		panic(taskErr)
	}
	if err != nil {
		panic(err)
	}

	// Wrapping it in a bullseye dataframe allows us to print it easily
	pool := memory.NewGoAllocator()
	df, err := dataframe.NewDataFrameFromTable(pool, table)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nArrow Table from Python now in Go:")
	fmt.Println(df.Display(0))

	// Arrow Table from Python now in Go:
	// rec[0]["f0"]: [1 2 3 4]
	// rec[0]["f1"]: ["foo" "bar" "baz" (null)]
	// rec[0]["f2"]: [true (null) false true]
	// rec[1]["f0"]: [1 2 3 4]
	// rec[1]["f1"]: ["foo" "bar" "baz" (null)]
	// rec[1]["f2"]: [true (null) false true]
	// rec[2]["f0"]: [1 2 3 4]
	// rec[2]["f1"]: ["foo" "bar" "baz" (null)]
	// rec[2]["f2"]: [true (null) false true]
	// rec[3]["f0"]: [1 2 3 4]
	// rec[3]["f1"]: ["foo" "bar" "baz" (null)]
	// rec[3]["f2"]: [true (null) false true]
	// rec[4]["f0"]: [1 2 3 4]
	// rec[4]["f1"]: ["foo" "bar" "baz" (null)]
	// rec[4]["f2"]: [true (null) false true]
}

func genPyTable(module *python3.PyObject) *python3.PyObject {
	pyTable := bridge.CallPyFunc(module, "zero_copy_chunks")
	if pyTable == nil {
		panic("pyTable is nil")
	}
	return pyTable
}
