package resource

import (
	"encoding/json"
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type Option func([]byte) []byte

type Marshallable interface {
	MarshalJSON() ([]byte, error)
}

func errorJson(err error) []byte {
	result, _ := json.Marshal(map[string]string{
		"Error": err.Error(),
	})

	return result
}

func WithRemoveKey(key string) Option {
	return func(jsonValue []byte) []byte {
		m := make(map[string]any)
		if err := json.Unmarshal(jsonValue, &m); err != nil {
			return errorJson(err)
		}

		delete(m, key)

		result, _ := json.Marshal(m)
		return result
	}
}

func BundleEntry(resource Marshallable, id string, name string, options ...Option) (*fhir.BundleEntry, error) {
	bytes, err := resource.MarshalJSON()
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		bytes = option(bytes)
	}

	entry := &fhir.BundleEntry{
		FullUrl: util.StrPtrFmt("urn:uuid:%s", id),
		Request: &fhir.BundleEntryRequest{
			Method: fhir.HTTPVerbPOST,
			Url:    name,
		},
		Resource: bytes,
	}

	return entry, nil

}
