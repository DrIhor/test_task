package errors

import "errors"

var (
	NoSuchStorage error = errors.New("no sucmh storage")

	DataNotExist   error = errors.New("such data do not exist")
	WrongInputData error = errors.New("wrong input data into db")
)
