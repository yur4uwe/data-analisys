package visualization

import "labs/labs/common"

const (
	BarChartID = "bar"

	BarGraphID = "bar-graph"
)

var (
	BarGraph = common.ChartDataset{
		Label:           "Bar Representation",
		Type:            common.ChartTypeBar,
		BorderColor:     common.Color1,
		BackgroundColor: "rgba(0, 0, 0, 0.1)",
		BorderWidth:     10,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
	}

	BarChart = common.Chart{
		ID:          BarChartID,
		Title:       "Bar Plot",
		XAxisLabel:  "Spending Type",
		YAxisLabel:  "Amount Spent",
		XAxisConfig: common.CategoryAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			BarGraphID: &BarGraph,
		},
	}

	BarMeta = BarChart.Meta()
)

func RenderBarPlot(req *common.RenderRequest) (res *common.RenderResponse) {
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

	chartCopy := common.CopyChart(BarChart)

	chartCopy.UpdatePointsForDataset(BarChartID, x, y)

	chartCopy.Labels = labels

	res.AddChart(BarChartID, &chartCopy)

	return
}
