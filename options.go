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

func applyCallOptions(req *http.Request, opts []CallOption) {
	for _, opt := range opts {
		opt(req)
	}
}
