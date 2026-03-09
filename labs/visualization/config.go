package visualization

import (
	"labs/charting"
)

const (
	LabID = "4"
)

var (
	Config = charting.LabConfig{
		Lab: Metadata,
		Charts: map[string]*charting.Chart{
			BarChartID:      &BarChart,
			FunctionChartID: &FunctionChart,
			LinearChartID:   &LinearChart,
			RadialChartID:   &RadialChart,
		},
	}

	Metadata = charting.LabMetadata{
		ID:   LabID,
		Name: "Visualization",
		Charts: map[string]charting.ChartMetadata{
			BarChartID:      BarMeta,
			FunctionChartID: FunctionMeta,
			LinearChartID:   LinearMeta,
			RadialChartID:   RadialMeta,
		},
	}
)
