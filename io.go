package godbf

import (
	"io/fs"
	"io/ioutil"

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

func (dt *DbfTable) Save(filename string, fileMode fs.FileMode) error {
	return ioutil.WriteFile(filename, dt.dataStore, fileMode)
}
