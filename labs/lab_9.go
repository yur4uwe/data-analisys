package labs

import (
	"labs/charting"
	forecastinglinparab "labs/labs/forecasting-lin-parab"
	"labs/labs/render"
)

type Lab9Provider struct{}

var _ charting.LabProvider = Lab9Provider{}

func NewLab9() *Lab9Provider {
	return &Lab9Provider{}
}

func (lp Lab9Provider) GetMetadata() charting.LabMetadata {
	return forecastinglinparab.LinParabConfig.Meta()
}

func (lp Lab9Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	return &charting.RenderResponse{Error: render.NewRenderError("not implemented")}
}

func (lp Lab9Provider) GetConfig() charting.LabConfig {
	return forecastinglinparab.LinParabConfig
}
