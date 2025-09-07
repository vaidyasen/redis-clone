package main

import (
	"fmt"
	"net"
	"time"

	"redis-learning/pkg/resp"
)

func main() {
	fmt.Println("=== Testing Redis Server ===")

	// Give server a moment to start if needed
	time.Sleep(100 * time.Millisecond)

	// Connect to server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		fmt.Println("Make sure the Redis server is running with: go run cmd/server/main.go")
		return
	}
	defer conn.Close()

	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)

	fmt.Println("Connected to Redis server!")

	// Test PING
	testPing(writer, parser)

	// Test SET and GET
	testSetGet(writer, parser)

	// Test DEL
	testDel(writer, parser)

	// Test error handling
	testErrors(writer, parser)

	fmt.Println("\n=== All tests completed! ===")
}

func sendCommand(writer *resp.Writer, parser *resp.Parser, command []string) resp.Value {
	var args []resp.Value
	for _, arg := range command {
		args = append(args, resp.NewBulkString(arg))
	}

	cmd := resp.NewArray(args)
	writer.Write(cmd)

	response, _ := parser.Read()
	return response
}

func testPing(writer *resp.Writer, parser *resp.Parser) {
	fmt.Println("\nTest 1: PING command")

	// Test simple PING
	response := sendCommand(writer, parser, []string{"PING"})
	fmt.Printf("PING -> %s\n", response.Str)

	// Test PING with message
	response = sendCommand(writer, parser, []string{"PING", "Hello Redis!"})
	fmt.Printf("PING Hello Redis! -> %s\n", response.Bulk)
}

func testSetGet(writer *resp.Writer, parser *resp.Parser) {
	fmt.Println("\nTest 2: SET and GET commands")

	// SET a key
	response := sendCommand(writer, parser, []string{"SET", "name", "Alice"})
	fmt.Printf("SET name Alice -> %s\n", response.Str)

	// GET the key
	response = sendCommand(writer, parser, []string{"GET", "name"})
	fmt.Printf("GET name -> %s\n", response.Bulk)

	// GET non-existent key
	response = sendCommand(writer, parser, []string{"GET", "nonexistent"})
	if response.Null {
		fmt.Printf("GET nonexistent -> (nil)\n")
	}
}

func testDel(writer *resp.Writer, parser *resp.Parser) {
	fmt.Println("\nTest 3: DEL command")

	// SET a key first
	sendCommand(writer, parser, []string{"SET", "temp", "value"})

	// DEL the key
	response := sendCommand(writer, parser, []string{"DEL", "temp"})
	fmt.Printf("DEL temp -> %d\n", response.Num)

	// Try to GET deleted key
	response = sendCommand(writer, parser, []string{"GET", "temp"})
	if response.Null {
		fmt.Printf("GET temp (after deletion) -> (nil)\n")
	}

	// DEL non-existent key
	response = sendCommand(writer, parser, []string{"DEL", "nonexistent"})
	fmt.Printf("DEL nonexistent -> %d\n", response.Num)
}

func testErrors(writer *resp.Writer, parser *resp.Parser) {
	fmt.Println("\nTest 4: Error handling")

	// Unknown command
	response := sendCommand(writer, parser, []string{"UNKNOWN"})
	fmt.Printf("UNKNOWN -> (error) %s\n", response.Str)

	// Wrong number of arguments
	response = sendCommand(writer, parser, []string{"SET", "key"})
	fmt.Printf("SET key (missing value) -> (error) %s\n", response.Str)

	response = sendCommand(writer, parser, []string{"GET"})
	fmt.Printf("GET (missing key) -> (error) %s\n", response.Str)
}
