package term

import (
	"os"

	"golang.org/x/term"
)

func GetTermSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
		height = 24
	}
	return width, height
}
