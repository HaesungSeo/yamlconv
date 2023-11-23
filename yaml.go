// Package yamlconv implements utility routines for manipulating
// yaml struct which converted from YAML-encoded string.
package yamlconv

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Print prints the YAML-encoded string from the yaml struct data
// to the standard out.
// tab is used to spacing the nested yaml structures.
func Print(data interface{}, tab string) {
	print(data, "", tab)
}

func print(data interface{}, tab, ntab string) {
	switch m := data.(type) {
	case []interface{}:
		nArr := len(m)
		for i, o := range m {
			fmt.Printf("\n%sA[%d/%d]", tab, i, nArr)
			print(o, ntab, ntab+ntab)
		}
	case map[string]interface{}:
		for k, v := range m {
			fmt.Printf("\n%sK[%s]", tab, k)
			print(v, ntab, ntab+ntab)
		}
	case map[interface{}]interface{}:
		for k, v := range m {
			fmt.Printf("\n%sK[%s]", tab, k)
			print(v, ntab, ntab+ntab)
		}
	case yaml.MapSlice:
		for _, o := range m {
			fmt.Printf("\n%sM[%s]", tab, o.Key)
			print(o.Value, ntab, ntab+ntab)
		}
	case string:
		fmt.Printf(" Str{%s}", m)
	case bool:
		if m {
			fmt.Printf(" Bool{true}")
		} else {
			fmt.Printf(" Bool{false}")
		}
	case int:
		fmt.Printf(" Int{%d}", m)
	case nil:
		fmt.Printf(" {}")
	default:
		fmt.Printf(" %T{%s}", m, m)
	}
}

// MarshalJson returns the JSON encoding of the sub yaml struct data.
// keys are used to filter the match sub yaml struct.
// a key in kyes must be a form of below:
// - '[' Unsigned Integer ']', e.g. '[0]', '[10]', etc
// - any string that can be used as golang map[] key.
func MarshalJson(data interface{}, keys []string) ([]byte, error) {
	sub, err := Search(data, keys)
	if err != nil {
		return nil, err
	}

	return []byte(marshalJson(sub)), nil
}

// UnmarshalJson parses the yaml struct data and stores the result in
// the value pointed to by v.
// keys are used to filter the match sub yaml struct.
// a key in kyes must be a form of below:
// - '[' Unsigned Integer ']', e.g. '[0]', '[10]', etc
// - any string that can be used as golang map[] key.
//
// If v is nil or not a pointer, it returns an InvalidUnmarshalError.
func UnmarshalJson(data interface{}, keys []string, v any) error {
	sub, err := Search(data, keys)
	if err != nil {
		return err
	}

	text := marshalJson(sub)
	return json.Unmarshal([]byte(text), v)
}

// marshalJson returns the JSON-encoded string from the yaml struct data.
func marshalJson(data interface{}) string {
	switch m := data.(type) {
	case []interface{}:
		items := ""
		sep := ""
		for _, o := range m {
			items = items + sep + marshalJson(o)
			sep = ","
		}
		return "[" + items + "]"
	case map[string]interface{}:
		items := ""
		sep := ""
		for k, v := range m {
			items = items + sep + fmt.Sprintf("\"%s\":%s", k, marshalJson(v))
			sep = ","
		}
		return "{" + items + "}"
	case map[interface{}]interface{}:
		items := ""
		sep := ""
		for k, v := range m {
			items = items + sep + fmt.Sprintf("\"%s\":%s", k, marshalJson(v))
			sep = ","
		}
		return "{" + items + "}"
	case yaml.MapSlice:
		items := ""
		sep := ""
		for _, o := range m {
			key := marshalJson(o.Key)
			value := marshalJson(o.Value)
			if len(key) > 0 {
				if len(value) > 0 {
					items = items + sep + fmt.Sprintf("%s:%s", key, value)
				} else {
					items = items + sep + key
				}
			} else {
				if len(value) > 0 {
					items = items + sep + value
				}
			}
			sep = ","
		}
		return "{" + items + "}"
	case bool:
		if m {
			return "true"
		} else {
			return "false"
		}
	case int:
		return fmt.Sprintf("%d", m)
	case nil:
		return "{}"
	case string:
		return fmt.Sprintf("\"%s\"", m)
	default:
		return fmt.Sprintf("\"%s\"", m)
	}
}

// Search returns the match sub-struct of yaml struct data.
// keys are used to filter the match sub yaml struct.
// a key in kyes must be a form of below:
// - '[' Unsigned Integer ']', e.g. '[0]', '[10]', etc
// - any string that can be used as golang map[] key.
//
// It returns the same yaml struct, if the keys is empty.
//
// An error is returned if there are no match keys or the length
// of keys are longer than the one of nesting of yaml struct data.
func Search(data interface{}, keys []string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	if len(keys) == 0 || len(keys[0]) == 0 {
		return data, nil
	}

	// index or pattern
	idx := -1
	search := ""
	switch {
	case keys[0][0] == '[':
		n, err := fmt.Sscanf(keys[0][1:], "%d", &idx)
		if err != nil {
			return data, &InvalidIndexError{
				fmt.Errorf("invalid index: %s: %w", keys[0],
					ErrInvalidIndexError)}
		}
		if n < 0 {
			return data, &IndexOutOfRangeError{
				fmt.Errorf("index out of range: %s: %w", keys[0],
					ErrIndexOutOfRangeError)}
		}
	default:
		search = keys[0]
	}

	switch m := data.(type) {
	case []interface{}:
		if idx == -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect key[%s], but []interface{}: %w", search,
					ErrInvalidIndexError)}
		}
		if idx >= len(m) {
			return nil, &NotFoundError{
				fmt.Errorf("index %d out of len(arr) %d: %w", idx, len(m),
					ErrNotFoundError)}
		}
		return Search(m[idx], keys[1:])
	case map[string]interface{}:
		if idx != -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect index %d, but map[]interface{}: %w", idx,
					ErrInvalidIndexError)}
		}
		i, ok := m[search]
		if !ok {
			var mkeys []interface{}
			for k := range m {
				mkeys = append(mkeys, k)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
		return Search(i, keys[1:])
	case map[interface{}]interface{}:
		if idx != -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect index %d, but map[]interface{}: %w", idx,
					ErrInvalidIndexError)}
		}
		i, ok := m[search]
		if !ok {
			var mkeys []interface{}
			for k := range m {
				mkeys = append(mkeys, k)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
		return Search(i, keys[1:])
	case yaml.MapSlice:
		if idx != -1 {
			if idx >= len(m) {
				return nil, &NotFoundError{
					fmt.Errorf("index %d out of len(MapSlice) %d: %w", idx, len(m),
						ErrNotFoundError)}
			}
			return Search(m[idx].Value, keys[1:])
		} else {
			var mkeys []interface{}
			for _, v := range m {
				if search == v.Key {
					return Search(v.Value, keys[1:])
				}
				mkeys = append(mkeys, v.Key)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
	default:
		if len(keys) > 0 {
			return nil, &SearchKeyTooLongError{
				fmt.Errorf("key left: %s: %w", keys,
					ErrSearchKeyTooLongError)}
		}
		return data, nil
	}
}

// Subtract returns the modified yaml struct data.
// keys are used to filter out the matching sub yaml struct.
// a key in kyes must be a form of below:
// - '[' Unsigned Integer ']', e.g. '[0]', '[10]', etc
// - any string that can be used as golang map[] key.
//
// It returns the same yaml struct, if the keys is empty.
//
// An error is returned if there are no match keys or the length
// of keys are longer than the one of nesting of yaml struct data.
func Subtract(data interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 || len(keys[0]) == 0 {
		// the end of search
		return data, nil
	}

	// index or pattern
	idx := -1
	search := ""
	switch {
	case keys[0][0] == '[':
		n, err := fmt.Sscanf(keys[0][1:], "%d", &idx)
		if err != nil {
			return nil, &InvalidIndexError{
				fmt.Errorf("invalid index: %s: %w", keys[0],
					ErrInvalidIndexError)}
		}
		if n < 0 {
			return nil, &IndexOutOfRangeError{
				fmt.Errorf("index out of range: %s: %w", keys[0],
					ErrIndexOutOfRangeError)}
		}
	default:
		search = keys[0]
	}

	switch m := data.(type) {
	case []interface{}:
		if idx == -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect key[%s], but []interface{}: %w", search,
					ErrInvalidIndexError)}
		}
		if idx >= len(m) {
			return nil, &NotFoundError{
				fmt.Errorf("index %d out of len(arr) %d: %w", idx, len(m),
					ErrNotFoundError)}
		}
		// the final
		if len(keys) == 1 {
			if idx == 0 {
				m = (m)[idx+1:]
			} else {
				m = append((m)[:idx], (m)[idx+1:]...)
			}
		} else {
			ret, err := Subtract(m[idx], keys[1:])
			if err != nil {
				return nil, err
			}
			m[idx] = ret
		}
		return m, nil
	case map[string]interface{}:
		if idx != -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect index %d, but map[]interface{}: %w", idx,
					ErrInvalidIndexError)}
		}
		i, ok := m[search]
		if !ok {
			var mkeys []interface{}
			for k := range m {
				mkeys = append(mkeys, k)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
		// the final
		if len(keys) == 1 {
			delete(m, search)
		} else {
			ret, err := Subtract(i, keys[1:])
			if err != nil {
				return nil, err
			}
			m[search] = ret
		}
		return m, nil
	case map[interface{}]interface{}:
		if idx != -1 {
			return nil, &InvalidIndexError{
				fmt.Errorf("expect index %d, but map[]interface{}: %w", idx,
					ErrInvalidIndexError)}
		}
		i, ok := m[search]
		if !ok {
			var mkeys []interface{}
			for k := range m {
				mkeys = append(mkeys, k)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
		// the final
		if len(keys) == 1 {
			delete(m, search)
		} else {
			ret, err := Subtract(i, keys[1:])
			if err != nil {
				return nil, err
			}
			m[search] = ret
		}
		return m, nil
	case yaml.MapSlice:
		if idx != -1 {
			if idx >= len(m) {
				return nil, &NotFoundError{
					fmt.Errorf("index %d out of len(MapSlice) %d: %w", idx, len(m),
						ErrNotFoundError)}
			}
			// the final
			if len(keys) == 1 {
				if idx == 0 {
					m = (m)[idx+1:]
				} else {
					m = append((m)[:idx], (m)[idx+1:]...)
				}
			} else {
				ret, err := Subtract(m[idx].Value, keys[1:])
				if err != nil {
					return nil, err
				}
				m[idx].Value = ret
			}
			return m, nil
		} else {
			var mkeys []interface{}
			for idx, v := range m {
				if search == v.Key {
					// the final
					if len(keys) == 1 {
						if idx == 0 {
							m = (m)[idx+1:]
						} else {
							m = append((m)[:idx], (m)[idx+1:]...)
						}
					} else {
						ret, err := Subtract(v.Value, keys[1:])
						if err != nil {
							return nil, err
						}
						m[idx].Value = ret
					}
					return m, nil
				}
				mkeys = append(mkeys, v.Key)
			}
			return nil, &NotFoundError{
				fmt.Errorf("search %s not in %s: %w", search, mkeys,
					ErrNotFoundError)}
		}
	default:
		if len(keys) > 0 {
			return nil, &SearchKeyTooLongError{
				fmt.Errorf("key left: %s: %w", keys,
					ErrSearchKeyTooLongError)}
		}
		return data, nil
	}
}
