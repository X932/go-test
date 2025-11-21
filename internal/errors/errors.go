package custom_errors

import "errors"

var NotFoundError = errors.New("Not Found")
var InvalidParamError = errors.New("Invalid Param")
var InvalidCredentialsError = errors.New("Invalid credentials")
