package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func TestHeaderParser(t *testing.T) {
	header := "Content-Type: application/json |Accept: application/json|accept-encoding: application/gzip|x-api-key: 09187304915uqiyoewue90832174109732y6985132"
	expected := [][]string{{"Content-Type", "application/json"}, {"Accept", "application/json"}, {"accept-encoding", "application/gzip"}, {"x-api-key", "09187304915uqiyoewue90832174109732y6985132"}}
	parsed := ParseHeaders(header)
	for i := range expected {
		if expected[i][0] != parsed[i][0] || expected[i][1] != parsed[i][1] {
			t.Fatalf("expected %s,%s got %s,%s", expected[i][0], expected[i][1], parsed[i][0], parsed[i][1])
		}
	}
}

type mockRoundTripper struct {
	f func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.f(req)
}
func newMockRoundTripper(f func(*http.Request) (*http.Response, error)) *mockRoundTripper {
	return &mockRoundTripper{f: f}
}
func TestHTTPGet(t *testing.T) {
	getFunc := NewHTTPGet()
	f := getFunc.(*get)
	url := "https://some-url.com/v1/some-endpoint.json"
	body := "OK"
	f.client.Transport = newMockRoundTripper(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != url {
			t.Fatalf("expected request URL: %s, got: %s", url, req.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}, nil
	})
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("http_get", getFunc); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT http_get($1)", url)
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}
	var res string
	err = row.Scan(&res)
	if err != nil {
		t.Fatal(err)
	}
	if res != body {
		t.Fatalf("expected response: %s, got: %s", body, res)
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
