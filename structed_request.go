package main

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

type structedRequest struct {
	Proto  string
	Method string
	Host   string
	Header http.Header
}

func NewStructedRequest(r *http.Request) structedRequest {
	return structedRequest{
		Proto:  r.Proto,
		Method: r.Method,
		Host:   r.Host,
		Header: r.Header,
	}
}

func (s structedRequest) GetAttributes() []attribute.KeyValue {
	attributes := make([]attribute.KeyValue, 0, len(s.Header)+3)
	attributes = append(attributes, []attribute.KeyValue{
		attribute.String("http.proto", s.Proto),
		attribute.String("http.method", s.Method),
		attribute.String("http.host", s.Host),
	}...)
	for k, v := range s.Header {
		attributes = append(attributes, attribute.String("http.header."+k, fmt.Sprintf("%v", v)))
	}
	return attributes
}
