package lines

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/augmentable-dev/vtab"
)

type iter struct {
	filePath      string
	file          *os.File
	scanner       *bufio.Scanner
	currentLineNo int
	delimiter     string
}

// taken from bufio/scan.go
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
func newIter(filePath string, delimiter string) (*iter, error) {
	var (
		err     error
		absPath string
		f       *os.File
	)

	if filePath != "" {
		absPath, err = filepath.Abs(filePath)
		if err != nil {
			return nil, err
		}

		f, err = os.Open(absPath)
		if err != nil {
			return nil, err
		}
		// TODO figure out a hook to close the file (f.Close())
		// probably needs to be exposed via the vtab pkg as a way to register an "on end"
	} else {
		f = os.Stdin
	}
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexAny(data, delimiter); i >= 0 {
			return i + 1, dropCR(data[0:i]), nil
		}
		if atEOF {
			return len(data), dropCR(data), nil
		}
		return 0, nil, nil
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(split)
	// TODO make the buffer size settable
	buf := make([]byte, 0, 1024*1024*512)
	scanner.Buffer(buf, 0)

	return &iter{
		filePath:      absPath,
		file:          f,
		scanner:       scanner,
		currentLineNo: 0,
		delimiter:     delimiter,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	switch c {
	case 0:
		return i.currentLineNo, nil
	case 1:
		return i.scanner.Text(), nil
	case 2:
		return i.filePath, nil
	case 3:
		return i.delimiter, nil
	}

	return nil, fmt.Errorf("unknown column")
}

func (i *iter) Next() (vtab.Row, error) {
	i.currentLineNo++
	keepGoing := i.scanner.Scan()
	if !keepGoing {
		err := i.scanner.Err()
		if err != nil {
			return nil, err
		}

		return nil, io.EOF
	}
	return i, nil
}
