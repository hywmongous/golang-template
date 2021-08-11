package merr

import "github.com/cockroachdb/errors"

var (
	ErrInputIsInvalid   = errors.New("input is invalid")
	ErrFailedInvocation = errors.New("invocation failed")

	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")

	ErrNotImplementedYet = errors.New("This function has not been implemented yet")

	ErrDomainConstraintsViolation = errors.New("domain constraints violation")

	ErrMaxRetriesExceeded  = errors.New("max retries exceeded")
	ErrConcurrencyConflict = errors.New("concurrency conflict")

	ErrMarshalingFailed   = errors.New("marshaling failed")
	ErrUnmarshalingFailed = errors.New("unmarshaling failed")
	ErrTechnical          = errors.New("technical")
)

func MarkAndWrapError(original, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}

func CreateInvalidInputError(funcName string, param string, err error) error {
	return MarkAndWrapError(
		err, ErrInputIsInvalid, funcName+param,
	)
}

func CreateFailedInvocation(funcName string, err error) error {
	return MarkAndWrapError(
		err, ErrFailedInvocation, funcName,
	)
}

func CreateFailedStructInvocation(structName string, funcName string, err error) error {
	return MarkAndWrapError(
		err, ErrFailedInvocation, structName+funcName,
	)
}

func CreateNotImplementedYet(funcName string) error {
	return CreateFailedInvocation(funcName, ErrNotImplementedYet)
}

func CreateNotImplementedYetStruct(structName string, funcName string) error {
	return CreateFailedStructInvocation(structName, funcName, ErrNotImplementedYet)
}
