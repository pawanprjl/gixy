package colors

const reset = "\033[0m"

func Red(s string) string    { return "\033[31m" + s + reset }
func Green(s string) string  { return "\033[32m" + s + reset }
func Yellow(s string) string { return "\033[33m" + s + reset }
func Cyan(s string) string   { return "\033[36m" + s + reset }
func Dim(s string) string    { return "\033[2m" + s + reset }
