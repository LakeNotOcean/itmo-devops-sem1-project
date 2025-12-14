package dto

import "fmt"

type FormatType string

const (
	FormatZip     FormatType = "zip"
	FormatTar     FormatType = "tar"
	DefaultFormat FormatType = FormatZip
)

func (f *FormatType) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		*f = DefaultFormat
		return nil
	}

	format := FormatType(str)

	switch format {
	case FormatZip, FormatTar:
		*f = format
		return nil
	default:
		return fmt.Errorf("invalid format: %s, allowed values: zip, tar", str)
	}
}

func (f FormatType) String() string {
	if f == "" {
		return string(DefaultFormat)
	}
	return string(f)
}

func (f FormatType) IsValid() bool {
	switch f {
	case FormatZip, FormatTar:
		return true
	default:
		return false
	}
}
