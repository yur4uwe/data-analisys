package polyapprox

import (
	"encoding/csv"
	"fmt"
	"io"
	"labs/labs/common"
	"labs/labs/render"
	"math"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	SampleDataID = "sample-data"

	sampleDataGraphID          = "sample"
	sampleApproximationGraphID = "sample-approx"

	approximationDegree = "degree"
	coeffsDisplayID     = "coeffs"
)

var (
	approxDegreeVariable = common.MutableField{
		ID:      approximationDegree,
		Label:   "Degree of Polynomial",
		Default: 1,
		Min:     0,
		Max:     10,
		Step:    1,
		Control: common.ControlRange,
	}

	coeffsDisplayVariable = common.MutableField{
		ID:      coeffsDisplayID,
		Label:   "Polynomial coefficients: ",
		Control: common.ControlNoControl,
	}

	sampleDataGraph = common.ChartDataset{
		Label:           "Sample Data",
		BorderColor:     common.Color10,
		BackgroundColor: []string{"rgba(0, 0, 0, 0.1)"},
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  []common.MutableField{},
	}

	sampleDataApproxGraph = common.ChartDataset{
		Label:           "Sample Data Approximation",
		BorderColor:     common.Color6,
		BackgroundColor: []string{"rgba(0, 0, 0, 0.1)"},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []common.MutableField{
			approxDegreeVariable,
			coeffsDisplayVariable,
		},
	}

	SampleDataChart = common.Chart{
		ID:          SampleDataID,
		Title:       "Sample Data (CSV)",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			OriginalDataID:             &sampleDataGraph,
			sampleApproximationGraphID: &sampleDataApproxGraph,
		},
		ChartVariables: []common.MutableField{},
	}

	SampleDataMetadata = SampleDataChart.Meta()
)

func sortXandY(x, y []float64) {
	slices.Sort(x)

	sort.SliceStable(y, func(i, j int) bool {
		return x[i] < x[j]
	})
}

func RenderSampleData(req *common.RenderRequest) *common.RenderResponse {
	x, y, err := ReadSampleCSV("../data/lab_3_var_12.csv")
	if err != nil {
		fmt.Println("failed to open file:", err)
		return &common.RenderResponse{
			Error: render.NewRenderError("failed to read sample data file"),
		}
	}

	chartCopy := common.CopyChart(SampleDataChart)
	chartCopy.UpdatePointsForDataset(OriginalDataID, x, y)

	degree, ok := req.GetGraphVariable(SampleDataID, sampleApproximationGraphID, approximationDegree)
	if !ok {
		degree = 2.0
	}

	coeffs, err := SolvePolynomialFit(x, y, int(degree))
	if err != nil {
		return &common.RenderResponse{
			Error: render.NewRenderErrorf("failed to solve polynomial fit: %v", err),
		}
	}

	minX, maxX := math.Inf(1), math.Inf(-1)
	for _, xi := range x {
		maxX = max(maxX, xi)
		minX = min(minX, xi)
	}

	step := (maxX - minX) / float64(len(x)-1)

	approx := make([]float64, 0, len(x))
	for i := minX; i < maxX; i += step {
		approx = append(approx, EvaluatePolynomial(coeffs, i))
	}

	chartCopy.UpdatePointsForDataset(sampleApproximationGraphID, x, approx)

	var str strings.Builder
	str.WriteString("Polynomial Coefficients (")
	for i, c := range coeffs {
		fmt.Fprintf(&str, "x%d=%.2f", i, c)
		if i != len(coeffs)-1 {
			fmt.Fprint(&str, ", ")
		}
	}
	str.WriteString(")")
	chartCopy.Datasets[sampleApproximationGraphID].GraphVariables[1].Label = str.String()

	return &common.RenderResponse{
		Charts: map[string]common.Chart{
			SampleDataID: chartCopy,
		},
	}
}
func ReadSampleCSV(filename string) ([]float64, []float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Skip header
	reader.Comma = ','

	var xVals, yVals []float64
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, nil, err
		}
		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, nil, err
		}

		xVals = append(xVals, x)
		yVals = append(yVals, y)
	}

	return xVals, yVals, nil
}
