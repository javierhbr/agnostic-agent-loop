package yamlish

import "fmt"

func AsMap(value any) map[string]any {
	if value == nil {
		return nil
	}
	m, _ := value.(map[string]any)
	return m
}

func AsSlice(value any) []any {
	if value == nil {
		return nil
	}
	s, _ := value.([]any)
	return s
}

func Lookup(value any, path ...string) any {
	current := value
	for _, key := range path {
		m := AsMap(current)
		if m == nil {
			return nil
		}
		current = m[key]
	}
	return current
}

func LookupString(value any, path ...string) string {
	current := Lookup(value, path...)
	s, _ := current.(string)
	return s
}

func LookupBool(value any, path ...string) bool {
	current := Lookup(value, path...)
	b, _ := current.(bool)
	return b
}

func RequireString(value any, path ...string) (string, error) {
	result := LookupString(value, path...)
	if result == "" {
		return "", fmt.Errorf("missing required string at %v", path)
	}
	return result, nil
}
