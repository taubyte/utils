package cast

import "fmt"

func Cast[T any](value any) (T, error) {
	v, ok := value.(T)
	if !ok {
		return v, fmt.Errorf("expected type %T got: %T", v, value)
	}

	return v, nil
}
