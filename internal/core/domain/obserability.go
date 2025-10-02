package domain

import (
	"net/http"
)

// Metrics
type MetricsHttp struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

type MetricsCounter struct {
	Name string
	Help string
}

type MetricsGauge struct {
	Name string
	Help string
}

// Log
type LogInfo struct {
	Path    string `json:"path,omitempty"`
	Job     string `json:"job,omitempty"`
	Message string `json:"msg,omitempty"`
}

type LogError struct {
	Path  string `json:"path,omitempty"`
	Job   string `json:"job,omitempty"`
	Error string `json:"error,omitempty"`
	File  string `json:"file,omitempty"`
	Line  int    `json:"line,omitempty"`
}

// Trace
type TracerTag struct {
	TypeVal   string
	Attribute any
}

type TracerStatus struct {
	Code any
}
