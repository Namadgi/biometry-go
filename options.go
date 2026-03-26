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

// WithUserFullname sets the default X-User-Fullname header for all requests.
func WithUserFullname(name string) Option {
	return func(c *Client) { c.userFullname = name }
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// CallOption modifies an individual HTTP request.
type CallOption func(*http.Request)

// WithFullname overrides X-User-Fullname for a single request.
func WithFullname(name string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-User-Fullname", name) }
}

// WithSessionID sets X-Session-ID for a single request.
func WithSessionID(id string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Session-ID", id) }
}

// WithActionTrigger sets X-Action-Trigger for a single request.
func WithActionTrigger(trigger string) CallOption {
	return func(r *http.Request) { r.Header.Set("X-Action-Trigger", trigger) }
}

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

// WithInhouseDocAuth controls whether to use in-house document authentication.
// Defaults to true on the server. Set to false to use IDScan.
func WithInhouseDocAuth(enabled bool) CallOption {
	return func(r *http.Request) {
		if enabled {
			r.Header.Set("X-Inhouse-Docauth", "true")
		} else {
			r.Header.Set("X-Inhouse-Docauth", "false")
		}
	}
}

// WithMRZValidation enables MRZ (Machine Readable Zone) validation in
// document authentication.
func WithMRZValidation(enabled bool) CallOption {
	return func(r *http.Request) {
		if enabled {
			r.Header.Set("X-Inhouse-MRZ", "true")
		}
	}
}

func applyCallOptions(req *http.Request, opts []CallOption) {
	for _, opt := range opts {
		opt(req)
	}
}
