package lazy_fd

import (
	"os"
)

type LazyFileReader interface {
	Read([]byte) (int, error)
}

type LazyFileReaderSimple struct {
	Filename          string
	LastReadFileIndex int64
	SeekCount         int
}

func NewLazyFileReaderSimple(filename string) *LazyFileReaderSimple {
	return &LazyFileReaderSimple{
		Filename:          filename,
		LastReadFileIndex: 0,
		SeekCount:         0,
	}
}

func (lfr *LazyFileReaderSimple) Read(p []byte) (int, error) {
	f, err := os.Open(lfr.Filename)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	f.Seek(lfr.LastReadFileIndex, os.SEEK_SET)
	lfr.SeekCount += 1

	n, err := f.Read(p)
	lfr.LastReadFileIndex += int64(n)

	return n, err
}

type LazyFileReaderBuffer struct {
	Filename          string
	LastReadFileIndex int64
	SeekCount         int64
	ReadCount         int64
	Buffer            []byte
	CurrentStartIndex int
	CurrentEndIndex   int
}

// Create a new file reader with a pre-allocated buffer.
func NewLazyFileReaderBuffer(filename string, buffer_size int) *LazyFileReaderBuffer {
	return &LazyFileReaderBuffer{
		Filename:          filename,
		LastReadFileIndex: 0,
		SeekCount:         0,
		ReadCount:         0,
		Buffer:            make([]byte, buffer_size),
		CurrentStartIndex: 0,
		CurrentEndIndex:   0,
	}
}

// the logic here is a little tricky, so lets outline it.
//
// 1. If we already have enough bytes in the buffer, just copy
//    those into the p array and return.
//
// 2. If we don't have enough bytes, copy all the data out of
//    the buffer and then do a full read into our buffer
//    this requires no allocs and we've already copied all the
//    the data we care about.
//
// Internally we need to keep track fo the buffer and what
// part of the internal buffer has been returned as apart of
// a previous call.
func (lfr *LazyFileReaderBuffer) Read(p []byte) (int, error) {
	var last_p_index int = 0

	if len(p) <= lfr.CurrentEndIndex-lfr.CurrentStartIndex {
		newEnd := lfr.CurrentStartIndex + len(p)
		copy(
			p[last_p_index:],
			lfr.Buffer[lfr.CurrentStartIndex:newEnd],
		)
		lfr.CurrentStartIndex = newEnd
		return len(p), nil
	} else {
		// we know p is larger than our buffer, copy all of what we have.
		copy(
			p[last_p_index:],
			lfr.Buffer[lfr.CurrentStartIndex:lfr.CurrentEndIndex],
		)
		last_p_index += lfr.CurrentEndIndex - lfr.CurrentStartIndex
		// We've now used up our internal buffer, so lets zero it out
		// and our internal representation.
		lfr.resetBuffer()
	}

	// open our file for reading, we may need to read multiple times
	// in order to fill p.
	f, err := os.Open(lfr.Filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// Seek to the end of the last read location from our last file read.
	f.Seek(lfr.LastReadFileIndex, os.SEEK_SET)
	lfr.SeekCount += 1

	// last location filled in p.
	var p_bytes_remain int

	for last_p_index <= len(p) {
		p_bytes_remain = len(p) - last_p_index

		// We've now used up our internal buffer, so lets zero it out
		// and our internal representation.
		lfr.resetBuffer()

		n, err := f.Read(lfr.Buffer)
		lfr.ReadCount += 1
		lfr.LastReadFileIndex += int64(n)
		lfr.CurrentEndIndex = n

		// If the number of bytes read is larger than p,
		// copy what we can and return it and an error
		// if one happened.
		if n > p_bytes_remain {
			copy(p[last_p_index:], lfr.Buffer[:p_bytes_remain])
			lfr.CurrentStartIndex = p_bytes_remain
			return last_p_index + p_bytes_remain, err
		}

		// p is greater than n, so we copy what we have
		// and the loop will run again.
		copy(p[last_p_index:], lfr.Buffer[:n])
		last_p_index += n

		if err != nil {
			return last_p_index, err
		}
	}

	return last_p_index, nil
}

func (lfr *LazyFileReaderBuffer) resetBuffer() {
	lfr.Buffer = lfr.Buffer[:]
	lfr.CurrentStartIndex = 0
	lfr.CurrentEndIndex = 0
}
