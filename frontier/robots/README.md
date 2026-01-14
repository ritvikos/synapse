# Robots Resolver

## Purpose

It resolves `robots.txt` from host, parses it in compliance with [robots exclusion protocol](https://en.wikipedia.org/wiki/Robots.txt).

Internally, it uses [`RobotsFetcher`](./types.go) interface to retrieve raw `robots.txt`. To prevent "thundering herd" scenarios where multiple callers target the same host (while the `robots.txt` for that host isn't fetched), it uses [**request coalescing**](./robots.go) via `singleflight`. Once fetched, the rules are persisted in backend. Finally, [**Compliance**](./types.go) is enforced via [`RobotsEntry`](./types.go) object, which provides helper methods to verify path permissions and retrieve `Crawl-Delay` directives.
