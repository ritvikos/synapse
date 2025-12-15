# HTTP Fetcher

## Purpose

It provides an abstraction layer over `http.Client` via [`HttpClient`](./types.go) interface for pluggable HTTP client implementations. The underlying HTTP client can be configured based on the requirement, while this abstraction would expose methods relevant for crawling/scraping.

Internally, the fetcher intercepts responses to apply the following transformations:

1. [**Decompression**](./decompress.go) on response bodies encoded with gzip, brotli, zstd, or deflate, based on the `Content-Encoding` header. (can be disabled via options, if the underlying client already handles it)

2. [**Charset normalization**](./charset.go) to convert the decompressed textual response bodies to UTF-8, determined via `Content-Type` header and fallbacks to [heuristic-based detection](https://www-archive.mozilla.org/projects/intl/universalcharsetdetection) on the first 1KB of the response body.
