package documents

import (
	"errors"
	"strings"
)

var (
	ErrInvalidOrderField = errors.New("invalid order field")

	ErrInvalidOrderDirection = errors.New("invalid order direction")

	ErrInvalidOrderValueType = errors.New("invalid order value type")
)

func normalizeOrder(order *SearchOrder) error {
	if order == nil {
		return nil
	}

	order.Field = strings.TrimSpace(order.Field)
	order.Direction = SortDirection(
		strings.ToLower(
			strings.TrimSpace(string(order.Direction)),
		),
	)
	order.ValueType = OrderValueType(
		strings.ToLower(
			strings.TrimSpace(string(order.ValueType)),
		),
	)

	if order.Field == "" {
		return ErrInvalidOrderField
	}

	switch order.Direction {
	case SortDirectionAsc, SortDirectionDesc:
	default:
		return ErrInvalidOrderDirection
	}

	switch order.ValueType {
	case OrderValueTypeString,
		OrderValueTypeNumber,
		OrderValueTypeDate:
	default:
		return ErrInvalidOrderValueType
	}

	return nil
}