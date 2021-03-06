package jsonstruct

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type JSONStruct map[string]interface{}

func New() JSONStruct {
	return JSONStruct(make(map[string]interface{}))
}

var (
	ErrValueNotFound = errors.New("Value not found")
)

func (s JSONStruct) String(dotPath string) (string, bool) {
	value, ok := s.FindElement(dotPath)
	if !ok {
		return "", false
	}

	switch value := value.(type) {
	case string:
		return value, true
	case int:
		return strconv.Itoa(value), true
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), true
	default:
		return "", false
	}
}

func (s JSONStruct) StringWithDefault(dotPath, defaultValue string) string {
	value, ok := s.String(dotPath)
	if !ok {
		return defaultValue
	}

	return value
}

func (s JSONStruct) Int(dotPath string) (int, bool) {
	value, ok := s.FindElement(dotPath)
	if !ok {
		return 0, false
	}

	switch value := value.(type) {
	case float64:
		// Parsed values are of value float64
		return int(value), true
	case int:
		// Set values may be of type int
		return value, true
	default:
		return 0, false
	}
}

func (s JSONStruct) IntWithDefault(dotPath string, defaultValue int) int {
	value, ok := s.Int(dotPath)
	if !ok {
		return defaultValue
	}

	return value
}

func (s JSONStruct) Duration(dotPath string) (time.Duration, error) {
	value, ok := s.String(dotPath)
	if !ok {
		return 0, ErrValueNotFound
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (s JSONStruct) DurationWithDefault(dotPath string, defaultValue time.Duration) (time.Duration, error) {
	value, err := s.Duration(dotPath)
	switch {
	case err == ErrValueNotFound:
		return defaultValue, nil
	case err != nil:
		return 0, err
	default:
		return value, nil
	}
}

func (s JSONStruct) List(dotPath string) ([]interface{}, bool) {
	value, ok := s.FindElement(dotPath)
	if !ok {
		return nil, false
	}

	list, ok := value.([]interface{})
	return list, ok
}

func (s JSONStruct) FindElement(dotPath string) (interface{}, bool) {
	if dotPath[0:1] != "." {
		return nil, false
	}

	keys := strings.Split(dotPath[1:], ".")

	parent := s

	var value interface{}
	for i, key := range keys {
		var ok bool
		value, ok = parent[key]
		if !ok {
			return nil, false
		}

		if i+1 < len(keys) {
			parent, ok = value.(map[string]interface{})
			if !ok {
				return nil, false
			}
		}
	}

	return value, true
}
