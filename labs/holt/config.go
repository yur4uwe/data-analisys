package holt

import (
	"labs/charting"
)

const (
	LabID = "8"

	ChartHoltTrainID   = "holt-train"
	ChartHoltTestID    = "holt-test"
	ChartHoltOptimalID = "holt-optimal"

	GraphTrainActualID   = "train-actual"
	GraphTrainForecastID = "train-forecast"
	GraphErrHeatmapID    = "error-heatmap"

	GraphTestActualID   = "test-actual"
	GraphTestForecastID = "test-forecast"

	VariableEpochsID       = "epochs"
	VariableLearningRateID = "learning-rate"
	VariableParamStepID    = "param-step"

	DisplayOptimalAlphaID = "optimal-alpha"
	DisplayOptimalBetaID  = "optimal-beta"
	DisplayTrainMSEID     = "train-mse"
	DisplayTestMSEID      = "test-mse"
	DisplayOptimalMSEID   = "optimal-mse"
)

type ExchangeRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

var (
	VariableEpochs = charting.MutableField{
		ID:      VariableEpochsID,
		Label:   "Gradient Descent Epochs",
		Default: 1000,
		Min:     10,
		Max:     10000,
		Step:    10,
		Control: charting.ControlRange,
	}

	VariableLearningRate = charting.MutableField{
		ID:      VariableLearningRateID,
		Label:   "Learning Rate",
		Default: 1.0,
		Min:     0.01,
		Max:     100.0,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	VariableHeatmapParamStep = charting.MutableField{
		ID:      VariableParamStepID,
		Label:   "Parameter Step Size",
		Default: 0.05,
		Min:     0.001,
		Max:     0.1,
		Step:    0.001,
		Control: charting.ControlRange,
	}

	OptimalAlphaField = charting.MutableField{ID: DisplayOptimalAlphaID, Label: "Optimal Alpha: -", Control: charting.ControlNoControl}
	OptimalBetaField  = charting.MutableField{ID: DisplayOptimalBetaID, Label: "Optimal Beta: -", Control: charting.ControlNoControl}
	TrainMSEField     = charting.MutableField{ID: DisplayTrainMSEID, Label: "Train MSE: -", Control: charting.ControlNoControl}
	TestMSEField      = charting.MutableField{ID: DisplayTestMSEID, Label: "Test MSE: -", Control: charting.ControlNoControl}
	OptimalMSEField   = charting.MutableField{ID: DisplayOptimalMSEID, Label: "Optimal MSE: -", Control: charting.ControlNoControl}

	TrainActualGraph = charting.ChartDataset{
		Label:           "Train Data",
		BorderColor:     charting.ColorTeal,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       false,
	}

	TrainForecastGraph = charting.ChartDataset{
		Label:           "Holt Forecast (Train)",
		BorderColor:     charting.ColorAmber,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []charting.MutableField{
			OptimalAlphaField,
			OptimalBetaField,
			TrainMSEField,
		},
	}

	TestActualGraph = charting.ChartDataset{
		Label:           "Test Data",
		BorderColor:     charting.ColorTeal,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       false,
	}

	TestForecastGraph = charting.ChartDataset{
		Label:           "Holt Forecast (Test)",
		BorderColor:     charting.ColorRed,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []charting.MutableField{
			TestMSEField,
		},
	}

	HeatmapGraph = charting.ChartDataset{
		Label:           "Holt error vs alpha and beta",
		BorderColor:     charting.ColorTransparent,
		BackgroundColor: []string{charting.ColorBlue, charting.ColorRed},
		BorderWidth:     0,
		PointRadius:     0,
		GraphVariables: []charting.MutableField{
			OptimalMSEField,
			OptimalAlphaField,
			OptimalBetaField,
		},
	}

	TrainChart = charting.Chart{
		ID:          ChartHoltTrainID,
		Title:       "Holt's Method - Training Phase",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			GraphTrainActualID:   &TrainActualGraph,
			GraphTrainForecastID: &TrainForecastGraph,
		},
		ChartVariables: []charting.MutableField{
			VariableEpochs,
			VariableLearningRate,
		},
	}

	TestChart = charting.Chart{
		ID:          ChartHoltTestID,
		Title:       "Holt's Method - Testing Phase",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			GraphTestActualID:   &TestActualGraph,
			GraphTestForecastID: &TestForecastGraph,
		},
	}

	OptimalChart = charting.Chart{
		ID:          ChartHoltOptimalID,
		Title:       "Heatmap of errors vs alpha and beta",
		Type:        charting.ChartTypeHeatmap,
		XAxisLabel:  "Alpha",
		YAxisLabel:  "Beta",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableHeatmapParamStep,
		},
		Datasets: map[string]*charting.ChartDataset{
			GraphErrHeatmapID: &HeatmapGraph,
		},
	}

	Config = charting.NewLabConfig(
		LabID,
		"Holt's Linear Trend Forecasting",
		map[string]*charting.Chart{
			ChartHoltTestID:    &TestChart,
			ChartHoltTrainID:   &TrainChart,
			ChartHoltOptimalID: &OptimalChart,
		},
	)

	Metadata = Config.Lab

	testExchangeRateData  = &ExchangeRateHistory{}
	trainExchangeRateData = &ExchangeRateHistory{}
)
