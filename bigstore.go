package bigstore

/*
#include "c_tick_store.h"
#include <stdlib.h>

#cgo amd64 CFLAGS: -DARCH_AMD64=1
#cgo arm64 CFLAGS: -DARCH_ARM64=1

#cgo CFLAGS: -I${SRCDIR}/include
#cgo #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib -ltick_store_amd64 -lstdc++ -lpthread
#cgo #cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib -ltick_store_arm64 -lstdc++ -lpthread
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Block size constants
const (
	BlockSz32     = int32(C.BLOCK_SZ_32)
	BlockSz64     = int32(C.BLOCK_SZ_64)
	BlockSz128    = int32(C.BLOCK_SZ_128)
	BlockSz256    = int32(C.BLOCK_SZ_256)
	BlockSz512    = int32(C.BLOCK_SZ_512)
	BlockSz1024   = int32(C.BLOCK_SZ_1024)
	BlockSz2048   = int32(C.BLOCK_SZ_2048)
	BlockSz4096   = int32(C.BLOCK_SZ_4096)
	BlockSz8192   = int32(C.BLOCK_SZ_8192)
	BlockSz16384  = int32(C.BLOCK_SZ_16384)
	BlockSz32768  = int32(C.BLOCK_SZ_32768)
	BlockSz65536  = int32(C.BLOCK_SZ_65536)
	BlockSz131072 = int32(C.BLOCK_SZ_131072)
	BlockSz262144 = int32(C.BLOCK_SZ_262144)
	BlockSz524288 = int32(C.BLOCK_SZ_524288)
)

// Open mode constants
const (
	ModeOpen         = int32(C.TICK_STORE_OPEN)
	ModeCreate       = int32(C.TICK_STORE_CREATE)
	ModeOpenOrCreate = int32(C.TICK_STORE_OPEN_OR_CREATE)
)

// TickStore is a Go wrapper around the C tick_store handle.
type TickStore struct {
	handle    C.TickStoreHandle
	blockSize int32
}

// New creates a new TickStore with the given block size.
// blockSize should be one of the BlockSz* constants.
func New(blockSize int32) (*TickStore, error) {
	h := C.tick_store_new(C.int32_t(blockSize))
	if h == nil {
		return nil, errors.New("bigstore: failed to create tick store, invalid block size")
	}
	return &TickStore{handle: h, blockSize: blockSize}, nil
}

// Open opens or creates the store file.
// mode should be one of ModeOpen, ModeCreate, ModeOpenOrCreate.
// initFileSize is the initial file size in bytes (only used when creating).
func (s *TickStore) Open(path string, mode int32, initFileSize uint64) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	rc := C.tick_store_open(s.handle, cpath, C.int32_t(mode), C.size_t(initFileSize))
	if rc != 0 {
		return errors.New("bigstore: open failed")
	}
	return nil
}

// Close closes the store file (flushes data to disk).
func (s *TickStore) Close() {
	if s.handle != nil {
		C.tick_store_close(s.handle)
	}
}

// Free releases all resources. Must be called when done.
func (s *TickStore) Free() {
	if s.handle != nil {
		C.tick_store_free(s.handle)
		s.handle = nil
	}
}

// AddCode registers a security code.
// maxRecords=0 means unlimited growth; >0 means ring buffer with that capacity.
func (s *TickStore) AddCode(code string, maxRecords uint32) error {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	rc := C.tick_store_add_code(s.handle, ccode, C.uint32_t(maxRecords))
	if rc != 0 {
		return errors.New("bigstore: add_code failed")
	}
	return nil
}

// PushBack appends a record for the given code.
// data must be exactly blockSize bytes (or will be truncated/zero-padded).
func (s *TickStore) PushBack(code string, data unsafe.Pointer) error {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	// buf := make([]byte, s.blockSize)
	// copy(buf, data)

	rc := C.tick_store_push_back(s.handle, ccode, data)
	if rc != 0 {
		return errors.New("bigstore: push_back failed")
	}
	return nil
}

// Size returns the number of records stored for the given code.
func (s *TickStore) Size(code string) uint64 {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return uint64(C.tick_store_size(s.handle, ccode))
}

// At returns a copy of the record at the given index for the given code.
// Returns nil if the index is out of range.
func (s *TickStore) At(code string, index int32) unsafe.Pointer {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	ptr := C.tick_store_at(s.handle, ccode, C.int32_t(index))
	if ptr == nil {
		return nil
	}
	// Copy from the mmap'd region into Go-managed memory
	// out := make([]byte, s.blockSize)
	// copy(out, unsafe.Slice((*byte)(ptr), s.blockSize))
	return ptr
}

func (s *TickStore) SetAt(code string, index int32, data unsafe.Pointer) error {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	rc := C.tick_store_set_at(s.handle, ccode, C.int32_t(index), data)
	if rc != 0 {
		return errors.New("bigstore: set_at failed")
	}
	return nil
}

func (s *TickStore) GetValue(code string, index int32, outRecord unsafe.Pointer) int32 {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	return int32(C.tick_store_get_value(s.handle, ccode, C.int32_t(index), outRecord))
}

// Flush forces data to be synced to disk.
func (s *TickStore) Flush() {
	C.tick_store_flush(s.handle)
}

// GetAllCodes returns all registered security codes.
func (s *TickStore) GetAllCodes() []string {
	var count C.size_t
	clist := C.tick_store_get_all_codes(s.handle, &count)
	if clist == nil || count == 0 {
		return nil
	}
	defer C.tick_store_free_code_list(clist, count)

	n := int(count)
	codes := make([]string, n)
	ptrs := unsafe.Slice(clist, n)
	for i := 0; i < n; i++ {
		codes[i] = C.GoString(ptrs[i])
	}
	return codes
}

// IsRingBuffer returns true if the code is configured as a ring buffer.
func (s *TickStore) IsRingBuffer(code string) bool {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return C.tick_store_is_ring_buffer(s.handle, ccode) != 0
}

// BlockSize returns the block size of this store.
func (s *TickStore) BlockSize() int32 {
	return s.blockSize
}
