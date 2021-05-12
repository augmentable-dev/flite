package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/augmentable-dev/vtab"
)

type iter struct {
	filePath  string
	file      *os.File
	scanner   *bufio.Scanner
	index     int
	delimiter string
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

	// default is a line splitter
	scanner := bufio.NewScanner(f)

	// if a delimiter is provided, see here: https://stackoverflow.com/questions/33068644/how-a-scanner-can-be-implemented-with-a-custom-split/33069759
	if delimiter != "" {
		split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			// Return nothing if at end of file and no data passed
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			// Find the index of the delimiter
			if i := strings.Index(string(data), delimiter); i >= 0 {
				return i + 1, data[0:i], nil
			}

			// If at end of file with data return the data
			if atEOF {
				return len(data), data, nil
			}

			return
		}
		scanner.Split(split)
	}

	// TODO make the buffer size settable
	buf := make([]byte, 0, 1024*1024*512)
	scanner.Buffer(buf, 0)

	return &iter{
		filePath:  absPath,
		file:      f,
		scanner:   scanner,
		index:     -1,
		delimiter: delimiter,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	switch c {
	case 0:
		return i.index, nil
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
	i.index++
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
