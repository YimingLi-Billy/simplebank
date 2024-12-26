package api

import (
	"github.com/YimingLi-Billy/simplebank/util"
	"github.com/go-playground/validator/v10"
)

// This custom validator then needs to be registed with GIN; check api/server.go
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check if currency is supported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
