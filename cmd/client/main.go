package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"redis-learning/pkg/resp"
)

func main() {
	// Connect to Redis server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to Redis server!")
	fmt.Println("Type commands like: SET key value, GET key, DEL key, PING")
	fmt.Println("Type 'quit' to exit")

	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("redis> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "quit" {
			break
		}

		// Parse command
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		// Convert to RESP array
		var args []resp.Value
		for _, part := range parts {
			args = append(args, resp.NewBulkString(part))
		}

		command := resp.NewArray(args)

		// Send command
		if err := writer.Write(command); err != nil {
			fmt.Printf("Error sending command: %v\n", err)
			break
		}

		// Read response
		response, err := parser.Read()
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			break
		}

		// Print response
		printResponse(response)
	}

	fmt.Println("Goodbye!")
}

func printResponse(value resp.Value) {
	switch value.Type {
	case "string":
		fmt.Printf("(string) %s\n", value.Str)
	case "error":
		fmt.Printf("(error) %s\n", value.Str)
	case "integer":
		fmt.Printf("(integer) %d\n", value.Num)
	case "bulk":
		if value.Null {
			fmt.Println("(nil)")
		} else {
			fmt.Printf("\"%s\"\n", value.Bulk)
		}
	case "array":
		if value.Null {
			fmt.Println("(nil)")
		} else {
			for i, v := range value.Array {
				fmt.Printf("%d) ", i+1)
				printResponse(v)
			}
		}
	}
}
