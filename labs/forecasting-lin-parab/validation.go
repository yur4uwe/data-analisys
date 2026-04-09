package forecastinglinparab

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
)

var (
	ValidationDegreeField = charting.MutableField{
		ID:      VariableValidationDegreeID,
		Label:   "Degree to Validate",
		Default: 3,
		Min:     1,
		Max:     15,
		Step:    1,
		Control: charting.ControlRange,
	}

	ModelValidationChart = charting.Chart{
		ID:          ChartModelValidationID,
		Type:        charting.ChartTypeLine,
		Title:       "Model Validation (Train vs Test)",
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphOriginalDataID:    &OriginalDataGraph,
			GraphTrainFitID:        &TrainFitGraph,
			GraphTestForecastID:    &TestForecastGraph,
			GraphLinearApproxID:    &LinearApproxGraph,
			GraphParabolicApproxID: &ParabolicApproxGraph,
			GraphDividerID:         &DividerGraph,
		},
		ChartVariables: []charting.MutableField{
			ValidationDegreeField,
		},
	}

	TrainFitGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Model Fit (Training)",
			BorderColor: charting.ColorBlue,
			BorderWidth: 3,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	TestForecastGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Extrapolated (Testing)",
			BorderColor: charting.ColorOrange,
			BorderWidth: 3,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	DividerGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Training/Testing Split",
			BorderColor: charting.ToColor("#000000"),
			BorderWidth: 2,
			Togglable:   false,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
		HideLine:        false,
	}

	LinearApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Linear Trend (Degree 1)",
			BorderColor: charting.ToColor("#16a34a"),
			BorderWidth: 1,
			Togglable:   true,
			Hidden:      true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	ParabolicApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Parabolic Trend (Degree 2)",
			BorderColor: charting.ToColor("#9333ea"),
			BorderWidth: 1,
			Togglable:   true,
			Hidden:      true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}
)

func init() {
	ModelValidationChart.RenderFunc = RenderModelValidation
}

func RenderModelValidation(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	degreeVal, ok := req.GetChartVariable(ChartModelValidationID, VariableValidationDegreeID)
	if !ok {
		degreeVal = ValidationDegreeField.Default
	}
	degree := int(degreeVal)

	fullRates := append(trainData.ExchangeRate, testData.ExchangeRate...)
	fullDates := append(trainData.Date, testData.Date...)
	splitIdx := len(trainData.ExchangeRate)

	trainX := make([]float64, splitIdx)
	for i := range splitIdx {
		trainX[i] = float64(i)
	}

	coefs, err := analysis.SolvePolynomialFit(trainX, trainData.ExchangeRate, degree)
	if err != nil {
		return res.NewErrorf("failed to solve polynomial fit (degree %d): %v", degree, err)
	}

	fitY := make([]*float64, len(fullRates))
	forecastY := make([]*float64, len(fullRates))
	for i := range fullRates {
		val := analysis.EvaluatePolynomial(coefs, float64(i))
		if i < splitIdx {
			fitY[i] = &val
			forecastY[i] = nil
		} else {
			fitY[i] = nil
			forecastY[i] = &val
		}
	}
	if splitIdx > 0 && splitIdx < len(fullRates) {
		val := analysis.EvaluatePolynomial(coefs, float64(splitIdx-1))
		forecastY[splitIdx-1] = &val
	}

	linCoefs, _ := analysis.SolvePolynomialFit(trainX, trainData.ExchangeRate, 1)
	linPredicted := make([]float64, len(fullRates))
	for i := range fullRates {
		linPredicted[i] = analysis.EvaluatePolynomial(linCoefs, float64(i))
	}

	parCoefs, _ := analysis.SolvePolynomialFit(trainX, trainData.ExchangeRate, 2)
	parPredicted := make([]float64, len(fullRates))
	for i := range fullRates {
		parPredicted[i] = analysis.EvaluatePolynomial(parCoefs, float64(i))
	}

	copyChart := charting.CopyChart(ModelValidationChart)
	copyChart.Labels = fullDates

	copyChart.UpdateDataPointsForDataset(GraphOriginalDataID, charting.F64ToPoints(fullRates))
	copyChart.UpdateDataPointsForDataset(GraphTrainFitID, charting.F64PtrToPoints(fitY))
	copyChart.UpdateDataPointsForDataset(GraphTestForecastID, charting.F64PtrToPoints(forecastY))
	copyChart.UpdateDataPointsForDataset(GraphLinearApproxID, charting.F64ToPoints(linPredicted))
	copyChart.UpdateDataPointsForDataset(GraphParabolicApproxID, charting.F64ToPoints(parPredicted))

	minY, maxY := analysis.MinMax(fullRates)
	padding := (maxY - minY) * 0.1
	divMin, divMax := minY-padding, maxY+padding
	divPoints := make([]charting.DataPoint, 0, len(fullRates)+1)
	for i := range fullRates {
		if i == splitIdx {
			divPoints = append(divPoints, charting.DataPoint{X: float64(i), Y: &divMin})
			divPoints = append(divPoints, charting.DataPoint{X: float64(i), Y: &divMax})
		} else {
			divPoints = append(divPoints, charting.DataPoint{X: float64(i), Y: nil})
		}
	}
	copyChart.UpdateDataPointsForDataset(GraphDividerID, divPoints)

	trainMSE := analysis.MSE(trainData.ExchangeRate, charting.ExtractF64(fitY[:splitIdx]))
	testMSE := analysis.MSE(testData.ExchangeRate, charting.ExtractF64(forecastY[splitIdx:]))
	copyChart.Title = fmt.Sprintf("Validation (Degree %d) | Train MSE: %.2e | Test MSE: %.2e", degree, trainMSE, testMSE)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}
