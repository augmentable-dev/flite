package lines

import (
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestSingleLineIter(t *testing.T) {
	i, err := newIter("./testdata/single-row.json", "\n")
	if err != nil {
		t.Fatal(err)
	}

	for {
		row, err := i.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}

		lineNo, err := row.Column(0)
		if err != nil {
			t.Fatal(err)
		}
		if lineNo != 1 {
			t.Fatalf("unexpected line no")
		}

		contents, err := row.Column(1)
		if err != nil {
			t.Fatal(err)
		}
		if contents != `{"hello": "world", "some_int": 42, "some_obj": { "hello": "world again" }}` {
			t.Fatalf("unexpected json row contents")
		}

		filePath, err := row.Column(2)
		if err != nil {
			t.Fatal(err)
		}
		if !filepath.IsAbs(filePath.(string)) {
			t.Fatalf("expected absolute file path")
		}
	}
}

func TestMultiLineIter(t *testing.T) {
	filePath := "./testdata/askgit-commits.ndjson"
	i, err := newIter(filePath, "\n")
	if err != nil {
		t.Fatal(err)
	}

	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := strings.Split(string(contentBytes), "\n")
	expectedLineNo := 1
	for {
		row, err := i.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}

		lineNo, err := row.Column(0)
		if err != nil {
			t.Fatal(err)
		}
		if lineNo != expectedLineNo {
			t.Fatalf("unexpected line no")
		}
		expectedLineNo++

		line, err := row.Column(1)
		if err != nil {
			t.Fatal(err)
		}

		expectedLine := expectedLines[lineNo.(int)-1]
		if line != expectedLine {
			t.Fatalf("unexpected line contents, want: %s got: %s", expectedLine, line)
		}
	}
}
