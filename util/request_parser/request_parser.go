package request_parser

import (
	"encoding/json"
	"encoding/xml"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"net/http"
)

var schemaDecoder = schema.NewDecoder()

func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// DecodeURLParam Parses url params to target
func DecodeURLParam(r *http.Request, target interface{}) error {
	schemaDecoder.IgnoreUnknownKeys(true)
	return schemaDecoder.Decode(target, r.URL.Query())
}

// DecodePayload Decode payload to target
func DecodePayload(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}

// DecodePayloadXML Decode payload to target in XML format
func DecodePayloadXML(r *http.Request, target interface{}) error {
	decoder := xml.NewDecoder(r.Body)
	return decoder.Decode(target)
}
