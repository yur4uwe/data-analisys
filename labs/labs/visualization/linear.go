package visualization

import (
	"fmt"
	"labs/labs/common"
	"labs/labs/render"
)

const (
	LinearChartID = "linear"

	LinearGraphID = "orig-data" // due to the default, baseline graph on frontend, one graph from chart should always be called orig-data
)

var (
	LinearGraph = common.ChartDataset{
		Label:           "Linear Data",
		BorderColor:     common.Color1,
		BackgroundColor: []string{"rgba(37, 99, 235, 0.1)"},
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	LinearChart = common.Chart{
		ID:          LinearChartID,
		Title:       "Linear Representation",
		Type:        common.ChartTypeLine,
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
	fmt.Printf("Rendering %s\n", req.ChartID)
	// x, y, err := polyapprox.ReadSampleCSV("../data/lab_3_var_12.csv")
	// if err != nil {
	// 	fmt.Println("failed to open file:", err)
	// 	return &common.RenderResponse{
	// 		Error: fmt.Errorf("encountered error while reading csv: %v", err),
	// 	}
	// }

	// chartCopy := common.CopyChart(LinearChart)
	// chartCopy.UpdatePointsForDataset(LinearGraphID, x, y)

	return &common.RenderResponse{
		Error: render.NewRenderError("impossible to visualize data"),
	}
}
