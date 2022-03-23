package helper

import (
	"fmt"
)

var (
	// ErrNotFound error when record not found
	ErrNotFound = fmt.Errorf("record Not Found")

	// ErrUnableToMarshalJSON error when json payload corrupt
	ErrUnableToMarshalJSON = fmt.Errorf("json payload corrupt")

	// ErrUpdateFailed error when update fails
	ErrMailSettingError = fmt.Errorf("mail template, to mail address or mail body null")

	// ErrBadParams error when bad params passed in
	ErrBadParams = fmt.Errorf("bad params error")

	// ErrForbidden error
	ErrForbidden = fmt.Errorf("forbidden")

	// ErrNotFound error can't convert type
	ErrCanNotConvertType = fmt.Errorf("type not converted")

	// ErrNotFound error can't convert struct
	ErrCanNotConvertStruct = fmt.Errorf("struct not converted because not exits")

	// ErrNotFound error Current version not lastest
	ErrCurentVersionNotLastest = fmt.Errorf("Current version not lastest")
)