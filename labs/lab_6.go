package labs

import (
	"labs/charting"
	"labs/labs/cluster"
)

const ()

type Lab6Provider struct{}

var _ charting.LabProvider = Lab6Provider{}

func NewLab6() *Lab6Provider {
	return &Lab6Provider{}
}

func (lp Lab6Provider) GetMetadata() charting.LabMetadata {
	return cluster.Metadata
}

func (lp Lab6Provider) GetConfig() charting.LabConfig {
	return cluster.Config
}

func (lp Lab6Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	res := &charting.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	chart, ok := cluster.Config.Charts[req.ChartID]
	if !ok {
		return res.NewErrorf("chart with id %q not found", req.ChartID)
	}

	return chart.RenderFunc(req)
}
