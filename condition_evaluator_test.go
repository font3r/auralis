package main

import "testing"

func TestEvaluateIntCondition(t *testing.T) {
	testCases := map[string]struct {
		condition Condition
		value     int16

		expected bool
	}{
		"1 = 1": {
			condition: Condition{sign: "=", value: 1}, value: 1,
			expected: true,
		},
		"1 = 2": {
			condition: Condition{sign: "=", value: 2}, value: 1,
			expected: false,
		},
		"1 != 1": {
			condition: Condition{sign: "!=", value: 1}, value: 1,
			expected: false,
		},
		"1 != 2": {
			condition: Condition{sign: "!=", value: 2}, value: 1,
			expected: true,
		},
		"2 < 1": {
			condition: Condition{sign: "<", value: 1}, value: 2,
			expected: false,
		},
		"1 < 3": {
			condition: Condition{sign: "<", value: 3}, value: 1,
			expected: true,
		},
		"1 > 3": {
			condition: Condition{sign: ">", value: 3}, value: 1,
			expected: false,
		},
		"5 > 3": {
			condition: Condition{sign: ">", value: 3}, value: 5,
			expected: true,
		},
		"1 >= 3": {
			condition: Condition{sign: ">=", value: 3}, value: 1,
			expected: false,
		},
		"5 >= 3": {
			condition: Condition{sign: ">=", value: 3}, value: 5,
			expected: true,
		},
		"5 >= 5": {
			condition: Condition{sign: ">=", value: 5}, value: 5,
			expected: true,
		},
		"1 <= 3": {
			condition: Condition{sign: "<=", value: 3}, value: 1,
			expected: true,
		},
		"5 <= 3": {
			condition: Condition{sign: "<=", value: 3}, value: 5,
			expected: false,
		},
		"5 <= 5": {
			condition: Condition{sign: "<=", value: 5}, value: 5,
			expected: true,
		},
	}

	for test, tC := range testCases {
		res := EvaluateIntCondition(tC.condition, tC.value)
		t.Run(test, func(t *testing.T) {
			if res != tC.expected {
				t.Errorf("\nexp %+v\ngot %+v", tC.expected, res)
			}
		})
	}
}
