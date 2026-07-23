package ordered

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOrdered(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ordered Suite")
}

var _ = Describe("OrderedMap", func() {
	It("marshals an empty map to {}", func() {
		om := NewOrderedMap()
		b, err := json.Marshal(om)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(b)).To(Equal("{}"))
	})

	It("preserves insertion order via Set and EntriesIter", func() {
		om := NewOrderedMap()
		om.Set("z", 1)
		om.Set("a", 2)
		om.Set("m", 3)

		var keys []string
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
		}
		Expect(keys).To(Equal([]string{"z", "a", "m"}))
	})

	It("updates existing key value without changing order", func() {
		om := NewOrderedMap()
		om.Set("a", 1)
		om.Set("b", 2)
		om.Set("c", 3)
		om.Set("b", 99)

		var keys []string
		var values []interface{}
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
			values = append(values, pair.Value)
		}
		Expect(keys).To(Equal([]string{"a", "b", "c"}))
		Expect(values).To(Equal([]interface{}{1, 99, 3}))
	})

	It("unmarshals JSON preserving key order", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"z":1,"a":2,"m":3}`), om)
		Expect(err).ToNot(HaveOccurred())

		var keys []string
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
		}
		Expect(keys).To(Equal([]string{"z", "a", "m"}))
	})

	It("unmarshals nested objects as *OrderedMap", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"outer":{"b":2,"a":1}}`), om)
		Expect(err).ToNot(HaveOccurred())

		iter := om.EntriesIter()
		pair, ok := iter()
		Expect(ok).To(BeTrue())
		Expect(pair.Key).To(Equal("outer"))

		nested, isOM := pair.Value.(*OrderedMap)
		Expect(isOM).To(BeTrue())

		var nestedKeys []string
		nIter := nested.EntriesIter()
		for {
			p, ok := nIter()
			if !ok {
				break
			}
			nestedKeys = append(nestedKeys, p.Key)
		}
		Expect(nestedKeys).To(Equal([]string{"b", "a"}))
	})

	It("handles duplicate JSON keys by keeping the last value without duplicate positions", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"a":1,"b":2,"a":3}`), om)
		Expect(err).ToNot(HaveOccurred())

		var keys []string
		var values []interface{}
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
			values = append(values, pair.Value)
		}
		Expect(keys).To(Equal([]string{"a", "b"}))
		Expect(values[0].(json.Number).String()).To(Equal("3"))
	})

	It("unmarshals arrays as []interface{}", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"arr":[1,2,3]}`), om)
		Expect(err).ToNot(HaveOccurred())

		iter := om.EntriesIter()
		pair, ok := iter()
		Expect(ok).To(BeTrue())
		Expect(pair.Key).To(Equal("arr"))

		arr, isSlice := pair.Value.([]interface{})
		Expect(isSlice).To(BeTrue())
		Expect(arr).To(HaveLen(3))
		Expect(arr[0].(json.Number).String()).To(Equal("1"))
		Expect(arr[1].(json.Number).String()).To(Equal("2"))
		Expect(arr[2].(json.Number).String()).To(Equal("3"))
	})

	It("uses UseNumber for numeric values", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"n":12345678901234567890}`), om)
		Expect(err).ToNot(HaveOccurred())

		iter := om.EntriesIter()
		pair, ok := iter()
		Expect(ok).To(BeTrue())

		_, isNumber := pair.Value.(json.Number)
		Expect(isNumber).To(BeTrue())
		Expect(pair.Value.(json.Number).String()).To(Equal("12345678901234567890"))
	})

	It("round-trips JSON with preserved key order", func() {
		input := `{"z":"last","a":"first","m":{"nested_b":2,"nested_a":1},"arr":[1,2]}`
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(input), om)
		Expect(err).ToNot(HaveOccurred())

		output, err := json.Marshal(om)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(output)).To(Equal(input))
	})

	It("rejects non-object JSON input", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`[1,2,3]`), om)
		Expect(err).To(HaveOccurred())
	})

	It("unmarshals into a zero-value receiver without panic", func() {
		var om OrderedMap
		err := json.Unmarshal([]byte(`{"x":1,"y":2}`), &om)
		Expect(err).ToNot(HaveOccurred())

		var keys []string
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
		}
		Expect(keys).To(Equal([]string{"x", "y"}))
	})

	It("discards old keys when unmarshalling into a reused receiver", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"a":1,"b":2}`), om)
		Expect(err).ToNot(HaveOccurred())

		err = json.Unmarshal([]byte(`{"x":10}`), om)
		Expect(err).ToNot(HaveOccurred())

		var keys []string
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			keys = append(keys, pair.Key)
		}
		Expect(keys).To(Equal([]string{"x"}))

		output, err := json.Marshal(om)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(output)).To(Equal(`{"x":10}`))
	})

	It("supports the redaction iteration pattern", func() {
		om := NewOrderedMap()
		err := json.Unmarshal([]byte(`{"user":"alice","token":"secret123","role":"admin"}`), om)
		Expect(err).ToNot(HaveOccurred())

		redact := map[string]bool{"token": true}
		iter := om.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			if redact[pair.Key] {
				om.Set(pair.Key, "***")
			}
		}

		output, err := json.Marshal(om)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(output)).To(Equal(`{"user":"alice","token":"***","role":"admin"}`))
	})
})
