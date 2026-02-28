package polyapprox

import (
	"fmt"
	"labs/labs/common"
	"labs/labs/render"
	"math/rand"
)

const (
	LabID = "3"

	RandomFitsID = "random-fits"

	OriginalDataID = "orig-data"
	NoisyDataID    = "noisy-data"
	LinearApproxID = "linear-approx"
	QuadApproxID   = "quad-approx"

	IntervalStartID = "start-interval"
	IntervalEndID   = "end-interval"
	IntervalStepID  = "step"
	NoiseAmpID      = "noise-amplifier"

	LinearFitCoefficientsID    = "linear-fit-coefficients"
	QuadraticFitCoefficientsID = "quadratic-fit-coefficients"
)

var (
	ChartVariables = []common.MutableField{
		{
			ID:      IntervalStartID,
			Label:   "Start",
			Default: -5.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    0.5,
			Control: common.ControlNumber,
		},
		{
			ID:      IntervalEndID,
			Label:   "End",
			Default: 5.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    0.5,
			Control: common.ControlNumber,
		},
		{
			ID:      IntervalStepID,
			Label:   "Step",
			Default: 0.1,
			Min:     0.01,
			Max:     1,
			Step:    0.01,
			Control: common.ControlRange,
		},
		{
			ID:      NoiseAmpID,
			Label:   "Noise Amplifier",
			Default: 1,
			Min:     1,
			Max:     100,
			Step:    1,
			Control: common.ControlRange,
		},
	}

	LinearFitCoefficients = common.MutableField{
		ID:      LinearFitCoefficientsID,
		Label:   "Linear Fit Coefficients",
		Control: common.ControlNoControl,
	}

	QuadraticFitCoefficients = common.MutableField{
		ID:      QuadraticFitCoefficientsID,
		Label:   "Quadratic Fit Coefficients",
		Control: common.ControlNoControl,
	}

	OriginalData = common.ChartDataset{
		Label:           "Original",
		Type:            common.ChartTypeLine,
		BorderColor:     "#2563eb",
		BackgroundColor: "rgba(37, 99, 235, 0.1)",
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	NoisyData = common.ChartDataset{
		Label:           "Noisy",
		Type:            common.ChartTypeLine,
		BorderColor:     "#dc2626",
		BackgroundColor: "rgba(220, 38, 38, 0.1)",
		PointRadius:     2,
		BorderWidth:     1,
		ShowLine:        true,
		Togglable:       true,
	}

	LinearApprox = common.ChartDataset{
		Label:           "Linear Approximation",
		Type:            common.ChartTypeLine,
		BorderColor:     "#16a34a",
		BackgroundColor: "rgba(22, 163, 74, 0.1)",
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  []common.MutableField{LinearFitCoefficients},
	}

	QuadApprox = common.ChartDataset{
		Label:           "Quadratic Approximation",
		Type:            common.ChartTypeLine,
		BorderColor:     "#9333ea",
		BackgroundColor: "rgba(147, 51, 234, 0.1)",
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  []common.MutableField{QuadraticFitCoefficients},
	}

	RandomFitsChart = common.Chart{
		ID:          RandomFitsID,
		Title:       "Random Data Fits",
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			OriginalDataID: &OriginalData,
			NoisyDataID:    &NoisyData,
			LinearApproxID: &LinearApprox,
			QuadApproxID:   &QuadApprox,
		},
		ChartVariables: ChartVariables,
	}

	RandomFitsMetadata = RandomFitsChart.Meta()

	Metadata = common.LabMetadata{
		ID:   LabID,
		Name: "Least Squares Approximation",
		Charts: map[string]common.ChartMetadata{
			RandomFitsID: RandomFitsMetadata,
			SampleDataID: SampleDataMetadata,
			RandomMSEID:  RandomMSEMetadata,
			SampleMSEID:  SampleMSEMetadata,
		},
	}
)

func RenderRandomFits(req *common.RenderRequest) *common.RenderResponse {
	start, hasStart := req.GetChartVariable(RandomFitsID, IntervalStartID)
	end, hasEnd := req.GetChartVariable(RandomFitsID, IntervalEndID)
	step, hasStep := req.GetChartVariable(RandomFitsID, IntervalStepID)
	noiseAmp, hasNoise := req.GetChartVariable(RandomFitsID, NoiseAmpID)

	if !hasStart {
		start = ChartVariables[0].Default
	}
	if !hasEnd {
		end = ChartVariables[1].Default
	}
	if !hasStep {
		step = ChartVariables[2].Default
	}
	if !hasNoise {
		noiseAmp = ChartVariables[3].Default
	}

	if step <= 0 {
		return &common.RenderResponse{Error: render.NewRenderError("step must be greater than 0")}
	}
	if start > end {
		return &common.RenderResponse{Error: render.NewRenderError("start interval must be less than or equal to end interval")}
	}

	seed := int64(230420067)
	x, y, origY := GenerateRandomSeries(start, end, step, noiseAmp, seed)

	if len(x) == 0 {
		return &common.RenderResponse{Error: render.NewRenderError("no data generated with given parameters")}
	}

	chartCopy := common.CopyChart(RandomFitsChart)

	chartCopy.UpdatePointsForDataset(OriginalDataID, x, origY)
	chartCopy.UpdatePointsForDataset(NoisyDataID, x, y)

	if coefs, err := SolvePolynomialFit(x, y, 1); err == nil {
		approx := make([]float64, 0, len(x))
		for _, xi := range x {
			approx = append(approx, EvaluatePolynomial(coefs, xi))
		}
		chartCopy.UpdatePointsForDataset(LinearApproxID, x, approx)
		chartCopy.Datasets[LinearApproxID].GraphVariables[0].Label = fmt.Sprintf("Linear Fit Coefficients (a=%.4f, b=%.4f) for y=bx+a", coefs[0], coefs[1])
	} else {
		fmt.Println("linear fit failed:", err)
	}

	if coefs, err := SolvePolynomialFit(x, y, 2); err == nil {
		approx := make([]float64, 0, len(x))
		for _, xi := range x {
			approx = append(approx, EvaluatePolynomial(coefs, xi))
		}
		chartCopy.UpdatePointsForDataset(QuadApproxID, x, approx)
		chartCopy.Datasets[QuadApproxID].GraphVariables[0].Label = fmt.Sprintf("Quadratic Fit Coefficients (a=%.4f, b=%.4f, c=%.4f) for y=cx^2+bx+a", coefs[0], coefs[1], coefs[2])
	} else {
		fmt.Println("quadratic fit failed:", err)
	}

	return &common.RenderResponse{
		Charts: map[string]common.Chart{
			RandomFitsID: chartCopy,
		},
	}
}

func GenerateRandomSeries(start, end, step, noiseAmp float64, seed int64) ([]float64, []float64, []float64) {
	r := rand.New(rand.NewSource(seed))
	n := int((end-start)/step) + 1
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)
	origY := make([]float64, 0, n)

	for i := start; i <= end; i += step {
		noise := r.NormFloat64() * 0.2 * noiseAmp
		curr := 0.8 - 4*i
		x = append(x, i)
		y = append(y, curr+noise)
		origY = append(origY, curr)
	}

	return x, y, origY
}
