package common

import "strings"

type ErrSliceErrors struct {
	errs []error
}

func (e *ErrSliceErrors) Error() string {
	if e == nil {
		return ""
	}
	ss := make([]string, 0, len(e.errs))
	for _, e := range e.errs {
		ss = append(ss, e.Error())
	}
	return strings.Join(ss, "\n")
}

func (e *ErrSliceErrors) Unwrap() error {
	if len(e.errs) == 0 {
		return nil
	}

	if len(e.errs) == 1 {
		return e.errs[0]
	}

	return &ErrSliceErrors{errs: e.errs[1:]}
}

func (e *ErrSliceErrors) Append(err error) *ErrSliceErrors {
	e.errs = append(e.errs, err)
	return e
}

func ReduceErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	e := &ErrSliceErrors{}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	if len(e.errs) == 0 {
		return nil
	}
	return e
}
