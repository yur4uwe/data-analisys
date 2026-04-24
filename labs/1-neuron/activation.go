package neuron

import (
	"labs/charting"
	"math"
)

var (
	ActivationChart = charting.Chart{
		ID:          ActivationChartID,
		Title:       "Activation Function",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Input (z)",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "Output (y)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphLandscapeID: &ActivationLandscapeDataset,
		},
		ChartVariables: append(SharedVariables, ActivationFuncField),
	}

	ActivationLandscapeDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "f(z)",
			BorderColor: charting.ColorBlue,
			BorderWidth: 3,
		},
		HideLine:    false,
		PointRadius: 0,
	}
)

func RenderActivation(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()
	chartCopy := charting.CopyChart(ActivationChart)

	alpha := req.GetVariable(VarAlphaID)
	actIdx := int(req.GetVariable(VarActivationID))
	act, _ := getActivation(actIdx, alpha)

	points := make([]charting.DataPoint, 0, 201)
	for i := -100; i <= 100; i++ {
		z := float64(i) / 10.0
		y := act(z)

		// JSON Safety
		if math.IsNaN(y) || math.IsInf(y, 0) {
			continue
		}

		points = append(points, charting.DataPoint{
			X: z,
			Y: &y,
		})
	}

	chartCopy.UpdateDataPointsForDataset(GraphLandscapeID, points)
	res.AddChart(ActivationChartID, &chartCopy)
	return res
}

func newSigmoid(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		return 1 / (1 + math.Exp(-alpha*x))
	}
	df := func(y float64) float64 {
		return alpha * y * (1 - y)
	}
	return f, df
}

func newTanh(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		numerator := math.Exp(alpha*x) - math.Exp(-alpha*x)
		denominator := math.Exp(alpha*x) + math.Exp(-alpha*x)
		return numerator / denominator
	}
	df := func(y float64) float64 {
		return alpha * (1 - y*y)
	}
	return f, df
}

func newReLU(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		if x > alpha {
			return x
		}
		return 0
	}
	df := func(y float64) float64 {
		if y > alpha {
			return 1
		}
		return 0
	}
	return f, df
}
