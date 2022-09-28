package util

type SoftError struct {
	Err error
}

func (t SoftError) Error() string {
	if t.Err != nil {
		return t.Err.Error()
	}
	return ""
}

func (t SoftError) Unwrap() error {
	if t.Err != nil {
		return t.Err
	}
	return nil
}
