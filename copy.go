package jsonstruct

func (s JSONStruct) DeepCopy() JSONStruct {
	return JSONStruct(deepCopy(s))
}

func deepCopy(msi map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range msi {
		switch t := value.(type) {
		case map[string]interface{}:
			copy[key] = deepCopy(t)
		default:
			copy[key] = value
		}
	}
	return copy
}
