package godbf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// For reference: https://en.wikipedia.org/wiki/.dbf#File_format_of_Level_5_DOS_dBASE

func TestDbfTable_New(t *testing.T) {
	tableUnderTest := New(nil)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Zero(t, tableUnderTest.Fields())
}

func TestDbfTable_AddBooleanField(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "testBool"
	err := tableUnderTest.AddBooleanField(expectedFieldName)
	require.Nil(t, err)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Len(t, tableUnderTest.Fields(), 1)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
	require.Equal(t, addedField.fieldType, Logical)
}

func TestDbfTable_AddBooleanField_TooLongGetsTruncated(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "FieldName!"
	suppliedFieldName := expectedFieldName + "shouldBeTruncated"

	err := tableUnderTest.AddBooleanField(suppliedFieldName)
	require.Nil(t, err)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
}

func TestDbfTable_AddBooleanField_SecondAttemptFails(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "FieldName!"

	err := tableUnderTest.AddBooleanField(expectedFieldName)
	require.Nil(t, err)

	err = tableUnderTest.AddBooleanField(expectedFieldName)
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_AddBooleanField_ErrorAfterDataEntryStart(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "goodField"

	err := tableUnderTest.AddBooleanField(expectedFieldName)
	require.Nil(t, err)

	_, err = tableUnderTest.AddNewRecord()
	require.Nil(t, err)

	postDataEntryField := "badField"

	err = tableUnderTest.AddBooleanField(postDataEntryField)
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_AddDateField(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "testDate"
	err := tableUnderTest.AddDateField(expectedFieldName)
	require.Nil(t, err)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Len(t, tableUnderTest.Fields(), 1)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
	require.Equal(t, addedField.fieldType, Date)
}

func TestDbfTable_AddTextField(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "testText"
	expectedFieldLength := byte(20)
	err := tableUnderTest.AddTextField(expectedFieldName, expectedFieldLength)
	require.Nil(t, err)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Len(t, tableUnderTest.Fields(), 1)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
	require.Equal(t, addedField.fieldType, Character)
	require.Equal(t, addedField.length, expectedFieldLength)
}

func TestDbfTable_AddNumericField(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "testNumber"
	expectedFieldLength := byte(20)
	expectedFDecimalPlaces := byte(2)
	err := tableUnderTest.AddNumberField(expectedFieldName, expectedFieldLength, expectedFDecimalPlaces)
	require.Nil(t, err)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Len(t, tableUnderTest.Fields(), 1)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
	require.Equal(t, addedField.fieldType, Numeric)
	require.Equal(t, addedField.length, expectedFieldLength)
	require.Equal(t, addedField.decimalPlaces, expectedFDecimalPlaces)
}

func TestDbfTable_AddFloatField(t *testing.T) {
	tableUnderTest := New(nil)
	expectedFieldName := "testFloat"
	expectedFieldLength := byte(20)
	expectedFDecimalPlaces := byte(2)
	err := tableUnderTest.AddFloatField(expectedFieldName, expectedFieldLength, expectedFDecimalPlaces)
	require.Nil(t, err)
	require.Zero(t, tableUnderTest.NumberOfRecords())
	require.Len(t, tableUnderTest.Fields(), 1)

	addedField := tableUnderTest.Fields()[0]
	require.Equal(t, addedField.name, expectedFieldName)
	require.Equal(t, addedField.fieldType, Float)
	require.Equal(t, addedField.length, expectedFieldLength)
	require.Equal(t, addedField.decimalPlaces, expectedFDecimalPlaces)
}

func TestDbfTable_FieldNames(t *testing.T) {
	tableUnderTest := New(nil)

	expectedFieldNames := []string{"first", "second"}

	for _, name := range expectedFieldNames {
		require.Nil(t, tableUnderTest.AddBooleanField(name))
	}
	fieldNamesUnderTest := tableUnderTest.FieldNames()
	require.Equal(t, fieldNamesUnderTest, expectedFieldNames)
}

func TestDbfTable_DecimalPlacesInField_ValidField(t *testing.T) {
	tableUnderTest := New(nil)

	numberFieldName := "numField"
	expectedNumberDecimalPlaces := uint8(0)
	err := tableUnderTest.AddNumberField(numberFieldName, 5, expectedNumberDecimalPlaces)
	require.Nil(t, err)

	actualNumberDecimalPlaces, err := tableUnderTest.DecimalPlacesInField(numberFieldName)
	require.Nil(t, err)
	require.Equal(t, actualNumberDecimalPlaces, expectedNumberDecimalPlaces)

	floatFieldName := "floatField"
	expectedFloatDecimalPlaces := uint8(2)
	err = tableUnderTest.AddFloatField(floatFieldName, 10, expectedFloatDecimalPlaces)
	require.Nil(t, err)

	actualFloatDecimalPlaces, err := tableUnderTest.DecimalPlacesInField(floatFieldName)
	require.Nil(t, err)
	require.Equal(t, actualFloatDecimalPlaces, expectedFloatDecimalPlaces)
}

func TestDbfTable_DecimalPlacesInField_NonExistentField(t *testing.T) {
	tableUnderTest := New(nil)

	_, err := tableUnderTest.DecimalPlacesInField("missingField")
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_DecimalPlacesInField_InvalidField(t *testing.T) {
	tableUnderTest := New(nil)

	textFieldName := "textField"
	err := tableUnderTest.AddTextField(textFieldName, 5)
	require.Nil(t, err)

	_, err = tableUnderTest.DecimalPlacesInField(textFieldName)
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_GetRowAsSlice_InitiallyEmptyStrings(t *testing.T) {
	tableUnderTest := New(nil)

	booldFieldName := "boolField"
	err := tableUnderTest.AddBooleanField(booldFieldName)
	require.Nil(t, err)

	textFieldName := "textField"
	err = tableUnderTest.AddBooleanField(textFieldName)
	require.Nil(t, err)

	dateFieldName := "dateField"
	err = tableUnderTest.AddBooleanField(dateFieldName)
	require.Nil(t, err)

	numFieldName := "numField"
	err = tableUnderTest.AddBooleanField(numFieldName)
	require.Nil(t, err)

	floatFieldName := "floatField"
	err = tableUnderTest.AddBooleanField(floatFieldName)
	require.Nil(t, err)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	fieldValues := tableUnderTest.GetRowAsSlice(recordIndex)

	for _, value := range fieldValues {
		require.Equal(t, value, "")
	}
}

func TestDbfTable_GetRowAsSlice(t *testing.T) {
	tableUnderTest := New(nil)

	boolFieldName := "boolField"
	expectedBoolFieldValue := "T"
	err := tableUnderTest.AddBooleanField(boolFieldName)
	require.Nil(t, err)

	textFieldName := "textField"
	expectedTextFieldValue := "some text"
	err = tableUnderTest.AddTextField(textFieldName, 10)
	require.Nil(t, err)

	dateFieldName := "dateField"
	expectedDateFieldValue := "20181201"
	err = tableUnderTest.AddDateField(dateFieldName)
	require.Nil(t, err)

	numFieldName := "numField"
	expectedNumFieldValue := "640"
	err = tableUnderTest.AddNumberField(numFieldName, 3, 0)
	require.Nil(t, err)

	floatFieldName := "floatField"
	expectedFloatFieldValue := "640.42"
	err = tableUnderTest.AddFloatField(floatFieldName, 6, 2)
	require.Nil(t, err)

	recordIndex, err := tableUnderTest.AddNewRecord()
	require.Nil(t, err)

	err = tableUnderTest.SetFieldValueByName(recordIndex, boolFieldName, expectedBoolFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, expectedTextFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, dateFieldName, expectedDateFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, numFieldName, expectedNumFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)
	require.Nil(t, err)

	fieldValues := tableUnderTest.GetRowAsSlice(recordIndex)
	require.Equal(t, fieldValues[0], expectedBoolFieldValue)
	require.Equal(t, fieldValues[1], expectedTextFieldValue)
	require.Equal(t, fieldValues[2], expectedDateFieldValue)
	require.Equal(t, fieldValues[3], expectedNumFieldValue)
	require.Equal(t, fieldValues[4], expectedFloatFieldValue)
}

func TestDbfTable_FieldValueByName(t *testing.T) {
	tableUnderTest := New(nil)

	boolFieldName := "boolField"
	expectedBoolFieldValue := "T"
	err := tableUnderTest.AddBooleanField(boolFieldName)
	require.Nil(t, err)

	textFieldName := "textField"
	expectedTextFieldValue := "some text"
	err = tableUnderTest.AddTextField(textFieldName, 10)
	require.Nil(t, err)

	dateFieldName := "dateField"
	expectedDateFieldValue := "20181201"
	err = tableUnderTest.AddDateField(dateFieldName)
	require.Nil(t, err)

	numFieldName := "numField"
	expectedNumFieldValue := "640"
	err = tableUnderTest.AddNumberField(numFieldName, 3, 0)
	require.Nil(t, err)

	floatFieldName := "floatField"
	expectedFloatFieldValue := "640.42"
	err = tableUnderTest.AddFloatField(floatFieldName, 6, 2)
	require.Nil(t, err)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	err = tableUnderTest.SetFieldValueByName(recordIndex, boolFieldName, expectedBoolFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, expectedTextFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, dateFieldName, expectedDateFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, numFieldName, expectedNumFieldValue)
	require.Nil(t, err)
	err = tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)
	require.Nil(t, err)

	bf, _ := tableUnderTest.FieldValueByName(recordIndex, boolFieldName)
	require.Equal(t, bf, expectedBoolFieldValue)

	tf, _ := tableUnderTest.FieldValueByName(recordIndex, textFieldName)
	require.Equal(t, tf, expectedTextFieldValue)

	df, _ := tableUnderTest.FieldValueByName(recordIndex, dateFieldName)
	require.Equal(t, df, expectedDateFieldValue)

	nf, _ := tableUnderTest.FieldValueByName(recordIndex, numFieldName)
	require.Equal(t, nf, expectedNumFieldValue)

	ff, _ := tableUnderTest.FieldValueByName(recordIndex, floatFieldName)
	require.Equal(t, ff, expectedFloatFieldValue)
}

func TestDbfTable_FieldValueByName_NonExistentField(t *testing.T) {
	tableUnderTest := New(nil)
	textFieldName := "textField"
	err := tableUnderTest.AddTextField(textFieldName, 10)
	require.Nil(t, err)

	_, err = tableUnderTest.FieldValueByName(0, "missingField")
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_SetFieldValueByName_NonExistentField(t *testing.T) {
	tableUnderTest := New(nil)
	err := tableUnderTest.SetFieldValueByName(0, "missingField", "someText")
	require.NotNil(t, err)
	t.Log(err)
}

func TestDbfTable_AddRecordWithNoFieldsDefined_Errors(t *testing.T) {
	tableUnderTest := New(nil)

	recordIndex, err := tableUnderTest.AddNewRecord()
	require.NotNil(t, err)
	require.Equal(t, recordIndex, -1)
}

func TestDbfTable_Int64FieldValueByName(t *testing.T) {
	tableUnderTest := New(nil)

	intFieldName := "intField"
	expectedIntValue := 640
	expectedIntFieldValue := fmt.Sprintf("%d", expectedIntValue)
	err := tableUnderTest.AddNumberField(intFieldName, 6, 2)
	require.Nil(t, err)

	recordIndex, err := tableUnderTest.AddNewRecord()
	require.Nil(t, err)

	err = tableUnderTest.SetFieldValueByName(recordIndex, intFieldName, expectedIntFieldValue)
	require.Nil(t, err)

	actualIntFieldValue, err := tableUnderTest.Int64FieldValueByName(recordIndex, intFieldName)
	require.Nil(t, err)
	require.EqualValues(t, actualIntFieldValue, expectedIntValue)
}

func TestDbfTable_Float64FieldValueByName(t *testing.T) {
	tableUnderTest := New(nil)

	floatFieldName := "floatField"
	expectedFloatValue := 640.42
	expectedFloatFieldValue := fmt.Sprintf("%.2f", expectedFloatValue)
	err := tableUnderTest.AddFloatField(floatFieldName, 10, 2)
	require.Nil(t, err)

	recordIndex, err := tableUnderTest.AddNewRecord()
	require.Nil(t, err)

	err = tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)
	require.Nil(t, err)

	actualFloatFieldValue, err := tableUnderTest.Float64FieldValueByName(recordIndex, floatFieldName)
	require.Nil(t, err)
	require.Equal(t, actualFloatFieldValue, expectedFloatValue)
}

func TestDbfTable_FieldDescriptor(t *testing.T) {
	tableUnderTest := New(nil)

	const (
		fieldName     = "floatField"
		fieldLength   = uint8(10)
		decimalPlaces = uint8(2)
	)

	floatFieldName := fieldName
	err := tableUnderTest.AddFloatField(floatFieldName, fieldLength, decimalPlaces)
	require.Nil(t, err)

	fieldUnderTest := tableUnderTest.Fields()[0]
	require.Equal(t, fieldUnderTest.Name(), fieldName)
	require.Equal(t, fieldUnderTest.FieldType(), Float)
	require.Equal(t, fieldUnderTest.Length(), fieldLength)
	require.Equal(t, fieldUnderTest.DecimalPlaces(), decimalPlaces)
	require.Equal(t, fieldUnderTest.Name(), fieldName)
}
