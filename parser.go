package gin_dump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	StringMaxLength = 0
	Newline         = "\n"
	Indent          = 4
)

func FormatJsonBytes(data []byte, hiddenFields []string, compactInArray bool) ([]byte, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	v = removeHiddenFields(v, hiddenFields)

	return []byte(format(v, 1, compactInArray)), nil
}

func FormatToJson(v any, hiddenFields []string, compactInArray bool) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return FormatJsonBytes(data, hiddenFields, compactInArray)
}

func format(v any, depth int, compactInArray bool) string {
	switch val := v.(type) {
	case string:
		return formatString(val)
	case float64:
		return fmt.Sprint(strconv.FormatFloat(val, 'f', -1, 64))
	case bool:
		return fmt.Sprint(strconv.FormatBool(val))
	case nil:
		return fmt.Sprint("null")
	case map[string]any:
		return formatMap(val, depth, compactInArray)
	case []any:
		return formatArray(val, depth, compactInArray)
	}

	return ""
}

func formatString(s string) string {
	r := []rune(s)
	if StringMaxLength > 0 && len(r) >= StringMaxLength {
		s = string(r[0:StringMaxLength]) + "..."
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(s)
	s = string(buf.Bytes())
	s = strings.TrimSuffix(s, "\n")

	return fmt.Sprint(s)
}

func formatMap(m map[string]any, depth int, compactInArray bool) string {
	if len(m) == 0 {
		return "{}"
	}

	currentIndent := generateIndent(depth - 1)
	nextIndent := generateIndent(depth)
	rows := make([]string, 0)
	keys := make([]string, 0)

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		k := fmt.Sprintf(`"%s"`, key)
		v := format(val, depth+1, compactInArray)

		valueIndent := " "
		if Newline == "" {
			valueIndent = ""
		}
		row := fmt.Sprintf("%s%s:%s%s", nextIndent, k, valueIndent, v)
		rows = append(rows, row)
	}

	return fmt.Sprintf("{%s%s%s%s}", Newline, strings.Join(rows, ","+Newline), Newline, currentIndent)
}

func formatArray(a []any, depth int, compact bool) string {
	if len(a) == 0 {
		return "[]"
	}

	if compact {

		elems := make([]string, 0)
		for _, val := range a {
			elem := format(val, depth+1, compact)
			elems = append(elems, elem)
		}
		return fmt.Sprintf("[ %s ]", strings.Join(elems, ", "))

	} else {

		currentIndent := generateIndent(depth - 1)
		nextIndent := generateIndent(depth)
		rows := make([]string, 0)

		for _, val := range a {
			c := format(val, depth+1, compact)
			row := nextIndent + c
			rows = append(rows, row)
		}
		return fmt.Sprintf("[%s%s%s%s]", Newline, strings.Join(rows, ","+Newline), Newline, currentIndent)

	}
}

func generateIndent(depth int) string {
	return strings.Repeat(" ", Indent*depth)
}

func removeHiddenFields(v any, hiddenFields []string) any {
	if _, ok := v.(map[string]any); !ok {
		return v
	}

	m := v.(map[string]any)

	for _, hiddenField := range hiddenFields {
		for k, _ := range m {
			// case-insensitive key deletion
			if strings.ToLower(k) == strings.ToLower(hiddenField) {
				delete(m, k)
			}
		}
	}

	return m
}
