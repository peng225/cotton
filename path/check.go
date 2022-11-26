package path

const (
	maxPathLength = 1024
)

func Valid(path string) bool {
	if len(path) > maxPathLength {
		return false
	}
	return true
}
