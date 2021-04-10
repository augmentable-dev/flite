package lines

import (
	"bufio"
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
}

func newIter(filePath string) (*iter, error) {
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

	scanner := bufio.NewScanner(f)
	// TODO make the buffer size settable
	buf := make([]byte, 0, 1024*1024*512)
	scanner.Buffer(buf, 0)

	return &iter{
		filePath:      absPath,
		file:          f,
		scanner:       scanner,
		currentLineNo: 0,
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
