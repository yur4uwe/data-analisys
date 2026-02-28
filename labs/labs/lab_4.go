package labs

import (
	"labs/labs/common"
	"labs/labs/visualization"
)

func NewLab4() *Lab4Provider {
	return &Lab4Provider{}
}

type Lab4Provider struct{}

var _ common.LabProvider = Lab4Provider{}

func (lp Lab4Provider) GetMetadata() common.LabMetadata {
	return visualization.Metadata
}

func (lp Lab4Provider) GetConfig() common.LabConfig {
	return visualization.Config
}

func (lp Lab4Provider) Render(req *common.RenderRequest) *common.RenderResponse {
	res := &common.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	switch req.ChartID {
	case visualization.BarChartID:
		return visualization.RenderBarPlot(req)
	case visualization.FunctionChartID:
		return visualization.RenderFunction(req)
	case visualization.LinearChartID:
		return visualization.RenderLinear(req)
	case visualization.RadialChartID:
		return visualization.RenderRadialPlot(req)
	default:
		return res.NewErrorf("unrecognized chart ID: %s", req.ChartID)
	}
}
