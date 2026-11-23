package main

/*
#include <stdlib.h>
*/
import "C"

import (
	core "PVote/mobile/core"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"runtime"
	"unsafe"
)

var engine core.Engine

func init() { runtime.GOMAXPROCS(1) }

//export PVoteCall
func PVoteCall(request *C.char) *C.char {
	if request == nil {
		return C.CString(`{"error":"empty native request"}`)
	}
	return C.CString(string(engine.Call([]byte(C.GoString(request)))))
}

//export PVoteFree
func PVoteFree(value *C.char) { C.free(unsafe.Pointer(value)) }

// PVoteExportReport exports an existing report without running or changing any
// measurement. The digest covers the original file bytes, including whitespace.
// This provides an offline fallback when iOS refuses developer-container copies.
//
//export PVoteExportReport
func PVoteExportReport(input *C.char) *C.char {
	if input == nil {
		return C.CString(`{"error":"missing report"}`)
	}
	raw := []byte(C.GoString(input))
	if len(raw) > 4*1024*1024 || !json.Valid(raw) {
		return C.CString(`{"error":"invalid report"}`)
	}
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(raw); err != nil {
		return C.CString(`{"error":"report compression failed"}`)
	}
	if err := writer.Close(); err != nil {
		return C.CString(`{"error":"report compression failed"}`)
	}
	digest := sha256.Sum256(raw)
	encoded, _ := json.Marshal(map[string]string{"sha256": hex.EncodeToString(digest[:]), "gzipBase64": base64.StdEncoding.EncodeToString(buffer.Bytes())})
	return C.CString(string(encoded))
}

func main() {}
