package forecastinglinparab

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
)

func RenderTrain(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	copyChart := charting.CopyChart(TrainDataChart)
	copyChart.Labels = trainData.Date
	copyChart.UpdateDataPointsForDataset(GraphOriginalDataID, charting.F64ToPoints(trainData.ExchangeRate))

	x, y := make([]float64, 0, len(trainData.Date)), make([]float64, 0, len(trainData.ExchangeRate))
	for i := range trainData.Date {
		x = append(x, float64(i))
		y = append(y, trainData.ExchangeRate[i])
	}

	lincoefs, err := analysis.SolvePolynomialFit(x, y, 1)
	if err != nil {
		return res.NewErrorf("failed to solve polynomial fit: %v", err)
	}

	bestLinA = lincoefs[0]
	bestLinB = lincoefs[1]

	predicted := make([]float64, len(trainData.Date))

	for i := range trainData.Date {
		predicted[i] = analysis.EvaluatePolynomial(lincoefs, float64(i))
	}

	error := analysis.MSE(trainData.ExchangeRate, predicted)

	copyChart.ChartVariables[0].Label = fmt.Sprintf("Linear Fit (len %d) (MSE: %.4e, a=%.4f, b=%.4f) for y=bx+a", len(lincoefs), error, lincoefs[0], lincoefs[1])

	copyChart.UpdateDataPointsForDataset(GraphLinearApproxID, charting.F64ToPoints(predicted))

	parcoefs, err := analysis.SolvePolynomialFit(x, y, 2)
	if err != nil {
		return res.NewErrorf("failed to solve polynomial fit: %v", err)
	}

	bestParA = parcoefs[0]
	bestParB = parcoefs[1]
	bestParC = parcoefs[2]

	predicted = make([]float64, len(trainData.Date))

	for i := range trainData.Date {
		predicted[i] = analysis.EvaluatePolynomial(parcoefs, float64(i))
	}

	error = analysis.MSE(trainData.ExchangeRate, predicted)

	copyChart.ChartVariables[1].Label = fmt.Sprintf("Parabolic Fit (MSE: %.4e, a=%.4f, b=%.4f, c=%.4f) for y=cx^2+bx+a", error, parcoefs[0], parcoefs[1], parcoefs[2])

	copyChart.UpdateDataPointsForDataset(GraphParabolicApproxID, charting.F64ToPoints(predicted))

	hasTrained = true

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	res.CachePolicy = charting.CachePolicyDontCache
	return res
}
