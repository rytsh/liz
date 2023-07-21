package mapx

// Get returns the value of the key in the map.
//
// If the key does not exist, the second return value is false.
//
// Example:
//
//	m := map[string]interface{}{
//		"def": map[string]interface{}{
//			"abc": 1,
//			"xyz": 2,
//		},
//	}
//	v, ok := Get(m, []string{"def", "abc"})
//	// v = 1, ok = true
func Get(m map[string]interface{}, key []string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}

	if len(key) == 1 {
		v, ok := m[key[0]]
		return v, ok
	}

	v, ok := m[key[0]]
	if !ok {
		return nil, false
	}

	if m, ok = v.(map[string]interface{}); !ok {
		return nil, false
	}

	return Get(m, key[1:])
}
