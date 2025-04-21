package service

type WrongArgumentError struct{}

func (wa WrongArgumentError) Error() string {
	return "wrong arguments"
}
