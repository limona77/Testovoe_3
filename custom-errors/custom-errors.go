package custom_errors

import "fmt"

var (
	ErrNotFound   = fmt.Errorf("not found")
	ErrNotAllowed = fmt.Errorf("not allowed")
)
