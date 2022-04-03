package city

import (
	"errors"
)

var (
	ErrCityExistWithName     = errors.New("City already exist with same name in database")
	ErrCityExistWithCityCode = errors.New("City already exist with same city code in database")
)
