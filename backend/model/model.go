package model

type LogEntry struct {
	Message   string
	Timestamp string //nano time. Remove last 9 digits to get the epch seconds time
	TimestampSeconds uint32

	Container struct {
		Name   string
		Id     string
		Labels map[string]string
	}
	Image struct {
		Id   string
		Name string
	}
	Fields map[string]string
}

type ServicePattern struct {
	Title      string
	Containers string
	Images     string
	Pattern    string
}

type GrokPatterns struct {
	Key string
	Val string
}