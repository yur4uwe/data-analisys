package cluster

import "labs/charting"

const (
	LabID = "6"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Data Clustering",
		map[string]*charting.Chart{
			ScatterPointsChartID: &ScatterPointsChart,
		},
	)

	Metadata = Config.Lab

	points = (*Points)(nil)
)

type Points struct {
	X []float64 `csv:"x"`
	Y []float64 `csv:"y"`
}

func init() {
	ScatterPointsChart.RenderFunc = RenderScatteredPoints
}
