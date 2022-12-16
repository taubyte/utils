package path

import (
	pathUtils "path"
)

/*
Splits a path on / and cleans up empty paths

example:

/home/taubyte with strings.Split would evaluate to []string{"", "home", "taubyte"}

Split will evaluate `/home/taubyte` as []string{"home", "taubyte"}
*/
func Split(path string) []string {
	names := make([]string, 0)
	for _parent, _cur := pathUtils.Split(pathUtils.Clean(path)); _parent != "/" || _cur != ""; _parent, _cur = pathUtils.Split(pathUtils.Clean(_parent)) {
		names = append([]string{_cur}, names...)
	}
	return names
}
