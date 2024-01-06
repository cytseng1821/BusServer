package constant

import (
	"BusServer/config"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestParam struct {
	Method  string
	URL     string
	Body    io.Reader
	Header  http.Header
	TimeOut time.Duration
}

func Request(ctx context.Context, param RequestParam) ([]byte, int, error) {
	if param.TimeOut == 0 {
		param.TimeOut = time.Duration(config.ReadTimeOut) * time.Second
	}
	client := &http.Client{
		Timeout: param.TimeOut,
	}

	req, err := http.NewRequestWithContext(ctx, param.Method, param.URL, param.Body)
	if err != nil {
		return nil, 0, err
	}
	for name, value := range param.Header {
		req.Header.Set(name, strings.Join(value, ","))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%s]%s %s", param.Method, param.URL, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	return body, http.StatusOK, err
}
