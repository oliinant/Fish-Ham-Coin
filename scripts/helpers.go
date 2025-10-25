package scripts

import (
	"fmt"
)

func WrapError(context string err error) error {
	if err == nil {
		return 
	}
	return fmt.Errorf("%s: %w", context, err)
}