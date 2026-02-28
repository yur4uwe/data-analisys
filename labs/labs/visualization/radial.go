package visualization

import (
	"labs/labs/common"
)

const (
	RadialChartID = "radial"

	RadialGraphID = "radial-graph"
)

var (
	RadialGraph = common.ChartDataset{
		Label:           "Radial Representation",
		Type:            common.ChartTypeBar,
		BorderColor:     common.Color2,
		BackgroundColor: "rgba(220, 38, 38, 0.1)",
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
	}

	RadialChart = common.Chart{
		ID:          RadialChartID,
		Title:       "Radial Plot",
		XAxisLabel:  "Category",
		YAxisLabel:  "Amount",
		XAxisConfig: common.CategoryAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			RadialGraphID: &RadialGraph,
		},
	}

	RadialMeta = RadialChart.Meta()
)

func RenderRadialPlot(req *common.RenderRequest) (res *common.RenderResponse) {
	values, err := ReadCategoricalCSV("../data/lab_4_var_12.csv")
	if err != nil {
		return res.NewErrorf("encountered error while reading csv: %v", err)
	}

	x := make([]float64, 0, len(values))
	y := make([]float64, 0, len(values))

	for i := range len(values) {
		x = append(x, float64(i))
	}

	labels := make([]string, 0, len(values))
	for k, v := range values {
		y = append(y, v)
		labels = append(labels, k)
	}

	chartCopy := common.CopyChart(RadialChart)

	chartCopy.UpdatePointsForDataset(RadialGraphID, x, y)

	chartCopy.Labels = labels

	res.AddChart(RadialChartID, &chartCopy)

	return
}
