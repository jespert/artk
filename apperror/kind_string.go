// Code generated by "stringer -type Kind ."; DO NOT EDIT.

package apperror

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownKind-0]
	_ = x[NotModifiedKind-304]
	_ = x[ValidationKind-400]
	_ = x[UnauthorizedKind-401]
	_ = x[ForbiddenKind-403]
	_ = x[NotFoundKind-404]
	_ = x[ConflictKind-409]
	_ = x[PreconditionFailedKind-412]
	_ = x[TooManyRequestsKind-429]
	_ = x[TimeoutKind-504]
}

const (
	_Kind_name_0 = "UnknownKind"
	_Kind_name_1 = "NotModifiedKind"
	_Kind_name_2 = "ValidationKindUnauthorizedKind"
	_Kind_name_3 = "ForbiddenKindNotFoundKind"
	_Kind_name_4 = "ConflictKind"
	_Kind_name_5 = "PreconditionFailedKind"
	_Kind_name_6 = "TooManyRequestsKind"
	_Kind_name_7 = "TimeoutKind"
)

var (
	_Kind_index_2 = [...]uint8{0, 14, 30}
	_Kind_index_3 = [...]uint8{0, 13, 25}
)

func (i Kind) String() string {
	switch {
	case i == 0:
		return _Kind_name_0
	case i == 304:
		return _Kind_name_1
	case 400 <= i && i <= 401:
		i -= 400
		return _Kind_name_2[_Kind_index_2[i]:_Kind_index_2[i+1]]
	case 403 <= i && i <= 404:
		i -= 403
		return _Kind_name_3[_Kind_index_3[i]:_Kind_index_3[i+1]]
	case i == 409:
		return _Kind_name_4
	case i == 412:
		return _Kind_name_5
	case i == 429:
		return _Kind_name_6
	case i == 504:
		return _Kind_name_7
	default:
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}