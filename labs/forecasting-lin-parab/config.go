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

	ChartMSEByDegreeID     = "mse-by-degree"
	ChartModelValidationID = "model-validation"

	GraphOriginalDataID    = "original-data"
	GraphLinearApproxID    = "linear-approx"
	GraphParabolicApproxID = "parabolic-approx"
	GraphTrainFitID        = "train-fit"
	GraphTestForecastID    = "test-forecast"
	GraphDividerID         = "divider"

	GraphTrainMSEID = "train-mse"
	GraphTestMSEID  = "test-mse"

	VariablePolyDegreeID       = "poly-degree"
	VariableMaxPolyDegreeID    = "max-poly-degree"
	VariableValidationDegreeID = "validation-degree"
)

var (
	exchangeRateData = &ExchangeRateHistory{}
	testData         = &ExchangeRateHistory{}
	trainData        = &ExchangeRateHistory{}

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

	LinParabConfig = charting.NewLabConfig(
		LabID,
		"Linear and Parabolic Approximation",
		map[string]*charting.Chart{
			ChartMSEByDegreeID:     &MSEByDegreeChart,
			ChartModelValidationID: &ModelValidationChart,
		},
	)
)

type ExchangeRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

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
