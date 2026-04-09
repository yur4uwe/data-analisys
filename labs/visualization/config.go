package visualization

import (
	"labs/charting"
)

const (
	LabID = "4"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Visualization",
		map[string]*charting.Chart{
			BarChartID:      &BarChart,
			FunctionChartID: &FunctionChart,
			LinearChartID:   &LinearChart,
			RadialChartID:   &RadialChart,
		},
	)

	Metadata = Config.Lab
)

func init() {
	BarChart.RenderFunc = RenderBarPlot
	FunctionChart.RenderFunc = RenderFunction
	LinearChart.RenderFunc = RenderLinear
	RadialChart.RenderFunc = RenderRadialPlot
}
