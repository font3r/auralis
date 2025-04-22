package main

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type DataType string

const (
	smallint         DataType = "smallint"         // 2B
	integer          DataType = "integer"          // 4
	bigint           DataType = "bigint"           // 8
	varchar          DataType = "varchar"          // fixed 16 for now
	uniqueidentifier DataType = "uniqueidentifier" // 16
	boolean          DataType = "boolean"          // 1
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

func getDataTypeByteSize(dataType DataType) int {
	switch dataType {
	case smallint:
		return 2 // int16
	case integer:
		return 4 // int32
	case bigint:
		return 8 // int64
	case varchar:
		return 16 // TODO: this should be configurable
	case uniqueidentifier:
		return 16
	case boolean:
		return 1
	default:
		panic("unhandled type")
	}
}
