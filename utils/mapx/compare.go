package mapx

import (
	"fmt"
	"reflect"
)

func IsSubset(set interface{}, subset interface{}) error {
	// if set is map
	if reflect.TypeOf(set).Kind() == reflect.Map {
		return IsMapSubset(set.(map[string]interface{}), subset.(map[string]interface{}))
	}

	// if set is slice
	if reflect.TypeOf(set).Kind() == reflect.Slice {
		return IsSliceSubset(set.([]interface{}), subset.([]interface{}))
	}

	return fmt.Errorf("set is not map or slice")
}

// IsMapSubset returns true if mapSubset is a subset of mapSet otherwise false
func IsMapSubset(mapSet map[string]interface{}, mapSubset map[string]interface{}) error {
	if len(mapSubset) > len(mapSet) {
		return fmt.Errorf("mapSubset length is greater than mapSet length")
	}

	for k := range mapSubset {
		if _, ok := mapSet[k]; !ok {
			return fmt.Errorf("key %s not found in mapSet", k)
		}

		mapSubsetType := reflect.TypeOf(mapSubset[k])
		mapSetType := reflect.TypeOf(mapSet[k])

		if mapSubsetType != mapSetType {
			return fmt.Errorf("type of key %s is not the same in mapSet and mapSubset", k)
		}

		if mapSubsetType.Comparable() && mapSetType.Comparable() {
			if mapSubset[k] != mapSet[k] {
				return fmt.Errorf("value of key %s is not equal", k)
			}

			continue
		}

		if mapSubsetType.Kind() == reflect.Map && mapSetType.Kind() == reflect.Map {
			if err := IsMapSubset(mapSet[k].(map[string]interface{}), mapSubset[k].(map[string]interface{})); err != nil {
				return fmt.Errorf("key %s: %v", k, err)
			}

			continue
		}

		if mapSubsetType.Kind() == reflect.Slice && mapSetType.Kind() == reflect.Slice {
			if err := IsSliceSubset(mapSet[k].([]interface{}), mapSubset[k].([]interface{})); err != nil {
				return fmt.Errorf("key %s: %v", k, err)
			}

			continue
		}

		return fmt.Errorf("value of key %s is not related", k)
	}

	return nil
}

// IsSliceSubset returns true if sliceSubset is a subset of sliceSet otherwise false
func IsSliceSubset(sliceSet []interface{}, sliceSubset []interface{}) error {
	for _, v := range sliceSubset {
		if err := IsSliceContains(sliceSet, v); err != nil {
			return fmt.Errorf("value %v not found in sliceSet", v)
		}
	}

	return nil
}

// IsSliceContains returns true if slice contains value otherwise false
func IsSliceContains(slice []interface{}, value interface{}) error {
	for _, v := range slice {
		// if value is comparable
		if reflect.TypeOf(value).Comparable() && reflect.TypeOf(v).Comparable() {
			if v == value {
				return nil
			}
		}

		// if value is map
		if reflect.TypeOf(value).Kind() == reflect.Map && reflect.TypeOf(v).Kind() == reflect.Map {
			if err := IsMapSubset(v.(map[string]interface{}), value.(map[string]interface{})); err == nil {
				return nil
			}
		}

		// if value is slice
		if reflect.TypeOf(value).Kind() == reflect.Slice && reflect.TypeOf(v).Kind() == reflect.Slice {
			if err := IsSliceSubset(v.([]interface{}), value.([]interface{})); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("value %v not found in slice", value)
}
