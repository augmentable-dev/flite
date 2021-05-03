package http

import "testing"

func TestHeaderParser(t *testing.T) {
	header := "Content-Type: application/json |Accept: application/json"
	expected := [][]string{{"Content-Type", "application/json"}, {"Accept", "application/json"}}
	parsed := ParseHeaders(header)
	for i, _ := range expected {
		if expected[i][0] != parsed[i][0] || expected[i][1] != parsed[i][1] {
			t.Fatalf("expected %s,%s got %s,%s", expected[i][0], expected[i][1], parsed[i][0], parsed[i][1])
		}
	}
}
