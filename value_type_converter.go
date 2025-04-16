package main

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrSmallintTypeConversion = AuraError{
		Message: "type smallint conversion error",
		Code:    "TYPE_CONV_ERROR",
	}
	ErrUUIDTypeConversion = AuraError{
		Message: "type UUID conversion error",
		Code:    "TYPE_CONV_ERROR",
	}
)

func ConvertToConcreteType(sourceType DataType, value any) (any, error) {
	switch sourceType {
	case smallint:
		{
			v, err := strconv.ParseInt(value.(string), 10, 16)
			if err != nil {
				return nil, ErrSmallintTypeConversion
			}

			return int16(v), nil
		}
	case varchar:
		{
			return strings.Trim(value.(string), "'"), nil
		}
	case uniqueidentifier:
		{
			v, err := uuid.Parse(value.(string))
			if err != nil {
				return nil, ErrUUIDTypeConversion
			}

			return v, nil
		}
	default:
		panic("invalid source type")
	}
}
