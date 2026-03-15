package yamlish

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type line struct {
	number int
	indent int
	text   string
}

type parser struct {
	lines []line
	pos   int
}

func Parse(data []byte) (any, error) {
	lines, err := tokenize(data)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return map[string]any{}, nil
	}

	p := &parser{lines: lines}
	value, err := p.parseBlock(lines[0].indent)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func tokenize(data []byte) ([]line, error) {
	var out []line
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		if strings.ContainsRune(raw, '\t') {
			return nil, fmt.Errorf("line %d: tabs are not supported", lineNo)
		}
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		withoutComment := stripComment(raw)
		if strings.TrimSpace(withoutComment) == "" {
			continue
		}
		indent := len(withoutComment) - len(strings.TrimLeft(withoutComment, " "))
		out = append(out, line{number: lineNo, indent: indent, text: strings.TrimSpace(withoutComment)})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func stripComment(s string) string {
	var b strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	for _, r := range s {
		switch r {
		case '\\':
			if inDouble && !escaped {
				escaped = true
				b.WriteRune(r)
				continue
			}
		case '\'':
			if !inDouble && !escaped {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle && !escaped {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return strings.TrimRight(b.String(), " ")
			}
		}
		b.WriteRune(r)
		if escaped && r != '\\' {
			escaped = false
		} else if escaped && r == '\\' {
			escaped = false
		}
	}
	return strings.TrimRight(b.String(), " ")
}

func (p *parser) parseBlock(indent int) (any, error) {
	if p.pos >= len(p.lines) {
		return nil, nil
	}
	if p.lines[p.pos].indent < indent {
		return nil, nil
	}
	if strings.HasPrefix(p.lines[p.pos].text, "-") {
		return p.parseSequence(indent)
	}
	return p.parseMap(indent)
}

func (p *parser) parseMap(indent int) (map[string]any, error) {
	result := map[string]any{}
	for p.pos < len(p.lines) {
		current := p.lines[p.pos]
		if current.indent < indent {
			break
		}
		if current.indent > indent {
			return nil, fmt.Errorf("line %d: unexpected indentation", current.number)
		}
		if strings.HasPrefix(current.text, "-") {
			break
		}
		key, rest, ok := splitKeyValue(current.text)
		if !ok {
			return nil, fmt.Errorf("line %d: expected key: value", current.number)
		}
		p.pos++
		if rest == "" {
			if p.pos < len(p.lines) && p.lines[p.pos].indent > current.indent {
				value, err := p.parseBlock(p.lines[p.pos].indent)
				if err != nil {
					return nil, err
				}
				result[key] = value
			} else {
				result[key] = nil
			}
			continue
		}
		value, err := parseScalar(rest)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", current.number, err)
		}
		result[key] = value
	}
	return result, nil
}

func (p *parser) parseSequence(indent int) ([]any, error) {
	var result []any
	for p.pos < len(p.lines) {
		current := p.lines[p.pos]
		if current.indent < indent {
			break
		}
		if current.indent > indent {
			return nil, fmt.Errorf("line %d: unexpected indentation in sequence", current.number)
		}
		if !strings.HasPrefix(current.text, "-") {
			break
		}

		content := strings.TrimSpace(strings.TrimPrefix(current.text, "-"))
		p.pos++
		if content == "" {
			if p.pos < len(p.lines) && p.lines[p.pos].indent > current.indent {
				value, err := p.parseBlock(p.lines[p.pos].indent)
				if err != nil {
					return nil, err
				}
				result = append(result, value)
			} else {
				result = append(result, nil)
			}
			continue
		}

		if key, rest, ok := splitKeyValue(content); ok {
			item := map[string]any{}
			if rest == "" {
				if p.pos < len(p.lines) && p.lines[p.pos].indent > current.indent {
					value, err := p.parseBlock(p.lines[p.pos].indent)
					if err != nil {
						return nil, err
					}
					item[key] = value
				} else {
					item[key] = nil
				}
			} else {
				value, err := parseScalar(rest)
				if err != nil {
					return nil, fmt.Errorf("line %d: %w", current.number, err)
				}
				item[key] = value
			}

			for p.pos < len(p.lines) {
				next := p.lines[p.pos]
				if next.indent <= current.indent {
					break
				}
				extra, err := p.parseMap(next.indent)
				if err != nil {
					return nil, err
				}
				for k, v := range extra {
					item[k] = v
				}
			}
			result = append(result, item)
			continue
		}

		value, err := parseScalar(content)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", current.number, err)
		}
		result = append(result, value)
	}
	return result, nil
}

func splitKeyValue(s string) (string, string, bool) {
	inSingle := false
	inDouble := false
	for i, r := range s {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case ':':
			if !inSingle && !inDouble {
				key := strings.TrimSpace(s[:i])
				rest := strings.TrimSpace(s[i+1:])
				if key == "" {
					return "", "", false
				}
				return key, rest, true
			}
		}
	}
	return "", "", false
}

func parseScalar(s string) (any, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	if (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) || (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) {
		if len(s) < 2 {
			return "", nil
		}
		return s[1 : len(s)-1], nil
	}
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	if s == "null" || s == "~" {
		return nil, nil
	}
	return s, nil
}
