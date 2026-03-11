package cluster

import (
	"labs/charting"
)

const (
	SimpleChartID = "simple"

	VariableThresholdID = "threshold"
)

var (
	VariableThreshold = charting.MutableField{
		ID:      VariableThresholdID,
		Label:   "Distance Threshold",
		Default: 5,
		Min:     0.0,
		Max:     10,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	SimpleChart = charting.Chart{
		ID:          SimpleChartID,
		Title:       "Simple Clusterization",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "X",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "Y",
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableThreshold,
		},
	}
)

func RenderSimple(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadPoints(); err != nil {
		return res.NewError(err.Error())
	}

	threshold, ok := req.GetChartVariable(SimpleChartID, VariableThresholdID)
	if !ok {
		threshold = VariableThreshold.Default
	}

	copyChart := charting.CopyChart(SimpleChart)

	labels, centroids := simpleClustering(points, threshold)

	clusterData(labels, len(centroids), &copyChart)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)

	return res
}

// simpleClustering assigns each point to the nearest existing cluster center,
// or creates a new cluster if all distances exceed the threshold T.
// Returns cluster labels (one per point) and the final cluster centers.
func simpleClustering(points []charting.DataPoint, T float64) (labels []int, centroids []charting.DataPoint) {
	if len(points) == 0 {
		return nil, nil
	}

	labels = make([]int, len(points))
	centroids = []charting.DataPoint{points[0]}

	for i, p := range points[1:] {
		best, minDist := 0, euclidianDist(p, centroids[0])
		for j := 1; j < len(centroids); j++ {
			if d := euclidianDist(p, centroids[j]); d < minDist {
				minDist, best = d, j
			}
		}

		if minDist > T {
			labels[i+1] = len(centroids)
			centroids = append(centroids, p)
		} else {
			labels[i+1] = best
		}
	}

	return labels, centroids
}
