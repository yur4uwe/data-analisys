package visualization

import (
	"fmt"
	"labs/labs/common"
)

const (
	BarChartID = "bar"

	BarGraphID = "orig-data"
)

var (
	BarGraph = common.ChartDataset{
		Label: "Bar Representation",
		BackgroundColor: []string{
			common.Color1,
			common.Color2,
			common.Color3,
			common.Color4,
			common.Color5,
		},
		BorderColor: "rgba(0, 0, 0, 0.1)",
		BorderWidth: 2,
		PointRadius: 0,
		ShowLine:    true,
		Togglable:   true,
	}

	BarChart = common.Chart{
		ID:          BarChartID,
		Title:       "Bar Plot",
		Type:        common.ChartTypeBar,
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
	fmt.Printf("Rendering %s\n", req.ChartID)
	values, err := ReadCategoricalCSV("../data/lab_4_var_12.csv")
	if err != nil {
		return res.NewErrorf("encountered error while reading csv: %v", err)
	}

	y := make([]float64, 0, len(values))
	labels := make([]string, 0, len(values))
	for k, v := range values {
		y = append(y, v)
		labels = append(labels, k)
	}

	chartCopy := common.CopyChart(BarChart)

	err = chartCopy.UpdateDataForDataset(BarGraphID, y)
	if err != nil {
		return res.NewErrorf("encountered error while updating points: %v", err)
	}

	chartCopy.Labels = labels

	res = common.NewRenderResponse()
	res.AddChart(BarChartID, &chartCopy)

	return res
}
