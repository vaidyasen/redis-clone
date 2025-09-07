package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"redis-learning/pkg/resp"
)

// Server represents our Redis server
type Server struct {
	host     string
	port     string
	listener net.Listener
	db       *Database
}

// Database represents our in-memory data store
type Database struct {
	data map[string]*RedisValue
	mu   sync.RWMutex
}

// NewDatabase creates a new database instance
func NewDatabase() *Database {
	return &Database{
		data: make(map[string]*RedisValue),
	}
}

// Set stores a key-value pair
func (db *Database) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = NewStringValue(value)
}

// Get retrieves a value by key
func (db *Database) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, exists := db.data[key]
	if !exists || val.IsExpired() {
		if exists && val.IsExpired() {
			// Clean up expired key
			db.mu.RUnlock()
			db.mu.Lock()
			delete(db.data, key)
			db.mu.Unlock()
			db.mu.RLock()
		}
		return "", false
	}
	if val.Type != "string" {
		return "", false
	}
	return val.String, true
}

// GetValue retrieves a RedisValue by key
func (db *Database) GetValue(key string) (*RedisValue, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, exists := db.data[key]
	if !exists || val.IsExpired() {
		if exists && val.IsExpired() {
			// Clean up expired key
			db.mu.RUnlock()
			db.mu.Lock()
			delete(db.data, key)
			db.mu.Unlock()
			db.mu.RLock()
		}
		return nil, false
	}
	return val, true
}

// SetValue stores a RedisValue
func (db *Database) SetValue(key string, value *RedisValue) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// Del deletes a key
func (db *Database) Del(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, exists := db.data[key]
	if exists {
		delete(db.data, key)
	}
	return exists
}

// NewServer creates a new Redis server
func NewServer(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
		db:   NewDatabase(),
	}
}

// Start starts the server and listens for connections
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}
	
	s.listener = listener
	log.Printf("Redis server listening on %s", addr)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		
		// Handle each client in a separate goroutine
		go s.handleClient(conn)
	}
}

// Stop stops the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleClient handles a client connection
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	
	log.Printf("Client connected: %s", conn.RemoteAddr())
	
	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)
	
	for {
		// Read command from client
		value, err := parser.Read()
		if err != nil {
			log.Printf("Error reading from client %s: %v", conn.RemoteAddr(), err)
			return
		}
		
		// Process the command
		response := s.processCommand(value)
		
		// Send response back to client
		if err := writer.Write(response); err != nil {
			log.Printf("Error writing to client %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
}

// processCommand processes a Redis command and returns a response
func (s *Server) processCommand(value resp.Value) resp.Value {
	if value.Type != "array" || len(value.Array) == 0 {
		return resp.NewError("ERR invalid command format")
	}
	
	// Extract command and arguments
	command := value.Array[0].Bulk
	args := value.Array[1:]
	
	// Convert command to uppercase for case-insensitive matching
	switch command {
	case "PING":
		return s.handlePing(args)
	case "SET":
		return s.handleSet(args)
	case "GET":
		return s.handleGet(args)
	case "DEL":
		return s.handleDel(args)
	case "LPUSH":
		return s.handleLPush(args)
	case "RPUSH":
		return s.handleRPush(args)
	case "LPOP":
		return s.handleLPop(args)
	case "RPOP":
		return s.handleRPop(args)
	case "LLEN":
		return s.handleLLen(args)
	case "TYPE":
		return s.handleType(args)
	case "QUIT":
		return resp.NewSimpleString("OK")
	default:
		return resp.NewError(fmt.Sprintf("ERR unknown command '%s'", command))
	}
}

// handlePing handles the PING command
func (s *Server) handlePing(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewSimpleString("PONG")
	}
	if len(args) == 1 {
		return resp.NewBulkString(args[0].Bulk)
	}
	return resp.NewError("ERR wrong number of arguments for 'ping' command")
}

// handleSet handles the SET command
func (s *Server) handleSet(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'set' command")
	}
	
	key := args[0].Bulk
	value := args[1].Bulk
	
	s.db.Set(key, value)
	return resp.NewSimpleString("OK")
}

// handleGet handles the GET command
func (s *Server) handleGet(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'get' command")
	}
	
	key := args[0].Bulk
	value, exists := s.db.Get(key)
	
	if !exists {
		return resp.NewNullBulkString()
	}
	
	return resp.NewBulkString(value)
}

// handleDel handles the DEL command
func (s *Server) handleDel(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'del' command")
	}
	
	key := args[0].Bulk
	deleted := s.db.Del(key)
	
	if deleted {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

// handleLPush handles the LPUSH command
func (s *Server) handleLPush(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'lpush' command")
	}
	
	key := args[0].Bulk
	
	// Get or create list
	val, exists := s.db.GetValue(key)
	if !exists {
		val = NewListValue()
		s.db.SetValue(key, val)
	} else if val.Type != "list" {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	
	// Push all values
	for i := 1; i < len(args); i++ {
		val.ListPush(args[i].Bulk, true) // true for left push
	}
	
	return resp.NewInteger(val.ListLength())
}

// handleRPush handles the RPUSH command
func (s *Server) handleRPush(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'rpush' command")
	}
	
	key := args[0].Bulk
	
	// Get or create list
	val, exists := s.db.GetValue(key)
	if !exists {
		val = NewListValue()
		s.db.SetValue(key, val)
	} else if val.Type != "list" {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	
	// Push all values
	for i := 1; i < len(args); i++ {
		val.ListPush(args[i].Bulk, false) // false for right push
	}
	
	return resp.NewInteger(val.ListLength())
}

// handleLPop handles the LPOP command
func (s *Server) handleLPop(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'lpop' command")
	}
	
	key := args[0].Bulk
	val, exists := s.db.GetValue(key)
	
	if !exists {
		return resp.NewNullBulkString()
	}
	
	if val.Type != "list" {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	
	value, popped := val.ListPop(true) // true for left pop
	if !popped {
		return resp.NewNullBulkString()
	}
	
	// If list is empty, delete the key
	if val.ListLength() == 0 {
		s.db.Del(key)
	}
	
	return resp.NewBulkString(value)
}

// handleRPop handles the RPOP command
func (s *Server) handleRPop(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'rpop' command")
	}
	
	key := args[0].Bulk
	val, exists := s.db.GetValue(key)
	
	if !exists {
		return resp.NewNullBulkString()
	}
	
	if val.Type != "list" {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	
	value, popped := val.ListPop(false) // false for right pop
	if !popped {
		return resp.NewNullBulkString()
	}
	
	// If list is empty, delete the key
	if val.ListLength() == 0 {
		s.db.Del(key)
	}
	
	return resp.NewBulkString(value)
}

// handleLLen handles the LLEN command
func (s *Server) handleLLen(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'llen' command")
	}
	
	key := args[0].Bulk
	val, exists := s.db.GetValue(key)
	
	if !exists {
		return resp.NewInteger(0)
	}
	
	if val.Type != "list" {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	
	return resp.NewInteger(val.ListLength())
}

// handleType handles the TYPE command
func (s *Server) handleType(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'type' command")
	}
	
	key := args[0].Bulk
	val, exists := s.db.GetValue(key)
	
	if !exists {
		return resp.NewSimpleString("none")
	}
	
	return resp.NewSimpleString(val.Type)
}
