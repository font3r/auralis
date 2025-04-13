package main

type Condition struct {
	target string
	sign   string
	value  any
}

func ConvertConditionType(td TableDescriptor, cond *Condition) error {
	for _, cd := range td.columnDescriptors {
		if cond.target == cd.name {
			value, err := ConvertToConcreteType(cd.dataType, cond.value)
			if err != nil {
				return err
			}

			cond.value = value
		}
	}

	return nil
}

func GetMatchingCondition(conditions []Condition, target string) []Condition {
	return filter(conditions, func(c Condition) bool {
		return c.target == target
	})
}

func EvaluateIntCondition[V int16 | int32 | int64](cond Condition, value V) bool {
	switch cond.sign {
	case "=":
		return value == V(cond.value.(int))
	case "!=":
		return value != V(cond.value.(int))
	case ">":
		return value > V(cond.value.(int))
	case ">=":
		return value >= V(cond.value.(int))
	case "<":
		return value < V(cond.value.(int))
	case "<=":
		return value <= V(cond.value.(int))
	default:
		panic("invalid condition sign")
	}
}

func filter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
