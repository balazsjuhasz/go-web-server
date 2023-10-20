package database

import "errors"

var ErrNotExist = errors.New("resource does not exist")
var ErrAlreadyExists = errors.New("already exists")
