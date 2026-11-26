package biometry

import "net/http"

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient sets a custom *http.Client (for timeouts, transport, etc.).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// CallOption modifies an individual HTTP request.
//
// Note: identity, session, trigger, document-provider, and prefilled-video
// selection are NOT headers in API v2 — they are fields of the per-call input
// structs (they are carried in the JSON `request` part). The call options below
// only cover request metadata headers that the gateway still reads directly.
type CallOption func(*http.Request)

// WithDeviceInfo sets X-Device-Info for a single request. The value should be
// a JSON-encoded string containing device metadata (identifier, model, OS, etc.).
func WithDeviceInfo(jsonInfo string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Device-Info", jsonInfo) }
}

// WithGeoLocation sets X-Geo-Location for a single request. The value should
// be a JSON-encoded string with lat/lng/country/city. If omitted, the server
// falls back to IP-based geolocation.
func WithGeoLocation(jsonGeo string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Geo-Location", jsonGeo) }
}

// WithRequestTags sets X-Request-Tags for a single request.
// Tags are comma-separated labels for categorizing the request.
func WithRequestTags(tags string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Request-Tags", tags) }
}

// WithClientApp sets X-Client-App for a single request, identifying the
// calling application.
func WithClientApp(app string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Client-App", app) }
}

// WithClientVersion sets X-Client-Version for a single request.
func WithClientVersion(version string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Client-Version", version) }
}

func applyCallOptions(req *http.Request, opts []CallOption) {
	for _, opt := range opts {
		opt(req)
	}
}
