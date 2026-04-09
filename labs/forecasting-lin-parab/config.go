package forecastinglinparab

import (
	"errors"
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	LabID = "9"

	ChartTrainDataID              = "train-data"
	ChartTestDataID               = "test-data"
	ChartOptimalParabolicParamsID = "optimal-parabolic-params"
	ChartOptimalLinearParamsID    = "optimal-linear-params"

	GraphOriginalDataID    = "original-data"
	GraphLinearApproxID    = "linear-approx"
	GraphParabolicApproxID = "parabolic-approx"

	VariableParabolicFitCoefficientsID = "parabolic-fit-coefficients"
	VariableLinearFitCoefficientsID    = "linear-fit-coefficients"
)

var (
	bestLinA = 0.0
	bestLinB = 0.0

	bestParA = 0.0
	bestParB = 0.0
	bestParC = 0.0

	hasTrained = false

	exchangeRateData = &ExchangeRateHistory{}
	testData         = &ExchangeRateHistory{}
	trainData        = &ExchangeRateHistory{}
)

type ExchangeRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

var (
	LinParabConfig = charting.NewLabConfig(
		LabID,
		"Linear and Parabolic Approximation",
		map[string]*charting.Chart{
			ChartTrainDataID: &TrainDataChart,
			ChartTestDataID:  &TestDataChart,
		},
	)

	TrainDataChart = charting.Chart{
		ID:          ChartTrainDataID,
		Type:        charting.ChartTypeLine,
		Title:       "Training Data",
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Amount",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphOriginalDataID:    &OriginalDataGraph,
			GraphLinearApproxID:    &LinearApproxGraph,
			GraphParabolicApproxID: &ParabolicApproxGraph,
		},
		ChartVariables: []charting.MutableField{
			LinearFitCoefficients,
			ParabolicFitCoefficients,
		},
	}

	TestDataChart = charting.Chart{
		ID:          ChartTestDataID,
		Type:        charting.ChartTypeLine,
		Title:       "Testing Data",
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Amount",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphOriginalDataID:    &OriginalDataGraph,
			GraphLinearApproxID:    &LinearApproxGraph,
			GraphParabolicApproxID: &ParabolicApproxGraph,
		},
		ChartVariables: []charting.MutableField{
			LinearFitCoefficients,
			ParabolicFitCoefficients,
		},
	}

	OriginalDataGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Original Data",
			BorderColor: charting.ColorEmerald,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     2,
	}

	LinearFitCoefficients = charting.MutableField{
		ID:      VariableLinearFitCoefficientsID,
		Label:   "Linear Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	ParabolicFitCoefficients = charting.MutableField{
		ID:      VariableParabolicFitCoefficientsID,
		Label:   "Parabolic Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	LinearApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Linear Approximation",
			BorderColor: charting.ToColor("#16a34a"),
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ToColor("rgba(22, 163, 74, 0.1)"),
		PointRadius:     0,
	}

	ParabolicApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Parabolic Approximation",
			BorderColor: charting.ToColor("#9333ea"),
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ToColor("rgba(147, 51, 234, 0.1)"),
		PointRadius:     0,
	}
)

func loadExchangeHistory() error {
	if len(trainData.ExchangeRate) > 0 && len(testData.ExchangeRate) > 0 {
		return nil
	}
	f, err := os.Open("./data/lab_9_var_12.csv")
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.Comma = ','

	exchangeRateData = &ExchangeRateHistory{}
	if err := d.Decode(exchangeRateData); err != nil {
		return err
	}
	n := len(exchangeRateData.ExchangeRate)
	if n < 4 {
		return errors.New("not enough data for training and testing")
	}

	splitIdx := n / 2
	trainData.ExchangeRate = exchangeRateData.ExchangeRate[:splitIdx]
	trainData.Date = exchangeRateData.Date[:splitIdx]
	testData.ExchangeRate = exchangeRateData.ExchangeRate[splitIdx:]
	testData.Date = exchangeRateData.Date[splitIdx:]

	return nil
}

func init() {
	TrainDataChart.RenderFunc = RenderTrain
	TestDataChart.RenderFunc = RenderTest
}
