package azuremutex

type LeaseAlreadyPresentError struct {
	Err error
}

func NewLeaseAlreadyPresentError(err error) *LeaseAlreadyPresentError {
	return &LeaseAlreadyPresentError{
		Err: err,
	}
}

func (e *LeaseAlreadyPresentError) Error() string {
	return "lease already present"
}
