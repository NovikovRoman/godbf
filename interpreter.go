package godbf

import (
	"bytes"
	"fmt"

	"golang.org/x/text/encoding"
)

// NewFromByteArray creates a DbfTable, reading it from a raw byte array, expecting the supplied encoding.
func NewFromByteArray(data []byte, enc encoding.Encoding) (table *DbfTable, err error) {
	table = new(DbfTable)
	table.useEncoding(enc)
	if err = unpackHeader(data, table); err != nil {
		return
	}

	table.dataStore = data
	expectedSize := int(table.numberOfBytesInHeader) + int(table.numberOfRecords)*int(table.lengthOfEachRecord)
	actualSize := len(data)
	if actualSize != expectedSize {
		err = fmt.Errorf("encoded content is %d bytes, but header expected %d", actualSize, expectedSize)
		return
	}

	lockSchema(table)
	return
}

func unpackHeader(s []byte, dt *DbfTable) error {
	dt.fileSignature = s[0]
	dt.SetLastUpdatedFromBytes(s[1:4])
	dt.SetNumberOfRecordsFromBytes(s[4:8])
	dt.SetNumberOfBytesInHeaderFromBytes(s[8:10])
	dt.SetLengthOfEachRecordFromBytes(s[10:12])

	return unpackFields(s, dt)
}

func unpackFields(s []byte, dt *DbfTable) (err error) {
	// create fieldMap to translate field name to index
	dt.fieldMap = make(map[string]int)

	// Number of fields in dbase table
	dt.numberOfFields = int((dt.numberOfBytesInHeader - 1 - 32) / 32)
	for i := 0; i < dt.numberOfFields; i++ {
		if err = unpackField(s, dt, i); err != nil {
			return
		}
	}
	return
}

func unpackField(s []byte, dt *DbfTable, fieldIndex int) (err error) {
	offset := (fieldIndex * 32) + 32

	var fieldName string
	if fieldName, err = deriveFieldName(s, dt, offset); err != nil {
		return
	}

	dt.fieldMap[fieldName] = fieldIndex

	switch s[offset+11] {
	case 'C':
		err = dt.AddTextField(fieldName, s[offset+16])
	case 'N':
		err = dt.AddNumberField(fieldName, s[offset+16], s[offset+17])
	case 'F':
		err = dt.AddFloatField(fieldName, s[offset+16], s[offset+17])
	case 'L':
		err = dt.AddBooleanField(fieldName)
	case 'D':
		err = dt.AddDateField(fieldName)
	}
	return
}

func deriveFieldName(s []byte, dt *DbfTable, offset int) (fieldName string, err error) {
	nameBytes := s[offset : offset+fieldNameByteLength]

	// Max usable field length is 10 bytes, where the 11th should contain the end of field marker.
	endOfFieldIndex := bytes.Index(nameBytes, []byte{endOfFieldNameMarker})
	if endOfFieldIndex == -1 {
		err = fmt.Errorf("end-of-field marker missing from field bytes, offset [%d,%d]",
			offset, offset+fieldNameByteLength)
		return
	}

	fieldName, err = dt.decodeString(string(nameBytes[:endOfFieldIndex]))
	return
}

func lockSchema(dt *DbfTable) {
	dt.schemaLocked = true // Schema changes no longer permitted
}

// New creates a new dbase table from scratch for the given character encoding
func New(enc encoding.Encoding) (table *DbfTable) {
	dt := new(DbfTable)

	// read dbase table header information
	dt.fileSignature = 0x03
	dt.RefreshLastUpdated()
	dt.numberOfRecords = 0
	dt.numberOfBytesInHeader = 32
	dt.lengthOfEachRecord = 0
	dt.fieldTerminator = 0x0D

	dt.useEncoding(enc)
	dt.createdFromScratch = true
	// create fieldMap to translate field name to index
	dt.fieldMap = make(map[string]int)
	dt.schemaLocked = false

	// Number of fields in dbase table
	dt.numberOfFields = int((dt.numberOfBytesInHeader - 1 - 32) / 32)

	s := make([]byte, dt.numberOfBytesInHeader+1) // +1 is for footer

	// set DbfTable dataStore slice that will store the complete file in memory
	dt.dataStore = s

	dt.dataStore[0] = dt.fileSignature
	dt.dataStore[1] = dt.updateYear
	dt.dataStore[2] = dt.updateMonth
	dt.dataStore[3] = dt.updateDay

	// no MDX file (index upon demand)
	dt.dataStore[28] = 0x00

	// set dbase language driver
	// Huston we have problem!
	// There is no easy way to deal with encoding issues. At least at the moment
	// I will try to find archaic encoding code defined by dbase standard (if there is any)
	// for given encoding. If none match I will go with default ANSI.
	//
	// Despite this flag in set in dbase file, I will continue to use provide encoding for
	// everything except this file encoding flag.
	//
	// Why? To make sure at least if you know the real encoding you can process text accordingly.
	dt.dataStore[29] = codePageID(enc)
	dt.updateHeader()
	return dt
}
