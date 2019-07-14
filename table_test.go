package bridge

import (
	"fmt"
	"testing"

	"github.com/DataDog/go-python3"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/go-bullseye/bullseye/dataframe"
	"github.com/nickpoorman/pytasks"
)

func BenchmarkAll(b *testing.B) {
	for i := 5; i <= 10; i += 2 {
		b.Run(fmt.Sprintf("BenchmarkZeroCopy_%d", i), zeroCopyBenchmarkN(i))
	}
	for i := 1000; i <= 10000; i += 500 {
		b.Run(fmt.Sprintf("BenchmarkZeroCopy_%d", i), zeroCopyBenchmarkN(i))
	}

	// At this point we know we won't need Python anymore in this
	// program, we can restore the state and lock the GIL to perform
	// the final operations before exiting.
	err := pytasks.GetPythonSingleton().Finalize()
	if err != nil {
		panic(err)
	}
}

// So the benchmarks don't get compiled out during optimization.
var benchTable array.Table

func zeroCopyBenchmarkN(numChunks int) func(b *testing.B) {
	return func(b *testing.B) {
		if numChunks <= 0 {
			b.Fatal("numChunks must be greater than zero")
		}

		py := pytasks.GetPythonSingleton()

		fooModule, err := py.ImportModule("foo")
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			err := pytasks.GetPythonSingleton().NewTaskSync(func() {
				fooModule.DecRef()
			})
			if err != nil {
				panic(err)
			}
		}()

		var pyTable *python3.PyObject
		taskErr := py.NewTaskSync(func() {
			pyNumChunks := python3.PyLong_FromLong(numChunks)
			defer pyNumChunks.DecRef()
			pyTable = CallPyFunc(fooModule, "zero_copy_chunks", pyNumChunks)
			if pyTable == nil {
				b.Fatal("pyTable is nil")
			}
		})
		if taskErr != nil {
			b.Fatal(taskErr)
		}
		defer func() {
			err := pytasks.GetPythonSingleton().NewTaskSync(func() {
				pyTable.DecRef()
			})
			if err != nil {
				panic(err)
			}
		}()

		var table array.Table
		b.ResetTimer()

		taskErr = py.NewTaskSync(func() {
			// Run the loop inside of the task so we don't benchmark
			// grabbing the GIL over and over. When the loop
			// is outside the task, the results are still consistent.
			for i := 0; i < b.N; i++ {
				table, err = PyTableToTable(pyTable)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		if taskErr != nil {
			b.Fatal(taskErr)
		}

		benchTable = table
	}
}

func TestTable(t *testing.T) {
	// Init Python
	_ = pytasks.GetPythonSingleton()

	t.Run("PyTableToTable", testPyTableToTable)

	// At this point we know we won't need Python anymore in this
	// program, we can restore the state and lock the GIL to perform
	// the final operations before exiting.
	err := pytasks.GetPythonSingleton().Finalize()
	if err != nil {
		t.Fatal(err)
	}
}

func testPyTableToTable(t *testing.T) {
	pool := memory.NewCheckedAllocator(memory.NewGoAllocator())
	defer pool.AssertSize(t, 0)

	py := pytasks.GetPythonSingleton()
	fooModule, err := py.ImportModule("foo")
	if err != nil {
		t.Fatal(err)
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
		table, err = PyTableToTable(pyTable)
		pyTable.DecRef()
	})
	if taskErr != nil {
		t.Fatal(taskErr)
	}
	if err != nil {
		t.Fatal(err)
	}

	// Wrapping it in a bullseye dataframe allows us to print it easily
	df, err := dataframe.NewDataFrameFromTable(pool, table)
	if err != nil {
		t.Fatal(err)
	}

	got := df.Display(0)
	want := `rec[0]["f0"]: [1 2 3 4]
rec[0]["f1"]: ["foo" "bar" "baz" (null)]
rec[0]["f2"]: [true (null) false true]
rec[1]["f0"]: [1 2 3 4]
rec[1]["f1"]: ["foo" "bar" "baz" (null)]
rec[1]["f2"]: [true (null) false true]
rec[2]["f0"]: [1 2 3 4]
rec[2]["f1"]: ["foo" "bar" "baz" (null)]
rec[2]["f2"]: [true (null) false true]
rec[3]["f0"]: [1 2 3 4]
rec[3]["f1"]: ["foo" "bar" "baz" (null)]
rec[3]["f2"]: [true (null) false true]
rec[4]["f0"]: [1 2 3 4]
rec[4]["f1"]: ["foo" "bar" "baz" (null)]
rec[4]["f2"]: [true (null) false true]
`
	if got != want {
		t.Fatalf("\ngot=\n%v\nwant=\n%v", got, want)
	}
}

func genPyTable(module *python3.PyObject) *python3.PyObject {
	pyTable := CallPyFunc(module, "zero_copy_chunks")
	if pyTable == nil {
		panic("pyTable is nil")
	}
	return pyTable
}
