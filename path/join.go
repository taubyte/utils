package path

import (
	"path"
)

/* Like "path".Join except it takes a []string rather than ...string */
func Join(names []string) string {
	return path.Join(names...)
}
