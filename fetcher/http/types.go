package http

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type EventHooks struct {
	OnRequest  func(*http.Request)
	OnResponse func(*http.Response)
	OnError    func(*http.Request, error)
	OnChunk    func([]byte)

	// TODO: expose parser
	OnScraped func(*http.Response)
}
