package stats

import (
	"labs/charting"
	"labs/uncsv"
	"os"
	"sort"
)

const (
	EmpiricalDistributionChartID = "empirical-distribution"
	EmpiricalDistributionGraphID = "empirical-distribution"
)

var (
	EmpiricalDistributionGraph = charting.ChartDataset{
		Label:           "Empirical Distribution Function F(x)",
		BorderColor:     charting.Color2,
		BackgroundColor: []string{charting.ColorTransparent},
		ShowLine:        true,
		PointRadius:     3,
		BorderWidth:     2,
		Togglable:       false,
	}

	EmpiricalDistributionChart = charting.Chart{
		ID:          EmpiricalDistributionChartID,
		Title:       "Empirical Distribution Function of Salaries",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Salary (USD)",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "F(x) - Cumulative Probability",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			EmpiricalDistributionGraphID: &EmpiricalDistributionGraph,
		},
	}
)

func RenderEmpiricalDistribution(req *charting.RenderRequest) (res *charting.RenderResponse) {
	// Load data if not already loaded
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("empirical distribution chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ';'
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("empirical distribution chart: error while decoding csv: %s", err.Error())
		}
	}

	salaries := make([]float64, len(salaryRecords.Salary))
	copy(salaries, salaryRecords.Salary)

	// Sort salaries for easier cumulative probability calculation
	sort.Float64s(salaries)

	n := float64(len(salaries))

	x := make([]float64, 0)
	y := make([]float64, 0)

	// Add 0 value for cumulative probability before first salary
	if len(salaries) > 0 {
		x = append(x, salaries[0]-1)
		y = append(y, 0)
	}

	for i, salary := range salaries {
		fx := float64(i+1) / n
		x = append(x, salary)
		y = append(y, fx)
	}

	copyChart := charting.CopyChart(EmpiricalDistributionChart)
	if err := copyChart.UpdatePointsForDataset(EmpiricalDistributionGraphID, x, y); err != nil {
		return res.NewErrorf("error updating dataset: %s", err.Error())
	}

	res = charting.NewRenderResponse()
	res.AddChart(EmpiricalDistributionChartID, &copyChart)
	return res
}
