// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package http

import "net/http"

// Configures the [HttpFetcher] instance
type HttpFetcherOptions func(*HttpFetcher)

func WithEventHooks(hooks EventHooks) HttpFetcherOptions {
	return func(f *HttpFetcher) {
		f.eventHook = hooks
	}
}

func WithCookieJar(jar http.CookieJar) HttpFetcherOptions {
	return func(f *HttpFetcher) {
		f.cookieJar = jar
	}
}

// Configures individual HTTP Requests made by [HttpFetcher]
type RequestOptions func(*http.Request)

func WithBasicAuth(username, password string) RequestOptions {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

func WithCookies(cookies []*http.Cookie) RequestOptions {
	return func(req *http.Request) {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
}

func WithHeaders(headers map[string]string) RequestOptions {
	return func(req *http.Request) {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
}

func WithReferer(referer string) RequestOptions {
	return func(req *http.Request) {
		req.Header.Add("Referer", referer)
	}
}

func WithUserAgent(userAgent string) RequestOptions {
	return func(req *http.Request) {
		req.Header.Add("User-Agent", userAgent)
	}
}
