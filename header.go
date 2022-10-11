package godbf

// dbase file header information
type header struct {
	fileSignature uint8 // Valid dBASE III PLUS table file (03h without a memo .DBT file; 83h with a memo)
	dateOfLastUpdate
	numberOfRecords       uint32 // Number of records in the table.
	numberOfBytesInHeader uint16 // Number of bytes in the header.
	lengthOfEachRecord    uint16 // Number of bytes in the record.

	// columns of dbase file
	fields          []FieldDescriptor
	fieldTerminator int8 // 0Dh stored as the field terminator.
}

// SetNumberOfRecordsFromBytes sets numberOfRecords from a byte array.
func (h *header) SetNumberOfRecordsFromBytes(s []byte) {
	h.numberOfRecords = uint32(s[0]) | (uint32(s[1]) << 8) | (uint32(s[2]) << 16) | (uint32(s[3]) << 24)
}

// SetNumberOfBytesInHeaderFromBytes sets numberOfBytesInHeader from a byte array.
func (h *header) SetNumberOfBytesInHeaderFromBytes(s []byte) {
	h.numberOfBytesInHeader = uint16(s[0]) | (uint16(s[1]) << 8)
}

// SetLengthOfEachRecordFromBytes sets lengthOfEachRecord from a byte array.
func (h *header) SetLengthOfEachRecordFromBytes(s []byte) {
	h.lengthOfEachRecord = uint16(s[0]) | (uint16(s[1]) << 8)
}
