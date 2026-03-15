package jsonrpc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"platform-truth-mcp-server/internal/service"
)

type Message struct {
	JSONRPC string          `json:"jsonrpc,omitempty"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolService interface {
	ToolDefinitions() []map[string]any
	Call(ctx context.Context, name string, args map[string]any) (map[string]any, error)
}

type Server struct {
	name          string
	version       string
	service       ToolService
	initialized   bool
	clientReady   bool
	supportedVers []string
}

func NewServer(name, version string, svc ToolService) *Server {
	return &Server{
		name:          name,
		version:       version,
		service:       svc,
		supportedVers: []string{"2024-11-05", "2025-03-26", "2025-06-18"},
	}
}

func (s *Server) Serve(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	writer := bufio.NewWriter(w)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		responses, err := s.handleLine(line)
		if err != nil {
			if encErr := encoder.Encode(errorResponse(nil, -32700, err.Error())); encErr != nil {
				return encErr
			}
			if flushErr := writer.Flush(); flushErr != nil {
				return flushErr
			}
			continue
		}
		for _, response := range responses {
			if response == nil {
				continue
			}
			if err := encoder.Encode(response); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (s *Server) handleLine(line string) ([]any, error) {
	if strings.HasPrefix(line, "[") {
		var messages []Message
		if err := json.Unmarshal([]byte(line), &messages); err != nil {
			return nil, err
		}
		responses := make([]any, 0, len(messages))
		for _, message := range messages {
			response := s.handleMessage(message)
			if response != nil {
				responses = append(responses, response)
			}
		}
		if len(responses) == 0 {
			return nil, nil
		}
		return []any{responses}, nil
	}

	var message Message
	if err := json.Unmarshal([]byte(line), &message); err != nil {
		return nil, err
	}
	response := s.handleMessage(message)
	if response == nil {
		return nil, nil
	}
	return []any{response}, nil
}

func (s *Server) handleMessage(message Message) any {
	if message.JSONRPC != "" && message.JSONRPC != "2.0" {
		return errorResponse(message.ID, -32600, "jsonrpc must be 2.0")
	}
	if message.Method == "" {
		return nil
	}

	switch message.Method {
	case "initialize":
		params := parseParams(message.Params)
		version := chooseVersion(stringValue(params, "protocolVersion"), s.supportedVers)
		s.initialized = true
		return successResponse(message.ID, map[string]any{
			"protocolVersion": version,
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    s.name,
				"version": s.version,
			},
		})
	case "notifications/initialized":
		s.clientReady = true
		return nil
	case "ping":
		return successResponse(message.ID, map[string]any{})
	case "tools/list":
		if !s.initialized {
			return errorResponse(message.ID, -32002, "server not initialized")
		}
		return successResponse(message.ID, map[string]any{"tools": s.service.ToolDefinitions()})
	case "tools/call":
		if !s.initialized {
			return errorResponse(message.ID, -32002, "server not initialized")
		}
		params := parseParams(message.Params)
		name := stringValue(params, "name")
		args := mapValue(params, "arguments")
		if name == "" {
			return errorResponse(message.ID, -32602, "tools/call requires params.name")
		}
		result, err := s.service.Call(context.Background(), name, args)
		if err != nil {
			return errorResponse(message.ID, -32601, err.Error())
		}
		return successResponse(message.ID, result)
	default:
		if strings.HasPrefix(message.Method, "notifications/") {
			return nil
		}
		return errorResponse(message.ID, -32601, fmt.Sprintf("method %q not found", message.Method))
	}
}

func successResponse(id any, result any) Message {
	return Message{JSONRPC: "2.0", ID: id, Result: result}
}

func errorResponse(id any, code int, message string) Message {
	return Message{JSONRPC: "2.0", ID: id, Error: &RPCError{Code: code, Message: message}}
}

func parseParams(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var params map[string]any
	if err := json.Unmarshal(raw, &params); err != nil {
		return map[string]any{}
	}
	return params
}

func stringValue(m map[string]any, key string) string {
	value, _ := m[key].(string)
	return strings.TrimSpace(value)
}

func mapValue(m map[string]any, key string) map[string]any {
	value, _ := m[key].(map[string]any)
	if value == nil {
		return map[string]any{}
	}
	return value
}

func chooseVersion(requested string, supported []string) string {
	for _, version := range supported {
		if version == requested {
			return version
		}
	}
	return supported[len(supported)-1]
}

var _ service.Tool = service.Tool{}
