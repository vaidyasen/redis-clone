package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Value represents a RESP value
type Value struct {
	Type   string
	Str    string
	Num    int
	Bulk   string
	Array  []Value
	Null   bool
}

// RESP data types
const (
	STRING  = "string"
	ERROR   = "error"  
	INTEGER = "integer"
	BULK    = "bulk"
	ARRAY   = "array"
)

// Parser handles RESP protocol parsing
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new RESP parser
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

// Read parses the next RESP value from the input
func (p *Parser) Read() (Value, error) {
	// Read the type byte
	typeByte, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case '+': // Simple String
		return p.readSimpleString()
	case '-': // Error  
		return p.readError()
	case ':': // Integer
		return p.readInteger()
	case '$': // Bulk String
		return p.readBulkString()
	case '*': // Array
		return p.readArray()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %c", typeByte)
	}
}

// readSimpleString reads a simple string (+OK\r\n)
func (p *Parser) readSimpleString() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	
	return Value{
		Type: STRING,
		Str:  line,
	}, nil
}

// readError reads an error (-ERR message\r\n)
func (p *Parser) readError() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	
	return Value{
		Type: ERROR,
		Str:  line,
	}, nil
}

// readInteger reads an integer (:42\r\n)
func (p *Parser) readInteger() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	
	num, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid integer: %s", line)
	}
	
	return Value{
		Type: INTEGER,
		Num:  num,
	}, nil
}

// readBulkString reads a bulk string ($5\r\nhello\r\n)
func (p *Parser) readBulkString() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	
	// Handle null bulk string
	if line == "-1" {
		return Value{
			Type: BULK,
			Null: true,
		}, nil
	}
	
	length, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid bulk string length: %s", line)
	}
	
	// Read the actual string
	bulk := make([]byte, length)
	_, err = io.ReadFull(p.reader, bulk)
	if err != nil {
		return Value{}, err
	}
	
	// Read the trailing \r\n
	p.reader.ReadByte() // \r
	p.reader.ReadByte() // \n
	
	return Value{
		Type: BULK,
		Bulk: string(bulk),
	}, nil
}

// readArray reads an array (*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n)
func (p *Parser) readArray() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	
	// Handle null array
	if line == "-1" {
		return Value{
			Type: ARRAY,
			Null: true,
		}, nil
	}
	
	length, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid array length: %s", line)
	}
	
	array := make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := p.Read()
		if err != nil {
			return Value{}, err
		}
		array[i] = val
	}
	
	return Value{
		Type:  ARRAY,
		Array: array,
	}, nil
}

// readLine reads a line ending with \r\n
func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	// Remove \r\n
	line = strings.TrimSuffix(line, "\r\n")
	return line, nil
}
