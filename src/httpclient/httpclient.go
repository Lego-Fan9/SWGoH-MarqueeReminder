package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

var Client *http.Client
var ErrMaxRetriesExceeded = errors.New("max retries exceeded")
var ErrUnknownHTTP = errors.New("unknown error making http.Do")

const (
	maxRetries     = 5
	retryDelay     = time.Second
	maxIdleConns   = 50
	fiveSeconds    = 5 * time.Second  //nolint:revive,staticcheck,nolintlint
	fifteenSeconds = 15 * time.Second //nolint:revive,staticcheck,nolintlint
	thirtySeconds  = 30 * time.Second //nolint:revive,staticcheck,nolintlint
	ninetySeconds  = 90 * time.Second //nolint:revive,staticcheck,nolintlint
)

func init() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   fiveSeconds,
			KeepAlive: thirtySeconds,
		}).DialContext,
		MaxIdleConns:        maxIdleConns,
		IdleConnTimeout:     ninetySeconds,
		TLSHandshakeTimeout: fifteenSeconds,
	}

	Client = &http.Client{
		Timeout:   thirtySeconds,
		Transport: transport,
	}
}

// Does an http call using retry logic. same syntax as net/http.Do.
func DoWithRetry(req *http.Request) (*http.Response, error) {
	var err error

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}

	for attempt := range maxRetries {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, errr := Client.Do(req)
		if errr == nil {
			return resp, nil
		}

		err = errr

		time.Sleep(retryDelay * time.Duration(attempt+1))
	}

	if err == nil {
		err = ErrMaxRetriesExceeded
	}

	return nil, err
}

// Does DoWithRetry() but without the retry logic. Same syntax as net/http.Do.
func DoWithoutRetry(req *http.Request) (*http.Response, error) {
	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnknownHTTP, err)
	}

	return resp, nil
}

// Makes http calls with retry using same syntax as net/http.Get.
func Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	return DoWithRetry(req)
}

// Makes http calls with retry using same syntax as net/http.Post.
func Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)

	return DoWithRetry(req)
}

// Does retry and rate limits for Discord. Same syntax as net/http.Do.
func Discord(req *http.Request) (*http.Response, error) {
	clone1, clone2, err := cloneRequestWithBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := discordInternal(clone1)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		resp, err = discordInternal(clone2)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func discordInternal(req *http.Request) (*http.Response, error) {
	resp, err := DoWithRetry(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return handle429(resp, req)
	}

	remaining := parseInt(resp.Header.Get("X-RateLimit-Remaining")) //nolint:canonicalheader,nolintlint
	if remaining == 0 {
		wait := parseFloat(resp.Header.Get("X-RateLimit-Reset-After")) //nolint:canonicalheader,nolintlint
		time.Sleep(time.Duration(wait * float64(time.Second)))
	}

	return resp, nil
}

func cloneRequestWithBody(req *http.Request) (*http.Request, *http.Request, error) {
	var bodyBytes []byte

	var err error

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, nil, err //nolint:wrapcheck
		}

		req.Body.Close()
	}

	makeClone := func() *http.Request {
		c := req.Clone(req.Context())
		if bodyBytes != nil {
			c.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		return c
	}

	return makeClone(), makeClone(), nil
}

func handle429(resp *http.Response, req *http.Request) (*http.Response, error) {
	var err error

	retryAfter := parseRetryAfter(resp.Header.Get("X-RateLimit-Reset-After")) //nolint:canonicalheader,nolintlint
	time.Sleep(retryAfter)

	resp, err = DoWithRetry(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		remaining := parseInt(resp.Header.Get("X-RateLimit-Remaining")) //nolint:canonicalheader,nolintlint
		if remaining == 0 {
			wait := parseFloat(resp.Header.Get("X-RateLimit-Reset-After")) //nolint:canonicalheader,nolintlint
			time.Sleep(time.Duration(wait * float64(time.Second)))
		}
	}

	return resp, err
}

// Internal function for Discord() to parse the RetryAfter header and returns a time.Duration to wait.
func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 0
	}

	v, err := strconv.ParseFloat(h, 64)
	if err != nil {
		return 0
	}

	return time.Duration(v * float64(time.Second))
}

// internal to Discord(). converts a string to a float.
func parseFloat(h string) float64 {
	if h == "" {
		return 0
	}

	v, err := strconv.ParseFloat(h, 64)
	if err != nil {
		return 0
	}

	return v
}

// internal to Discord(). converts a string to an int.
func parseInt(h string) int {
	if h == "" {
		return 0
	}

	v, err := strconv.Atoi(h)
	if err != nil {
		return 0
	}

	return v
}
