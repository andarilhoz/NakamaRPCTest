package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeEmptyPayload(t *testing.T) {
	var payload string = ``
	payloadRequest, err := DeserializePayload(payload)
	if err != nil {
		assert.Empty(t, payloadRequest)
		assert.Equal(t, err.Error(), "unexpected end of JSON input")
	}
}

func TestDeserializeEmptyJsonPayload(t *testing.T) {
	var payload string = `{}`
	payloadRequest, err := DeserializePayload(payload)
	if err != nil {
		t.Fatalf("Error on deserialization")
	}

	assert.Equal(t, payloadRequest.RequestType, "core", "Request type should be filled with 'core'")
	assert.Equal(t, payloadRequest.RequestVersion, "1.0.0", "Request version should be filled with '1.0.0'")
	assert.Nil(t, payloadRequest.RequestHash)
}

func TestDeserializePartialPayload(t *testing.T) {
	var payload string = `{"type":"level"}`
	payloadRequest, err := DeserializePayload(payload)
	if err != nil {
		t.Fatalf("Error on deserialization")
	}

	assert.Equal(t, payloadRequest.RequestType, "level", "Request type should remain 'level'")
	assert.Equal(t, payloadRequest.RequestVersion, "1.0.0", "Request version should be filled with '1.0.0'")
	assert.Nil(t, payloadRequest.RequestHash)
}

func TestDeserializeFullPayload(t *testing.T) {
	var payload string = `{"type":"level","version":"1.2.3","hash":"123ab321"}`
	payloadRequest, err := DeserializePayload(payload)
	if err != nil {
		t.Fatalf("Error on deserialization")
	}

	assert.Equal(t, payloadRequest.RequestType, "level", "Request type should remain 'level'")
	assert.Equal(t, payloadRequest.RequestVersion, "1.2.3", "Request version should remain '1.2.3'")
	assert.Equal(t, *payloadRequest.RequestHash, "123ab321", "Request version should remain '123ab321'")
}
