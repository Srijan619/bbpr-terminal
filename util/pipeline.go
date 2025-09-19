package util

import (
	"fmt"
	"simple-git-terminal/types"

	"github.com/gdamore/tcell/v2"
)

// Color mapping
func GetColorForStatus(status types.PipelineStatus) tcell.Color {
	switch {
	case status.Failed():
		return tcell.ColorRed
	case status.Passed():
		return tcell.ColorGreen
	case status.Running():
		return tcell.ColorYellow
	case status.Successful():
		return tcell.ColorGreen
	case status.Pending():
		return tcell.ColorBlue
	case status.Error():
		return tcell.ColorDarkRed
	case status.InProgress():
		return tcell.ColorOrange
	default:
		return tcell.ColorGray
	}
}

func GetIconForStatus(status types.PipelineStatus) string {
	switch {
	case status.Passed(), status.Successful():
		return "\u2714" // ✔ Check mark
	case status.Failed():
		return "\u2716" // ✖ Cross mark
	case status.Pending():
		return "\u23F3" // ⏳ Hourglass
	case status.Running():
		return "\u25B6" // ▶ Play button
	case status.Stopped():
		return "\u25A0" // ■ Stop square
	case status.Error():
		return "\u26A0" // ⚠ Warning sign
	case status.InProgress():
		return "\u23F3" // ⏳ Hourglass (same as pending, but you can pick another)
	default:
		return "\u2753" // ❓ Question mark
	}
}

func tcellColorToDynamicColor(c tcell.Color) string {
	switch c {
	case tcell.ColorGreen:
		return "green"
	case tcell.ColorRed:
		return "red"
	case tcell.ColorYellow:
		return "yellow"
	case tcell.ColorGray:
		return "gray"
	case tcell.ColorDarkRed:
		return "maroon"
	case tcell.ColorOrange:
		return "orange"
	case tcell.ColorBlue:
		return "blue"
	default:
		return "white"
	}
}

func GetIconForStatusWithColor(status types.PipelineStatus) string {
	icon := GetIconForStatus(status)
	color := GetColorForStatus(status)

	return fmt.Sprintf("[%s]%s[-]", tcellColorToDynamicColor(color), icon)
}
