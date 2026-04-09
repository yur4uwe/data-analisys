package my_lab

import (
	"labs/charting"
)

const LabID = "new_lab_id"

var Template = charting.Chart{
	ID:         "chart_1",
	Title:      "New Visualization",
	Type:       charting.ChartTypeLine,
	XAxisLabel: "X Axis",
	YAxisLabel: "Y Axis",
	Datasets: map[string]charting.Dataset{
		"ds_1": &charting.GridDataset{
			BaseDataset: charting.BaseDataset{
				Label:       "Dataset 1",
				BorderColor: charting.ToColor("#36A2EB"),
				BorderWidth: 2,
			},
		},
	},
}

var Config = charting.NewLabConfig(
	LabID,
	"Lab Title",
	map[string]*charting.Chart{
		Template.ID: &Template,
	},
)

func init() {
	Template.RenderFunc = Render
}

func Render(req *charting.RenderRequest) *charting.RenderResponse {
	res := charting.NewRenderResponse()
	copyChart := charting.CopyChart(Template)

	// Logic for data population goes here
	// example: val, _ := req.GetChartVariable(Template.ID, "var_id")

	res.AddChart(copyChart.ID, &copyChart)
	return res
}
