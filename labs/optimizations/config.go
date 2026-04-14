package optimizations

import "labs/charting"

const (
	LabID = "10"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Optimizations",
		map[string]*charting.Chart{
			OneDimChartID: &OneDimChart,
			TwoDimChartID: &TwoDimChart,
		},
	)

	Metadata = Config.Lab
)

func init() {
	OneDimChart.RenderFunc = RenderOneDim
	TwoDimChart.RenderFunc = RenderTwoDim
}
