package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderParser(t *testing.T) {
	header := "Content-Type: application/json |Accept: application/json|accept-encoding: application/gzip|x-api-key: 09187304915uqiyoewue90832174109732y6985132"
	expected := [][]string{{"Content-Type", "application/json"}, {"Accept", "application/json"}, {"accept-encoding", "application/gzip"}, {"x-api-key", "09187304915uqiyoewue90832174109732y6985132"}}
	parsed := ParseHeaders(header)
	for i, _ := range expected {
		if expected[i][0] != parsed[i][0] || expected[i][1] != parsed[i][1] {
			t.Fatalf("expected %s,%s got %s,%s", expected[i][0], expected[i][1], parsed[i][0], parsed[i][1])
		}
	}
}
func TestHttpRequest(t *testing.T) {
	method := "GET"
	url := "http://api.citybik.es/v2/networks"
	responseRecorder := httptest.NewRecorder()
	headers := [][]string{{"Content-Type", "application/json"}}
	req, err := HttpRequest(url, headers, method)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(httpHandler)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Fatal(err)
	}
	expected := `{"alive": true}`
	if responseRecorder.Body.String() != expected {
		t.Fatalf("received %s expected %s", responseRecorder.Body.String(), expected)
	}

}
func httpHandler(w http.ResponseWriter, r *http.Request) {
	println(r.URL.String())
	if r.URL.String() != "http://api.citybik.es/v2/networks" {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "incorrect url"`)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"response": "incorrect header"`)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
