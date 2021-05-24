package http

import (
	"net/http"
	"strings"
)

func ParseHeaders(headers string) [][]string {
	headerList := strings.Split(headers, "|")
	var kvHeaders [][]string
	for _, s := range headerList {
		st := strings.Split(s, ":")
		for i := range st {
			st[i] = strings.TrimSpace(st[i])
		}
		kvHeaders = append(kvHeaders, st)
	}
	return kvHeaders
}
func HttpRequest(requestUrl string, headers [][]string, requestType string) (*http.Request, error) {
	request, err := http.NewRequest(requestType, requestUrl, nil)
	if err != nil {
		return nil, err
	}
	for _, header := range headers {
		request.Header.Add(header[0], header[1])
	}
	return request, nil
}
