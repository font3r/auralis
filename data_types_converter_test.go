package main

import (
	"testing"
)

func TestConvertToConcreteType(t *testing.T) {
	testCases := map[string]struct {
		sourceType DataType
		value      any

		expectedRes any
		expectedErr error
	}{
		"invalid smallint convertsion": {
			sourceType:  smallint,
			value:       string("12312312312313123"),
			expectedErr: ErrSmallintTypeConversion,
		},
		"valid smallint conversion with int16 max": {
			sourceType:  smallint,
			value:       string("32767"),
			expectedRes: int16(32767),
		},
		"valid smallint conversion with int16 min": {
			sourceType:  smallint,
			value:       string("-32768"),
			expectedRes: int16(-32768),
		},
	}
	for test, tC := range testCases {
		val, err := ConvertToConcreteType(tC.sourceType, tC.value)
		t.Run(test, func(t *testing.T) {
			if err != nil && err.Error() != tC.expectedErr.Error() {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedErr, err)
			} else if val != tC.expectedRes {
				t.Errorf("\nexp %+v\ngot %+v", tC.expectedRes, val)
			}
		})
	}
}
