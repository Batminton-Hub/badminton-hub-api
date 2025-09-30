package domain

import "net/http"

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
