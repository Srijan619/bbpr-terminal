package util

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func HexColor(c tcell.Color) string {
	return fmt.Sprintf("%06x", c.Hex())
}
