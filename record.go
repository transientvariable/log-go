package log

import (
	"context"
	"math"
	"strings"
	"sync"
	"time"
)

type kind uint

const (
	kindAny kind = iota
	kindBool
	kindDuration
	kindFloat64
	kindFloat32
	kindInt64
	kindString
	kindTime
	kindUint64
)

var recPool = sync.Pool{
	New: func() any {
		return &Record{}
	},
}

func acquireRecord() *Record {
	return recPool.Get().(*Record)
}

func releaseRecord(r *Record) {
	if r != nil {
		r.attrs = nil
		r.ctx = nil
		r.err = nil
		r.msg = ""
		recPool.Put(r)
	}
}

type attr struct {
	value any
	kind  kind
}

// Record ...
type Record struct {
	attrs map[string]attr
	ctx   context.Context
	err   error
	msg   string
}

func (r *Record) addAttr(key string, value any, kind kind) {
	if key = strings.TrimSpace(key); key != "" && value != nil {
		if r.attrs == nil {
			r.attrs = make(map[string]attr)
		}

		r.attrs[key] = attr{
			value: value,
			kind:  kind,
		}
	}
}

// Any adds an opaque attribute to the log Record.
func Any(key string, value any) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindAny)
	}
}

// Bool adds a bool attribute to the log Record.
func Bool(key string, value bool) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindBool)
	}
}

// Context sets the context attribute for the log Record.
func Context(ctx context.Context) func(*Record) {
	return func(r *Record) {
		r.ctx = ctx
	}
}

// Duration adds a time.Duration attribute to the log Record.
func Duration(key string, value time.Duration) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, uint64(value.Nanoseconds()), kindDuration)
	}
}

// Err sets the error attribute for the log Record.
func Err(err error) func(*Record) {
	return func(r *Record) {
		r.err = err
	}
}

// Int converts an integer to a 64-bit integer and adds an attribute to the log Record with the converted value.
func Int(key string, value int) func(*Record) {
	return Int64(key, int64(value))
}

// Int64 adds a 64-bit integer attribute to the log Record.
func Int64(key string, value int64) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindInt64)
	}
}

// Float32 adds a 32-bit floating-point number attribute to the log Record.
func Float32(key string, value float32) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, math.Float32bits(value), kindFloat32)
	}
}

// Float64 adds a 64-bit floating-point number attribute to the log Record.
func Float64(key string, value float64) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, math.Float64bits(value), kindFloat64)
	}
}

// String adds a string attribute to the log Record.
func String(key string, value string) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindString)
	}
}

// Time adds a time.Time attribute to the log Record. The monotonic portion is discarded.
func Time(key string, value time.Time) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindTime)
	}
}

// Uint64 adds an unsigned 64-bit integer attribute to the log Record.
func Uint64(key string, value uint64) func(*Record) {
	return func(r *Record) {
		r.addAttr(key, value, kindUint64)
	}
}
