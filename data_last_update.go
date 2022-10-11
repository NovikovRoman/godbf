package godbf

import "time"

// dateOfLastUpdate holds the date of last update; in YYMMDD format where each is stored in a single byte, as per dBase.
// A consequence of this is that the lowest level of granularity supported is a whole 24-hour day.
// Because measures smaller than a day cannot be encoded, for any time.Time conversion, the 'time of day' for a given
// encoded day is assumed to be 12:00:00AM of that day.
//
// Timezones are also not supported. All date manipulation done via this struct assume the time.Local location applies.
// Callers should thus be careful to ensure that they are using time.Local as their location  when interfacing
// with the various decorator methods.
//
// The updateYear byte encodes the year number (0-255). The actual gregorian calendar year is derived by adding
// 1900 to the byte's value. Consequently, the range of years supported is [1900-2155] inclusive.
//
// The updateMonth byte encodes the 0-indexed month number (0-11).
//
// The updateDay byte encodes the 0-indexed day of the month (0-30).
type dateOfLastUpdate struct {
	updateYear  uint8 // YY + yearOffset (1900) = actual year.
	updateMonth uint8
	updateDay   uint8
}

// RefreshLastUpdated refreshes the dateOfLastUpdate to the YYMMDD byte encoding of today, assuming the local timezone.
// SetLastUpdated() is used by this method, and the same restrictions for it apply here.
func (ud *dateOfLastUpdate) RefreshLastUpdated() {
	ud.SetLastUpdated(time.Now())
}

// SetLastUpdated sets the dateOfLastUpdate to the YYMMDD byte encoding of the time.Time specified.
// See dateOfLastUpdate for the various limitations present in interpreting time.Time.
func (ud *dateOfLastUpdate) SetLastUpdated(updateTime time.Time) {
	ud.updateYear = byte(updateTime.Year() - yearOffset)
	ud.updateMonth = byte(updateTime.Month())
	ud.updateDay = byte(updateTime.Day())
}

// SetLastUpdatedFromBytes sets the dateOfLastUpdate to the YYMMDD byte encoding of the time specified.
// The 0-index is assigned to updateYear, the 1-index byte to updateMonth, and the 2-index byte to updateDay.
//
// See dateOfLastUpdate for further detail on appropriate byte values.
func (ud *dateOfLastUpdate) SetLastUpdatedFromBytes(timeBytes []byte) {
	ud.updateYear = timeBytes[0]
	ud.updateMonth = timeBytes[1]
	ud.updateDay = timeBytes[2]
}

// LastUpdated interprets the byte trio in dateOfLastUpdate, returning as close a time.Time value as possible.
// As no hours, minutes, seconds, etc.  are supported in the encoding, we assume 12:00:00AM for the return time.
// Similarly, time.Local is assumed for the location.
func (ud *dateOfLastUpdate) LastUpdated() time.Time {
	updateTime := time.Date(
		int(ud.updateYear)+yearOffset,
		time.Month(ud.updateMonth),
		int(ud.updateDay),
		0, 0, 0, 0,
		time.Local)

	return updateTime
}

// LowDefTime takes a time.Time and returns a low-definition time.Time equivalent that follows the same simplification
// approach as LastUpdated().
func (ud *dateOfLastUpdate) LowDefTime(highDefTime time.Time) time.Time {
	lowDefTime := time.Date(
		highDefTime.Year(),
		highDefTime.Month(),
		highDefTime.Day(),
		0, 0, 0, 0,
		time.Local)
	return lowDefTime
}
