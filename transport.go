package biometry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// FileUpload represents a file to be uploaded via multipart.
type FileUpload struct {
	FieldName string
	FileName  string
	Content   io.Reader
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	if c.userFullname != "" {
		req.Header.Set("X-User-Fullname", c.userFullname)
	}
	return req, nil
}

// do executes a request, unwraps the stdResponse envelope, and decodes
// the "data" field into dst. Pass nil for dst to discard the response body.
func (c *Client) do(req *http.Request, dst any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, body)
	}

	if dst == nil {
		return nil
	}

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decoding envelope: %w", err)
	}
	if len(envelope.Data) > 0 {
		if err := json.Unmarshal(envelope.Data, dst); err != nil {
			return fmt.Errorf("decoding data: %w", err)
		}
	}
	return nil
}

// doJSON sends a JSON-encoded request and decodes the response data into dst.
func (c *Client) doJSON(ctx context.Context, method, path string, payload, dst any, callOpts ...CallOption) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encoding request: %w", err)
		}
		body = bytes.NewReader(b)
	}
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	applyCallOptions(req, callOpts)
	return c.do(req, dst)
}

// doMultipart sends a multipart/form-data request and decodes the response
// data into dst.
func (c *Client) doMultipart(ctx context.Context, path string, fields map[string]string, files []FileUpload, dst any, callOpts ...CallOption) error {
	req, err := c.buildMultipartRequest(ctx, path, fields, files)
	if err != nil {
		return err
	}
	applyCallOptions(req, callOpts)
	return c.do(req, dst)
}

// doMultipartRaw sends a multipart/form-data request and returns the full
// Response envelope (for endpoints that include scoring data).
func (c *Client) doMultipartRaw(ctx context.Context, path string, fields map[string]string, files []FileUpload, callOpts ...CallOption) (*Response, error) {
	req, err := c.buildMultipartRequest(ctx, path, fields, files)
	if err != nil {
		return nil, err
	}
	applyCallOptions(req, callOpts)
	return c.doRaw(req)
}

// doRaw executes a request and returns the full Response envelope.
func (c *Client) doRaw(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, body)
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	r.Header = resp.Header
	return &r, nil
}

func (c *Client) buildMultipartRequest(ctx context.Context, path string, fields map[string]string, files []FileUpload) (*http.Request, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("writing field %s: %w", k, err)
		}
	}

	for _, f := range files {
		part, err := w.CreateFormFile(f.FieldName, f.FileName)
		if err != nil {
			return nil, fmt.Errorf("creating form file %s: %w", f.FieldName, err)
		}
		if _, err := io.Copy(part, f.Content); err != nil {
			return nil, fmt.Errorf("copying file %s: %w", f.FieldName, err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}
