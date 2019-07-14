# foo.py
import random
import pandas as pd
import pyarrow as pa

random.seed(3)


def zero_copy_chunks(num_chunks=5):
    a1 = pa.array([1, 2, 3, 4])
    a2 = pa.array(['foo', 'bar', 'baz', None])
    a3 = pa.array([True, None, False, True])
    data = [a1, a2, a3]
    batch = pa.RecordBatch.from_arrays(data, ['f0', 'f1', 'f2'])
    batches = [batch] * num_chunks
    table = pa.Table.from_batches(batches)
    return table
