package reader

import (
	"io/ioutil"
	"os"
)

// Read is for reading file to string
func Read(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	return string(fd), err
}
