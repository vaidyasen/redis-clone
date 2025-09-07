package resp

import (
	"fmt"
	"io"
)

// Writer handles RESP protocol serialization
type Writer struct {
	writer io.Writer
}

// NewWriter creates a new RESP writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

// Write serializes a Value to RESP format
func (w *Writer) Write(v Value) error {
	switch v.Type {
	case STRING:
		return w.writeSimpleString(v.Str)
	case ERROR:
		return w.writeError(v.Str)
	case INTEGER:
		return w.writeInteger(v.Num)
	case BULK:
		if v.Null {
			return w.writeNullBulkString()
		}
		return w.writeBulkString(v.Bulk)
	case ARRAY:
		if v.Null {
			return w.writeNullArray()
		}
		return w.writeArray(v.Array)
	default:
		return fmt.Errorf("unknown value type: %s", v.Type)
	}
}

// writeSimpleString writes a simple string (+OK\r\n)
func (w *Writer) writeSimpleString(s string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", s)
	return err
}

// writeError writes an error (-ERR message\r\n)
func (w *Writer) writeError(s string) error {
	_, err := fmt.Fprintf(w.writer, "-%s\r\n", s)
	return err
}

// writeInteger writes an integer (:42\r\n)
func (w *Writer) writeInteger(n int) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", n)
	return err
}

// writeBulkString writes a bulk string ($5\r\nhello\r\n)
func (w *Writer) writeBulkString(s string) error {
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(s), s)
	return err
}

// writeNullBulkString writes a null bulk string ($-1\r\n)
func (w *Writer) writeNullBulkString() error {
	_, err := fmt.Fprintf(w.writer, "$-1\r\n")
	return err
}

// writeArray writes an array (*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n)
func (w *Writer) writeArray(arr []Value) error {
	_, err := fmt.Fprintf(w.writer, "*%d\r\n", len(arr))
	if err != nil {
		return err
	}
	
	for _, val := range arr {
		if err := w.Write(val); err != nil {
			return err
		}
	}
	
	return nil
}

// writeNullArray writes a null array (*-1\r\n)
func (w *Writer) writeNullArray() error {
	_, err := fmt.Fprintf(w.writer, "*-1\r\n")
	return err
}

// Helper functions to create common Values

// NewSimpleString creates a simple string value
func NewSimpleString(s string) Value {
	return Value{Type: STRING, Str: s}
}

// NewError creates an error value
func NewError(s string) Value {
	return Value{Type: ERROR, Str: s}
}

// NewInteger creates an integer value
func NewInteger(n int) Value {
	return Value{Type: INTEGER, Num: n}
}

// NewBulkString creates a bulk string value
func NewBulkString(s string) Value {
	return Value{Type: BULK, Bulk: s}
}

// NewNullBulkString creates a null bulk string value
func NewNullBulkString() Value {
	return Value{Type: BULK, Null: true}
}

// NewArray creates an array value
func NewArray(arr []Value) Value {
	return Value{Type: ARRAY, Array: arr}
}

// NewNullArray creates a null array value
func NewNullArray() Value {
	return Value{Type: ARRAY, Null: true}
}
