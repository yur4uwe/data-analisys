package charting

import (
	"fmt"
	"regexp"
	"strings"
)

// DataPoint represents a point with x and y coordinates for scatter/bubble charts
type DataPoint struct {
	X float64 `json:"x" csv:"x"`
	Y float64 `json:"y" csv:"y"`
}

// FieldControl specifies what UI element should render a field
// Values correspond directly to HTML input types
type FieldControl string

const (
	ControlRange     FieldControl = "range"     // Slider input
	ControlNumber    FieldControl = "number"    // Number input
	ControlCheckbox  FieldControl = "checkbox"  // Checkbox input
	ControlSelect    FieldControl = "select"    // Select dropdown
	ControlText      FieldControl = "text"      // Text input
	ControlNoControl FieldControl = "nocontrol" // Special case for displaying a label only
)

type GraphType string

const (
	ChartTypeLine    GraphType = "line"
	ChartTypeBar     GraphType = "bar"
	ChartTypeScatter GraphType = "scatter"
	ChartTypeBubble  GraphType = "bubble"
	ChartTypePie     GraphType = "pie"
	ChartTypeHeatmap GraphType = "heatmap"

	ChartTypeMultiLine    GraphType = "multi-line"
	ChartTypeMultiBar     GraphType = "multi-bar"
	ChartTypeMultiScatter GraphType = "multi-scatter"
	ChartTypeMultiBubble  GraphType = "multi-bubble"
	ChartTypeMultiPie     GraphType = "multi-pie"
	ChartTypeMultiHeatmap GraphType = "multi-heatmap"
)

type AxisConfig string

const (
	LinearAxis      AxisConfig = "linear"
	LogarithmicAxis AxisConfig = "logarithmic"
	TimeAxis        AxisConfig = "time"
	CategoryAxis    AxisConfig = "category"
)

type Color string

func isHex(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	l := c | 32
	return l >= 'a' && l <= 'f'
}

func ToColor(colorlike string) Color {
	if strings.HasPrefix(colorlike, "#") {
		if len(colorlike) != 7 {
			panic(fmt.Sprintf("cannot turn %s into hex color with length %d", colorlike, len(colorlike)))
		}
		for i := range 6 {
			if !isHex(colorlike[i+1]) {
				panic(fmt.Sprintf("cannot turn %q into a hex color, character at index %d isn't hex valid"))
			}
		}
		return Color(colorlike)
	} else if strings.HasPrefix(colorlike, "rgb") {
		colorRegex := regexp.MustCompile(`(?i)rgba?\(\s*\d{1,3}\s*(?:,\s*\d{1,3}\s*){2}(?:,\s*(?:\d*\.)?\d+\s*)?\)`)
		if !colorRegex.Match([]byte(colorlike)) {
			panic(fmt.Sprintf("failed to parse %q with an rgba regex"))
		}

		return Color(colorlike)

	} else {
		panic(fmt.Sprintf("impossible to turn %q into color", colorlike))
	}
}

const (
	ColorBlue        = "#1d4ed8"
	ColorRed         = "#b91c1c"
	ColorAmber       = "#d97706"
	ColorGreen       = "#16a34a"
	ColorViolet      = "#6d28d9"
	ColorPurple      = "#7c3aed"
	ColorFuchsia     = "#c026d3"
	ColorOrange      = "#ea580c"
	ColorLightPurple = "#9333ea"
	ColorCrimson     = "#be123c"
	ColorEmerald     = "#059669"
	ColorCyan        = "#0891b2"
	ColorPink        = "#db2777"
	ColorLime        = "#65a30d"
	ColorTeal        = "#0d9488"
	ColorIndigo      = "#4f46e5"
	ColorRose        = "#e11d48"
	ColorSky         = "#0284c7"
	ColorYellow      = "#ca8a04"
	ColorSlate       = "#475569"

	ColorTransparent = "rgba(0, 0, 0, 0.1)"
)
