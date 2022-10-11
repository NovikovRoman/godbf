package godbf

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/charmap"
)

const validTestFile = "testdata/validFile.dbf"
const lessThanActualRecordsFile = "testdata/lessThanActualRecords.dbf"

const realFile = "testdata/122016B1.DBF"

// For reference: https://en.wikipedia.org/wiki/.dbf#File_format_of_Level_5_DOS_dBASE

func TestNewFromFile_ValidFile_NoError(t *testing.T) {
	_, err := NewFromFile(validTestFile, nil)
	require.Nil(t, err)
}

func TestNewFromFile_ValidFile_TableIsCorrect(t *testing.T) {
	tableUnderTest, _ := NewFromFile(validTestFile, nil)

	t.Logf("DbfReader:\n%#v\n", tableUnderTest)
	t.Logf("tableUnderTest.FieldNames() = %v\n", tableUnderTest.FieldNames())
	t.Logf("tableUnderTest.NumberOfRecords() = %v\n", tableUnderTest.NumberOfRecords())
	t.Logf("tableUnderTest.lengthOfEachRecord  %v\n", tableUnderTest.lengthOfEachRecord)

	verifyTableIsCorrect(tableUnderTest, t)
}

func TestNewFromByteArray_TableIsCorrect(t *testing.T) {
	rawFileBytes, err := ioutil.ReadFile(validTestFile)
	require.Nil(t, err)

	tableUnderTest, err := NewFromByteArray(rawFileBytes, nil)
	require.Nil(t, err)

	verifyTableIsCorrect(tableUnderTest, t)
}

func TestSaveToFile_LoadOfSavedIsCorrect(t *testing.T) {
	rawFileBytes, err := ioutil.ReadFile(validTestFile)
	require.Nil(t, err)

	var tableFromBytes *DbfTable
	tableFromBytes, err = NewFromByteArray(rawFileBytes, nil)
	require.Nil(t, err)

	tempFilename := filepath.Join("testdata", "tempSavedTable.dbf")
	err = ioutil.WriteFile(tempFilename, tableFromBytes.dataStore, os.ModePerm)
	require.Nil(t, err)

	tableUnderTest, err := NewFromFile(tempFilename, nil)
	require.Nil(t, err)

	verifyTableIsCorrect(tableUnderTest, t)

	err = os.Remove(tempFilename)
	require.Nil(t, err)
}

func TestSaveToFile_FromNew(t *testing.T) {
	var table *DbfTable

	table = New(nil)
	sampleTime := table.LowDefTime(time.Now())
	require.Zero(t, table.NumberOfRecords())
	require.Zero(t, table.Fields())
	require.Equal(t, table.LastUpdated(), sampleTime)

	tempFilename := filepath.Join("testdata", "tempSavedTable.dbf")
	err := table.Save(tempFilename, fs.ModePerm)
	require.Nil(t, err)

	table, err = NewFromFile(tempFilename, nil)
	require.Nil(t, err)
	require.Zero(t, table.NumberOfRecords())
	require.Zero(t, table.Fields())
	require.Equal(t, table.LastUpdated(), sampleTime)

	err = os.Remove(tempFilename)
	require.Nil(t, err)
}

func TestNewFromByteArray_EndOfFieldMarkerMissing_TableParsingError(t *testing.T) {
	rawFileBytes, err := ioutil.ReadFile(validTestFile)
	require.Nil(t, err)

	// Pad entire name byte range, including the final 11th byte, with non-terminating characters.
	const startByteOfFirstFieldName = 32
	for i := startByteOfFirstFieldName; i <= startByteOfFirstFieldName+fieldNameByteLength; i++ {
		rawFileBytes[i] = 0x41 // UTF-8 'A'
	}

	_, err = NewFromByteArray(rawFileBytes, nil)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "end-of-field marker missing from field bytes")
}

func TestNewFromFile_NewFromLessThanActualRecords_Errors(t *testing.T) {
	_, err := NewFromFile(lessThanActualRecordsFile, nil)
	require.NotNil(t, err)
}

func verifyTableIsCorrect(tableUnderTest *DbfTable, t *testing.T) {
	verifyFieldDescriptorsAreCorrect(tableUnderTest, t)
	verifyRecordsAreCorrect(tableUnderTest, t)
}

func verifyFieldDescriptorsAreCorrect(tableUnderTest *DbfTable, t *testing.T) {
	expectedFieldNumber := 5
	fields := tableUnderTest.Fields()
	require.Equal(t, len(fields), expectedFieldNumber)

	expectedFieldNames := []string{"TESTBOOL", "TESTTEXT", "TESTDATE", "TESTNUM", "TESTFLOAT"}
	require.Equal(t, tableUnderTest.FieldNames(), expectedFieldNames)

	boolField := tableUnderTest.Fields()[0]
	require.Equal(t, boolField.fieldType, Logical)
	require.EqualValues(t, boolField.length, 1)

	textField := tableUnderTest.Fields()[1]
	require.Equal(t, textField.fieldType, Character)
	require.EqualValues(t, textField.length, 10)

	dateField := tableUnderTest.Fields()[2]
	require.Equal(t, dateField.fieldType, Date)
	require.EqualValues(t, dateField.length, 8)

	numField := tableUnderTest.Fields()[3]
	require.Equal(t, numField.fieldType, Numeric)
	require.EqualValues(t, numField.length, 10)
	require.EqualValues(t, numField.decimalPlaces, 0)

	floatField := tableUnderTest.Fields()[4]
	require.Equal(t, floatField.fieldType, Float)
	require.EqualValues(t, floatField.length, 10)
	require.EqualValues(t, floatField.decimalPlaces, 2)
}

func verifyRecordsAreCorrect(tableUnderTest *DbfTable, t *testing.T) {
	expectedRecordNumber := 3
	actualRecordNumber := tableUnderTest.NumberOfRecords()
	require.Equal(t, actualRecordNumber, expectedRecordNumber)

	expectedRecordData := []string{"T", "test0", "20180101", "42", "42.01000"}
	require.Equal(t, tableUnderTest.GetRowAsSlice(0), expectedRecordData)

	expectedRecordData = []string{"F", "test1", "20180102", "43", "43.02000"}
	require.Equal(t, tableUnderTest.GetRowAsSlice(1), expectedRecordData)

	expectedRecordData = []string{"T", "test2", "20180103", "44", "44.03000"}
	require.Equal(t, tableUnderTest.GetRowAsSlice(2), expectedRecordData)
}

func TestFieldsNameCorrectDetect(t *testing.T) {
	tableUnderTest, _ := NewFromFile(realFile, charmap.CodePage866)
	fields := tableUnderTest.Fields()
	require.Len(t, fields, 18)

	expectedFieldNames := []string{"REGN", "PLAN", "NUM_SC", "A_P", "VR", "VV", "VITG", "ORA", "OVA", "OITGA", "ORP", "OVP", "OITGP", "IR", "IV", "IITG", "DT", "PRIZ"}

	require.Equal(t, tableUnderTest.FieldNames(), expectedFieldNames)
}

func TestNewFromFile_OpenErrors(t *testing.T) {
	_, err := NewFromFile("not_found"+lessThanActualRecordsFile, nil)
	require.NotNil(t, err)
	t.Log(err)
}

func TestSaveToFile_CreateErrors_Errors(t *testing.T) {
	rawFileBytes, err := ioutil.ReadFile(validTestFile)
	require.Nil(t, err)

	tableFromBytes, _ := NewFromByteArray(rawFileBytes, nil)
	tempFilename := filepath.Join("testdata_not_exists", "tempSavedTable.dbf")
	err = tableFromBytes.Save(tempFilename, fs.ModePerm)
	require.NotNil(t, err)
	t.Log(err)
}
