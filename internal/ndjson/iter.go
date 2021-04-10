package ndjson

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
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	// TODO figure out a hook to close the file (f.Close())
	// probably needs to be exposed via the vtab pkg as a way to register an "on end"

	return &iter{
		filePath:      absPath,
		file:          f,
		scanner:       bufio.NewScanner(f),
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
		return nil, io.EOF
	}
	err := i.scanner.Err()
	if err != nil {
		return nil, err
	}
	return i, nil
}
