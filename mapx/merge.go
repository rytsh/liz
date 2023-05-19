package mapx

func Merge(value map[string]interface{}, to map[string]interface{}) map[string]interface{} {
	for k := range value {
		if _, ok := to[k]; ok {
			// check if to[k] is map
			if _, ok := to[k].(map[string]interface{}); ok {
				// check if value[k] is map
				if _, ok := value[k].(map[string]interface{}); ok {
					// merge
					to[k] = Merge(value[k].(map[string]interface{}), to[k].(map[string]interface{}))
				} else {
					to[k] = value[k]
				}
			} else {
				to[k] = value[k]
			}
		} else {
			to[k] = value[k]
		}
	}

	return to
}
