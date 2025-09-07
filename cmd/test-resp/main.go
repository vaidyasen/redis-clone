package main

import (
	"bytes"
	"fmt"
	"strings"

	"redis-learning/pkg/resp"
)

func main() {
	fmt.Println("=== Testing RESP Protocol Implementation in Go ===")
	fmt.Println()

	// Test 1: Simple String
	testSimpleString()

	// Test 2: Error
	testError()

	// Test 3: Integer
	testInteger()

	// Test 4: Bulk String
	testBulkString()

	// Test 5: Array
	testArray()

	// Test 6: Redis Command (SET key value)
	testRedisCommand()

	fmt.Println()
	fmt.Println("=== All RESP tests completed! ===")
}

func testSimpleString() {
	fmt.Println("Test 1: Simple String (+OK)")
	
	// Parse
	input := "+OK\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Printf("Parsed: Type=%s, Str=%s\n", value.Type, value.Str)

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	writer.Write(resp.NewSimpleString("OK"))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}

func testError() {
	fmt.Println("Test 2: Error (-ERR unknown command)")
	
	// Parse
	input := "-ERR unknown command\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Printf("Parsed: Type=%s, Str=%s\n", value.Type, value.Str)

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	writer.Write(resp.NewError("ERR unknown command"))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}

func testInteger() {
	fmt.Println("Test 3: Integer (:42)")
	
	// Parse
	input := ":42\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Printf("Parsed: Type=%s, Num=%d\n", value.Type, value.Num)

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	writer.Write(resp.NewInteger(42))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}

func testBulkString() {
	fmt.Println("Test 4: Bulk String ($5\\r\\nhello\\r\\n)")
	
	// Parse
	input := "$5\r\nhello\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Printf("Parsed: Type=%s, Bulk=%s\n", value.Type, value.Bulk)

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	writer.Write(resp.NewBulkString("hello"))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}

func testArray() {
	fmt.Println("Test 5: Array (*2\\r\\n$5\\r\\nhello\\r\\n$5\\r\\nworld\\r\\n)")
	
	// Parse
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Printf("Parsed: Type=%s, Array=[", value.Type)
	for i, v := range value.Array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s", v.Bulk)
	}
	fmt.Printf("]\n")

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	arr := []resp.Value{
		resp.NewBulkString("hello"),
		resp.NewBulkString("world"),
	}
	writer.Write(resp.NewArray(arr))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}

func testRedisCommand() {
	fmt.Println("Test 6: Redis Command (SET key value)")
	
	// Parse a SET command: *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	input := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	parser := resp.NewParser(strings.NewReader(input))
	value, err := parser.Read()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	
	fmt.Printf("Parsed Redis Command: ")
	for i, v := range value.Array {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%s", v.Bulk)
	}
	fmt.Printf("\n")

	// Serialize
	var buf bytes.Buffer
	writer := resp.NewWriter(&buf)
	cmd := []resp.Value{
		resp.NewBulkString("SET"),
		resp.NewBulkString("key"),
		resp.NewBulkString("value"),
	}
	writer.Write(resp.NewArray(cmd))
	fmt.Printf("Serialized: %q\n\n", buf.String())
}
