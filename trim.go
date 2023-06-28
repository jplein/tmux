package tmux

func Trim(s string) string {
	switch {
	case len(s) == 0:
		return s
	case s[len(s)-1] == '\n':
		return s[:len(s)-1]
	default:
		return s
	}
}
