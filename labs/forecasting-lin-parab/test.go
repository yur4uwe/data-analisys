package forecastinglinparab

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
)

func RenderTest(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if !hasTrained {
		return res.NewError("Run Training on 'Training Data' chart first to determine optimal values.")
	}

	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	copyChart := charting.CopyChart(TestDataChart)
	copyChart.Labels = testData.Date
	copyChart.UpdateDataPointsForDataset(GraphOriginalDataID, charting.F64ToPoints(testData.ExchangeRate))

	lincoefs := []float64{bestLinA, bestLinB}
	parcoefs := []float64{bestParA, bestParB, bestParC}

	predicted := make([]float64, len(testData.Date))
	for i := range testData.Date {
		predicted[i] = analysis.EvaluatePolynomial(lincoefs, float64(i))
	}

	error := analysis.MSE(testData.ExchangeRate, predicted)
	copyChart.ChartVariables[0].Label = fmt.Sprintf("Linear Fit (MSE: %.4e, a=%.4f, b=%.4f) for y=bx+a", error, lincoefs[0], lincoefs[1])

	copyChart.UpdateDataPointsForDataset(GraphLinearApproxID, charting.F64ToPoints(predicted))

	predicted = make([]float64, len(testData.Date))
	for i := range testData.Date {
		predicted[i] = analysis.EvaluatePolynomial(parcoefs, float64(i))
	}

	error = analysis.MSE(testData.ExchangeRate, predicted)
	copyChart.ChartVariables[1].Label = fmt.Sprintf("Parabolic Fit (MSE: %.4e, a=%.4f, b=%.4f, c=%.4f) for y=cx^2+bx+a", error, parcoefs[0], parcoefs[1], parcoefs[2])

	copyChart.UpdateDataPointsForDataset(GraphParabolicApproxID, charting.F64ToPoints(predicted))

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	res.CachePolicy = charting.CachePolicyDontCache
	return res
}
