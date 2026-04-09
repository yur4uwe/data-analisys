package forecastinglinparab

import (
	"labs/analysis"
	"labs/charting"
)

var (
	MaxPolyDegreeField = charting.MutableField{
		ID:      VariableMaxPolyDegreeID,
		Label:   "Max Polynomial Degree",
		Default: 10,
		Min:     1,
		Max:     20,
		Step:    1,
		Control: charting.ControlRange,
	}

	MSEByDegreeChart = charting.Chart{
		ID:          ChartMSEByDegreeID,
		Type:        charting.ChartTypeMultiLine,
		Title:       "MSE vs Polynomial Degree",
		XAxisLabel:  "Degree",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "MSE",
		YAxisConfig: charting.LogarithmicAxis,
		Datasets: map[string]charting.Dataset{
			GraphTrainMSEID: &TrainMSEGraph,
			GraphTestMSEID:  &TestMSEGraph,
		},
		ChartVariables: []charting.MutableField{
			MaxPolyDegreeField,
		},
	}

	TrainMSEGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Training MSE",
			BorderColor: charting.ColorBlue,
			BorderWidth: 2,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     4,
	}

	TestMSEGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Validation (Test) MSE",
			BorderColor: charting.ColorRed,
			BorderWidth: 2,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     4,
	}
)

func init() {
	MSEByDegreeChart.RenderFunc = RenderMSEByDegree
}

func RenderMSEByDegree(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	maxDegreeVal, ok := req.GetChartVariable(ChartMSEByDegreeID, VariableMaxPolyDegreeID)
	if !ok {
		maxDegreeVal = MaxPolyDegreeField.Default
	}
	maxDegree := int(maxDegreeVal)

	splitIdx := len(trainData.ExchangeRate)
	trainX := make([]float64, splitIdx)
	for i := range splitIdx {
		trainX[i] = float64(i)
	}

	degrees := make([]float64, 0, maxDegree)
	trainMSEs := make([]float64, 0, maxDegree)
	testMSEs := make([]float64, 0, maxDegree)

	for d := 1; d <= maxDegree; d++ {
		coefs, err := analysis.SolvePolynomialFit(trainX, trainData.ExchangeRate, d)
		if err != nil {
			continue
		}

		trainPredicted := make([]float64, splitIdx)
		for i := range splitIdx {
			trainPredicted[i] = analysis.EvaluatePolynomial(coefs, float64(i))
		}
		trainMSE := analysis.MSE(trainData.ExchangeRate, trainPredicted)

		testPredicted := make([]float64, len(testData.ExchangeRate))
		for i := range len(testData.ExchangeRate) {
			testPredicted[i] = analysis.EvaluatePolynomial(coefs, float64(splitIdx+i))
		}
		testMSE := analysis.MSE(testData.ExchangeRate, testPredicted)

		degrees = append(degrees, float64(d))
		trainMSEs = append(trainMSEs, trainMSE)
		testMSEs = append(testMSEs, testMSE)
	}

	copyChart := charting.CopyChart(MSEByDegreeChart)
	copyChart.UpdatePointsForDataset(GraphTrainMSEID, degrees, trainMSEs)
	copyChart.UpdatePointsForDataset(GraphTestMSEID, degrees, testMSEs)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	res.CachePolicy = charting.CachePolicyDontCache
	return res
}
