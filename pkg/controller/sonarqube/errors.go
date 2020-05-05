package sonarqube

import (
	"fmt"
)

type ErrorReason string

const (
	ErrorReasonSpecUpdate     ErrorReason = "SpecUpdate"
	ErrorReasonResourceCreate ErrorReason = "ResourceCreate"
	ErrorReasonResourceUpdate ErrorReason = "ResourceUpdate"
	ErrorReasonUnknown        ErrorReason = "Unknown"
)

type SQError interface {
	Reason() ErrorReason
}

type Error struct {
	reason  ErrorReason
	message string
}

func (r *Error) Reason() ErrorReason {
	return r.reason
}

func (r *Error) Error() string {
	return fmt.Sprintf("%s: %s", r.reason, r.message)
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) ErrorReason {
	switch t := err.(type) {
	case SQError:
		return t.Reason()
	}
	return ErrorReasonUnknown
}
