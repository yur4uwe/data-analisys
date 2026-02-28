package visualization

import (
	"fmt"
	"labs/labs/common"
	"labs/labs/polyapprox"
)

const (
	LinearChartID = "linear"

	LinearGraphID = "linear-graph"
)

var (
	LinearGraph = common.ChartDataset{
		Label:           "Linear Data",
		Type:            common.ChartTypeLine,
		BorderColor:     common.Color1,
		BackgroundColor: "rgba(37, 99, 235, 0.1)",
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	LinearChart = common.Chart{
		ID:          LinearChartID,
		Title:       "Linear Representation",
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			LinearGraphID: &LinearGraph,
		},
	}

	LinearMeta = LinearChart.Meta()
)

func RenderLinear(req *common.RenderRequest) *common.RenderResponse {
	x, y, err := polyapprox.ReadSampleCSV("../data/lab_3_var_12.csv")
	if err != nil {
		fmt.Println("failed to open file:", err)
		return &common.RenderResponse{
			Error: fmt.Errorf("encountered error while reading csv: %v", err),
		}
	}

	chartCopy := common.CopyChart(LinearChart)
	chartCopy.UpdatePointsForDataset(LinearGraphID, x, y)

	return &common.RenderResponse{
		Charts: map[string]common.Chart{
			LinearChartID: chartCopy,
		},
	}
}
