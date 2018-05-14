# Lazy FD

When dealing with large amounts of data it can be helpful to use something like an external merge sort

   https://en.wikipedia.org/wiki/External_sorting#External_merge_sort

However, it can be hard to get the details right. The simple case can use `LazyFileReaderSimple` which is enough for most cases, check out the tests and benchmarks for example uses.

For more complex cases, there's a buffered version, which keeps a buffer full as needed, but only opens the file descriptor for as long as needed. This keeps reads high and the number of file descriptors to a minimum, which can be an issue if you have to deal with hundreds of files.

