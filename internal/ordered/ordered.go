package ordered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type KVPair struct {
	Key   string
	Value interface{}
}

type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		values: make(map[string]interface{}),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap) EntriesIter() func() (*KVPair, bool) {
	i := 0
	return func() (*KVPair, bool) {
		if i < len(om.keys) {
			key := om.keys[i]
			i++
			return &KVPair{Key: key, Value: om.values[key]}, true
		}
		return nil, false
	}
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, key := range om.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')
		valJSON, err := json.Marshal(om.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valJSON)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	om.keys = nil
	om.values = make(map[string]interface{})

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	if err := om.parseObject(dec); err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (om *OrderedMap) parseObject(dec *json.Decoder) error {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		value, err := handleDelim(t, dec)
		if err != nil {
			return err
		}

		om.Set(key, value)
	}

	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func parseArray(dec *json.Decoder) ([]interface{}, error) {
	arr := make([]interface{}, 0)
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		value, err := handleDelim(t, dec)
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}

	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		return nil, fmt.Errorf("expect JSON array close with ']'")
	}

	return arr, nil
}

func handleDelim(t json.Token, dec *json.Decoder) (interface{}, error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			om := NewOrderedMap()
			if err := om.parseObject(dec); err != nil {
				return nil, err
			}
			return om, nil
		case '[':
			return parseArray(dec)
		default:
			return nil, fmt.Errorf("unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}
