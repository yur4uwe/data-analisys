package cluster

import (
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	ScatterPointsGraphID = "scatter"
	ScatterPointsChartID = "scatter"
)

var (
	ScatterPointsGraph = charting.ChartDataset{
		Label:           "Points awaiting clusterization",
		BorderColor:     charting.ColorAmber,
		BackgroundColor: []string{charting.ColorTransparent},
		PointRadius:     3,
		BorderWidth:     2,
	}

	ScatterPointsChart = charting.Chart{
		ID:          ScatterPointsChartID,
		Title:       "Scattered points",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "X",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "Y",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			ScatterPointsGraphID: &ScatterPointsGraph,
		},
	}
)

func RenderScatteredPoints(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if points == nil {
		f, err := os.Open("./data/lab_6_var_12.csv")
		if err != nil {
			return res.NewErrorf("clustering points chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ','
		points = &Points{}
		if err := d.Decode(points); err != nil {
			return res.NewErrorf("clustering points chart: error while decoding csv: %s", err.Error())
		}
	}

	copyChart := charting.CopyChart(ScatterPointsChart)

	if err := copyChart.UpdatePointsForDataset(ScatterPointsGraphID, points.X, points.Y); err != nil {
		return res.NewErrorf("clustering points chart: error while updating points: %s", err.Error())
	}

	res = charting.NewRenderResponse()
	res.AddChart(ScatterPointsChartID, &copyChart)

	return res
}
