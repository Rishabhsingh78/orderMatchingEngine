package utils

import "errors"

var (
	ErrInvalidOrder          = errors.New("invalid order")
	ErrOrderNotFound         = errors.New("order not found")
	ErrInsufficientLiquidity = errors.New("insufficient liquidity")
	ErrInvalidSymbol         = errors.New("invalid symbol")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrInvalidQuantity       = errors.New("invalid quantity")
)
