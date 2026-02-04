package encoding

import (
	"gopkg.in/yaml.v3"
)

// ToonEncoder wraps the encoding logic.
// In a full implementation, this would import github.com/toon-format/toon-go.
// For this MVP, we will use YAML as a stand-in for TOON to ensure compilation
// without external dependencies that might be missing.
type ToonEncoder struct{}

func NewToonEncoder() *ToonEncoder {
	return &ToonEncoder{}
}

func (e *ToonEncoder) Encode(v interface{}) ([]byte, error) {
	// Mock TOON encoding using YAML
	return yaml.Marshal(v)
}

func (e *ToonEncoder) Decode(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
