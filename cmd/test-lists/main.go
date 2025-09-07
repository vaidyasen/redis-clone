package main

import (
	"fmt"
	"log"
	"net"
	"redis-learning/pkg/resp"
)

func main() {
	// Connect to Redis server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis server: %v", err)
	}
	defer conn.Close()

	writer := resp.NewWriter(conn)
	parser := resp.NewParser(conn)

	fmt.Println("=== Testing Redis List Operations ===")
	fmt.Println("Connected to Redis server!")
	fmt.Println()

	// Test LPUSH command
	fmt.Println("Test 1: LPUSH operations")
	sendCommand(writer, parser, []string{"LPUSH", "mylist", "world"})
	sendCommand(writer, parser, []string{"LPUSH", "mylist", "hello"})
	sendCommand(writer, parser, []string{"LLEN", "mylist"})
	fmt.Println()

	// Test RPUSH command
	fmt.Println("Test 2: RPUSH operations")
	sendCommand(writer, parser, []string{"RPUSH", "mylist", "from"})
	sendCommand(writer, parser, []string{"RPUSH", "mylist", "Go!"})
	sendCommand(writer, parser, []string{"LLEN", "mylist"})
	fmt.Println()

	// Test LPOP and RPOP
	fmt.Println("Test 3: POP operations")
	sendCommand(writer, parser, []string{"LPOP", "mylist"})
	sendCommand(writer, parser, []string{"RPOP", "mylist"})
	sendCommand(writer, parser, []string{"LLEN", "mylist"})
	fmt.Println()

	// Test TYPE command
	fmt.Println("Test 4: TYPE command")
	sendCommand(writer, parser, []string{"SET", "mystring", "value"})
	sendCommand(writer, parser, []string{"TYPE", "mystring"})
	sendCommand(writer, parser, []string{"TYPE", "mylist"})
	sendCommand(writer, parser, []string{"TYPE", "nonexistent"})
	fmt.Println()

	// Test wrong type operations
	fmt.Println("Test 5: Wrong type operations")
	sendCommand(writer, parser, []string{"LPUSH", "mystring", "value"}) // Should fail
	sendCommand(writer, parser, []string{"GET", "mylist"})              // Should fail
	fmt.Println()

	// Test empty list operations
	fmt.Println("Test 6: Empty list operations")
	sendCommand(writer, parser, []string{"LPOP", "mylist"}) // Pop remaining items
	sendCommand(writer, parser, []string{"LPOP", "mylist"})
	sendCommand(writer, parser, []string{"LLEN", "mylist"}) // Should be 0
	sendCommand(writer, parser, []string{"LPOP", "emptylist"}) // Should return nil
	fmt.Println()

	fmt.Println("=== All list tests completed! ===")
}

func sendCommand(writer *resp.Writer, parser *resp.Parser, args []string) {
	// Create command array
	values := make([]resp.Value, len(args))
	for i, arg := range args {
		values[i] = resp.NewBulkString(arg)
	}
	command := resp.NewArray(values)

	// Send command
	if err := writer.Write(command); err != nil {
		log.Printf("Error sending command: %v", err)
		return
	}

	// Read response
	response, err := parser.Read()
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	// Format and print response
	commandStr := fmt.Sprintf("%s", args[0])
	for i := 1; i < len(args); i++ {
		commandStr += " " + args[i]
	}
	
	responseStr := formatResponse(response)
	fmt.Printf("%s -> %s\n", commandStr, responseStr)
}

func formatResponse(value resp.Value) string {
	switch value.Type {
	case "string":
		return value.Str
	case "bulk":
		if value.Bulk == "" && value.Null {
			return "(nil)"
		}
		return value.Bulk
	case "integer":
		return fmt.Sprintf("(integer) %d", value.Num)
	case "error":
		return fmt.Sprintf("(error) %s", value.Str)
	case "array":
		if len(value.Array) == 0 {
			return "(empty array)"
		}
		result := "["
		for i, v := range value.Array {
			if i > 0 {
				result += ", "
			}
			result += formatResponse(v)
		}
		result += "]"
		return result
	default:
		return fmt.Sprintf("Unknown type: %s", value.Type)
	}
}
