package fstore

import "github.com/Masterminds/sprig/v3"

func definedFuncMaps() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"sprig": sprig.GenericFuncMap(),
	}
}
