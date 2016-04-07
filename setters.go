package jsonstruct

import (
	"errors"
	"strings"
	"time"
)

func (s JSONStruct) SetString(dotPath, value string) error {
	parent, lastKey, err := s.findParent(dotPath)
	if err != nil {
		return err
	}
	parent[lastKey] = value
	return nil
}

func (s JSONStruct) SetInt(dotPath string, value int) error {
	parent, lastKey, err := s.findParent(dotPath)
	if err != nil {
		return err
	}
	parent[lastKey] = value
	return nil
}

func (s JSONStruct) SetDuration(dotPath string, value time.Duration) error {
	parent, lastKey, err := s.findParent(dotPath)
	if err != nil {
		return err
	}
	parent[lastKey] = value.String()
	return nil
}

func (s JSONStruct) SetList(dotPath string, value []interface{}) error {
	parent, lastKey, err := s.findParent(dotPath)
	if err != nil {
		return err
	}
	parent[lastKey] = value
	return nil
}

func (s JSONStruct) findParent(dotPath string) (JSONStruct, string, error) {
	if dotPath[0:1] != "." {
		return nil, "", errors.New("Only . paths are currently supported")
	}

	keys := strings.Split(dotPath[1:], ".")
	if len(keys) == 1 {
		return s, keys[0], nil
	}
	lastKey := keys[len(keys)-1]
	keys = keys[0 : len(keys)-1]

	value := s
	for _, key := range keys {
		var ok bool
		child, ok := value[key]
		if !ok {
			child = make(map[string]interface{})
			value[key] = child
		}

		value, ok = child.(map[string]interface{})
		if !ok {
			newValue := make(map[string]interface{})
			value[key] = newValue
		}
	}

	return value, lastKey, nil
}
