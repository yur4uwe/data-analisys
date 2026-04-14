package optimizations

import (
	"labs/charting"
	"math"
)

const (
	OneDimChartID = "one-dim"
	TwoDimChartID = "two-dim"

	OneDimGraphID = "orig-data"
	TwoDimGraphID = "orig-data"

	VarStartID = "start"
	VarEndID   = "end"
	VarStepID  = "step"

	VarXStartID = "x-start"
	VarXEndID   = "x-end"
	VarXStepID  = "x-step"
	VarYStartID = "y-start"
	VarYEndID   = "y-end"
	VarYStepID  = "y-step"
)

var (
	VarStart = charting.MutableField{
		ID:      VarStartID,
		Label:   "Start (x > 0, x != 1)",
		Default: 1.1,
		Min:     0.1,
		Max:     100.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarEnd = charting.MutableField{
		ID:      VarEndID,
		Label:   "End",
		Default: 10.0,
		Min:     0.2,
		Max:     200.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarStep = charting.MutableField{
		ID:      VarStepID,
		Label:   "Step",
		Default: 0.1,
		Min:     0.01,
		Max:     5.0,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	VarXStart = charting.MutableField{
		ID:      VarXStartID,
		Label:   "X Start",
		Default: -5.0,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarXEnd = charting.MutableField{
		ID:      VarXEndID,
		Label:   "X End",
		Default: 5.0,
		Min:     0.2,
		Max:     20.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarXStep = charting.MutableField{
		ID:      VarXStepID,
		Label:   "X Step",
		Default: 0.2,
		Min:     0.05,
		Max:     2.0,
		Step:    0.05,
		Control: charting.ControlRange,
	}

	VarYStart = charting.MutableField{
		ID:      VarYStartID,
		Label:   "Y Start",
		Default: -5.0,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarYEnd = charting.MutableField{
		ID:      VarYEndID,
		Label:   "Y End",
		Default: 5.0,
		Min:     -10.0,
		Max:     20.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarYStep = charting.MutableField{
		ID:      VarYStepID,
		Label:   "Y Step",
		Default: 0.2,
		Min:     0.05,
		Max:     2.0,
		Step:    0.05,
		Control: charting.ControlRange,
	}

	OneDimChart = charting.Chart{
		ID:          OneDimChartID,
		Title:       "One-Dimensional Function",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "x",
		YAxisLabel:  "f(x)",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarStart,
			VarEnd,
			VarStep,
		},
		Datasets: map[string]charting.Dataset{
			OneDimGraphID: &OneDimGraph,
		},
	}

	OneDimGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "f(x) = x + 1/ln(x)",
			BorderColor: charting.ColorBlue,
			BorderWidth: 2,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	TwoDimChart = charting.Chart{
		ID:          TwoDimChartID,
		Title:       "Two-Dimensional Function",
		Type:        charting.ChartTypeSurface,
		XAxisLabel:  "x",
		YAxisLabel:  "y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarXStart, VarXEnd, VarXStep,
			VarYStart, VarYEnd, VarYStep,
		},
		Datasets: map[string]charting.Dataset{
			TwoDimGraphID: &TwoDimGraph,
		},
	}

	TwoDimGraph = charting.HeatmapDataset{
		BaseDataset: charting.BaseDataset{
			Label: "f(x, y)",
		},
	}
)

func onedimf(x float64) float64 {
	return x + (1 / math.Log(x))
}

// z=-\ln\left(\left(\left(x-c_{1}\right)^{2}-\left(y-c_{2}\right)^{2}\right)^{2}+c_{1}^{2}-c_{2}^{2}\right)
func twodimf(x, y float64) float64 {
	c1 := 1.9
	c2 := 0.0

	term1 := (x - c1) * (x - c1)
	term2 := (y - c2) * (y - c2)

	inner := term1 - term2
	res := (inner * inner) + (c1 * c1) - (c2 * c2)

	if res <= 0 {
		return 0
	}

	return -math.Log(res)
}

func RenderOneDim(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()

	start, _ := req.GetChartVariable(OneDimChartID, VarStartID)
	end, _ := req.GetChartVariable(OneDimChartID, VarEndID)
	step, _ := req.GetChartVariable(OneDimChartID, VarStepID)

	if start <= 0 || (start <= 1 && end >= 1) {
		return res.NewError("invalid start value: x must be > 0 and != 1")
	}

	n := int((end-start)/step) + 1
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)

	for i := range n {
		xVal := start + float64(i)*step
		x = append(x, xVal)
		y = append(y, onedimf(xVal))
	}

	chartCopy := charting.CopyChart(OneDimChart)
	chartCopy.UpdatePointsForDataset(OneDimGraphID, x, y)
	chartCopy.GenerateLabels(2)
	res.AddChart(OneDimChartID, &chartCopy)

	return res
}

func RenderTwoDim(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()

	xStart, _ := req.GetChartVariable(TwoDimChartID, VarXStartID)
	xEnd, _ := req.GetChartVariable(TwoDimChartID, VarXEndID)
	xStep, _ := req.GetChartVariable(TwoDimChartID, VarXStepID)

	yStart, _ := req.GetChartVariable(TwoDimChartID, VarYStartID)
	yEnd, _ := req.GetChartVariable(TwoDimChartID, VarYEndID)
	yStep, _ := req.GetChartVariable(TwoDimChartID, VarYStepID)

	nx := int((xEnd-xStart)/xStep) + 1
	ny := int((yEnd-yStart)/yStep) + 1

	points := make([]any, 0, nx*ny)

	for i := range nx {
		for j := range ny {
			xVal := xStart + float64(i)*xStep
			yVal := yStart + float64(j)*yStep
			zVal := twodimf(xVal, yVal)

			yCopy := yVal
			zCopy := zVal

			points = append(points, charting.HeatmapPoint{
				DataPoint: charting.DataPoint{X: xVal, Y: &yCopy},
				Value:     &zCopy,
			})
		}
	}

	chartCopy := charting.CopyChart(TwoDimChart)
	chartCopy.UpdateDataForDataset(TwoDimGraphID, points)
	chartCopy.GenerateLabels(2)
	res.AddChart(TwoDimChartID, &chartCopy)

	return res
}
