package util

import (
	"errors"
	"fmt"
	"log"
)

// FlatData holds data in flat structure
type FlatData struct {
	Items map[string]interface{}
}

// NewFlatData creates new NewFlatData object
func NewFlatData() *FlatData {
	fd := FlatData{Items: map[string]interface{}{}}
	return &fd
}

func (t *FlatData) String() string {
	s := ""
	for key, value := range t.Items {
		s += fmt.Sprintf("[%s]%v\n", key, value)
	}
	return s
}

// Search some desc
func (t *FlatData) Search(path string) (interface{}, error) {
	if path == "" {
		return "", errors.New("Path not found")
	}

	if value, ok := t.Items[path]; ok {
		return value, nil
	}

	return "", errors.New("Path not found")
}

// BuildFrom some desc
func (t *FlatData) BuildFrom(data map[string]interface{}) {
	for key, value := range data {
		t.recursiveBuild(key, value)
	}
}

// recursiveBuild build
func (t *FlatData) recursiveBuild(prefix string, data interface{}) {
	switch data.(type) {
	// Bool related
	case bool:
		t.Items[prefix] = data
		break
	// Int related
	case int:
		t.Items[prefix] = data
		break
	case []int:
		t.Items[prefix] = data
		break
	case map[string]int:
		for key, value := range data.(map[string]int) {
			t.recursiveBuild(prefix+"."+key, value)
		}
		break

	// Float related
	case float32:
		t.Items[prefix] = data
		break
	case []float32:
		t.Items[prefix] = data
		break
	case map[string]float32:
		for key, value := range data.(map[string]float32) {
			t.recursiveBuild(prefix+"."+key, value)
		}
		break

	case float64:
		t.Items[prefix] = data
		break
	case []float64:
		t.Items[prefix] = data
		break
	case map[string]float64:
		for key, value := range data.(map[string]float64) {
			t.recursiveBuild(prefix+"."+key, value)
		}
		break

	// String related
	case string:
		t.Items[prefix] = data
		break
	case []string:
		t.Items[prefix] = data
		break
	case map[string]string:
		for key, value := range data.(map[string]string) {
			t.recursiveBuild(prefix+"."+key, value)
		}
	case map[string][]string:
		for key, value := range data.(map[string][]string) {
			t.recursiveBuild(prefix+"."+key, value)
		}
		break

	// Interface related
	case []interface{}:
		slice := []string{}
		for _, x := range data.([]interface{}) {
			slice = append(slice, fmt.Sprintf("%v", x))
		}
		t.Items[prefix] = slice
		break
	case map[string]interface{}:
		for key, value := range data.(map[string]interface{}) {
			t.recursiveBuild(prefix+"."+key, value)
		}
		break

	// Unknown
	default:
		log.Fatal(&map[string]string{"file": "flatdata.go", "Function": "recursiveBuild", "error": "Unknown type \"" + fmt.Sprintf("%T", data) + "\"", "data": fmt.Sprintf("%+v", data)})
	}
}
