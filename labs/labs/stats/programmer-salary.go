package stats

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
)

const (
	ProgrammerSalaryBarChartID   = "programmer-salary"
	ProgrammerSalaryBarGraphID   = "programmer-salary"
	EmpyricalDestributionChartID = "distribution"
)

var (
	ProgrammerSalaryGraph = charting.ChartDataset{
		Label:           "Programmer Salary",
		BackgroundColor: []string{charting.Color3, charting.Color4, charting.Color5, charting.Color6, charting.Color7},
		PointRadius:     0,
		ShowLine:        true,
	}

	ProgrammerSalaryChart = charting.Chart{
		ID:          ProgrammerSalaryBarChartID,
		Title:       "Programmer Salary",
		Type:        charting.ChartTypeBar,
		XAxisLabel:  "amount, $",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "people, n",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			ProgrammerSalaryBarGraphID: &ProgrammerSalaryGraph,
		},
	}

	ProgrammerSalaryMeta = ProgrammerSalaryChart.Meta()

	salaryRecords = (*SalaryRecord)(nil)
)

func RenderProgrammerSalary(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("programmer salary chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ';'
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("programmer salary chart: error while decoding csv: %s", err.Error())
		}
	}

	buckets := make([]float64, 5)
	min_salary := math.Inf(1)
	max_salary := math.Inf(-1)
	for i := range salaryRecords.ID {
		min_salary = math.Min(min_salary, salaryRecords.Salary[i])
		max_salary = math.Max(max_salary, salaryRecords.Salary[i])
	}

	bucket_size := (max_salary - min_salary) / float64(len(buckets))
	for i := range salaryRecords.ID {
		if salaryRecords.Position[i] != Programmer {
			continue
		}

		bucket_index := int((salaryRecords.Salary[i] - min_salary) / bucket_size)
		if bucket_index >= len(buckets) {
			bucket_index = len(buckets) - 1
		}
		buckets[bucket_index]++
	}

	x := make([]float64, len(buckets))
	for i := range buckets {
		x[i] = min_salary + bucket_size*float64(i+1)
	}

	copyChart := charting.CopyChart(ProgrammerSalaryChart)
	copyChart.UpdateDataForDataset(ProgrammerSalaryBarGraphID, buckets)

	copyChart.Labels = make([]string, len(buckets))
	for i := range buckets {
		copyChart.Labels[i] = fmt.Sprintf("%.0f-%.0f", x[i], x[i]+bucket_size)
	}

	res = charting.NewRenderResponse()
	res.AddChart(ProgrammerSalaryBarChartID, &copyChart)
	return res
}
