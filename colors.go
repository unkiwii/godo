package godo

import "fmt"

const (
	ColorDefault = "\x1b[39m"
	ColorGray    = "\x1b[90m"
	ColorRed     = "\x1b[91m"
	ColorGreen   = "\x1b[92m"
	ColorYellow  = "\x1b[93m"
	ColorBlue    = "\x1b[94m"
)

func nocolor(s string) string {
	return fmt.Sprintf("%s%s", ColorDefault, s)
}

func red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorDefault)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorDefault)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorDefault)
}

func yellow(s string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorDefault)
}
