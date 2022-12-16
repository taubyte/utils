package hex

import "strconv"

/* Converts a hex string to an int64 */
func Int(hex string) (int64, error) {
	output, err := strconv.ParseInt(Strip(hex), 16, 64)
	if err != nil {
		return 0, err
	}
	return output, nil
}
