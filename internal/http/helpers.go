package http

import "strings"

func ParseHeaders(headers string) [][]string {
	headerList := strings.Split(headers, "|")
	var kvHeaders [][]string
	for _, s := range headerList {
		st := strings.Split(s, ":")
		for i, _ := range st {
			st[i] = strings.TrimSpace(st[i])
		}
		kvHeaders = append(kvHeaders, st)
	}
	return kvHeaders
}
