package runtime

import "encoding/json"

type JSONContextReader struct {
	Inputs        []byte
	ContextInputs []byte
}

func (r *JSONContextReader) ReadInputs(v interface{}) error {
	return json.Unmarshal(r.Inputs, v)
}

func (r *JSONContextReader) ReadContextInputs(v interface{}) error {
	return json.Unmarshal(r.ContextInputs, v)
}
