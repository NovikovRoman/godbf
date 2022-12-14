package godbf

import (
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding"
)

// NewFromFile creates a DbfTable, reading it from a file with the given file name, expecting the supplied encoding.
func NewFromFile(fileName string, enc encoding.Encoding) (table *DbfTable, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(fileName); err != nil {
		return
	}
	return NewFromByteArray(data, enc)
}

// Save saves the supplied DbfTable to a file of the specified filename
func (dt *DbfTable) Save(filename string, fileMode os.FileMode) error {
	return ioutil.WriteFile(filename, dt.dataStore, fileMode)
}
