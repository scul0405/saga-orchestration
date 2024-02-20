package pg_repo

import "errors"

var (
	ErrDuplicateEntry        = errors.New("duplicate entry")
	ErrInvalidIdempotency    = errors.New("invalid idempotency")
	ErrInsufficientInventory = errors.New("insufficient inventory")
)
