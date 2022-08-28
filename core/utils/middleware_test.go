package utils

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

type TestDecodeBodyScenario struct {
	description string
	body        io.ReadCloser
	out         interface{}
	decoder     *json.Decoder
	decodeErr   error
	noErr       bool
}

func TestDecodeBody(t *testing.T) {
	decoder := &json.Decoder{}
	scenarios := []TestDecodeBodyScenario{
		{
			description: "unmarshal err",
			body:        nil,
			out:         "out",
			decoder:     decoder,
			decodeErr:   &json.UnmarshalFieldError{},
			noErr:       false,
		},
		{
			description: "other errs",
			body:        nil,
			out:         "out",
			decoder:     decoder,
			decodeErr:   errors.New("other errs"),
			noErr:       false,
		},
		{
			description: "no err",
			body:        nil,
			out:         "out",
			decoder:     decoder,
			decodeErr:   nil,
			noErr:       true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(json.NewDecoder, func(r io.Reader) *json.Decoder {
				return decoder
			})
			patches.ApplyMethod(reflect.TypeOf(decoder), "DisallowUnknownFields", func(decoder *json.Decoder) {})
			patches.ApplyMethod(reflect.TypeOf(decoder), "Decode", func(decoder *json.Decoder, v interface{}) error {
				return scenario.decodeErr
			})
			err := DecodeBody(scenario.body, scenario.out)
			assert.Equal(t, err == nil, scenario.noErr)
		})
	}
}
