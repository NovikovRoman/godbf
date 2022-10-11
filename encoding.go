package godbf

import (
	"fmt"

	"golang.org/x/text/encoding"
)

func codePageID(enc encoding.Encoding) byte {
	const defCodePageID = 0x57 // ANSI
	if enc == nil {
		return defCodePageID
	}

	switch fmt.Sprint(enc) {
	case "Big5":
		return 0x78
	case "ISO 8859-2":
		return 0x1b
	case "IBM Code Page 865":
		return 0x66
	case "IBM Code Page 863":
		return 0x6c
	case "IBM Code Page 852":
		return 0x87
	case "IBM Code Page 860":
		return 0x24
	case "IBM Code Page 866":
		return 0x65
	case "IBM Code Page 850":
		return 0x37
	case "Windows 874":
		return 0x7c
	case "ISO 8859-9":
		return 0x88
	case "Windows 1250":
		return 0xc8
	case "Windows 1251":
		return 0xc9
	case "Windows 1252":
		return 0x59
	case "Windows 1253":
		return 0xcb
	case "Windows 1254":
		return 0xca
	case "Windows 1257":
		return 0xcc
	case "Shift JIS":
		return 0x7b
	}

	return defCodePageID
}
