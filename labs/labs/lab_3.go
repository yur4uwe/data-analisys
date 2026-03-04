package labs

import (
	"labs/charting"
	"labs/labs/polyapprox"
	"labs/labs/render"
)

type Lab3Provider struct{}

var _ charting.LabProvider = Lab3Provider{}

func NewLab3() *Lab3Provider {
	return &Lab3Provider{}
}

func (lp Lab3Provider) GetMetadata() charting.LabMetadata {
	return polyapprox.Metadata
}

func (lp Lab3Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	if req == nil {
		return &charting.RenderResponse{Error: render.NewRenderError("empty render request")}
	}

	switch req.ChartID {
	case polyapprox.RandomFitsID:
		return polyapprox.RenderRandomFits(req)
	case polyapprox.SampleDataID:
		return polyapprox.RenderSampleData(req)
	case polyapprox.RandomMSEID:
		return polyapprox.RenderRandomPolynomialMSE(req)
	case polyapprox.SampleMSEID:
		return polyapprox.RenderSamplePolynomialMSE(req)
	default:
		return &charting.RenderResponse{Error: render.NewRenderErrorf("unrecognised Chart: %s", req.ChartID)}
	}
}

func (lp Lab3Provider) GetConfig() charting.LabConfig {
	return charting.LabConfig{
		Lab: polyapprox.Metadata,
		Charts: map[string]*charting.Chart{
			polyapprox.RandomFitsID: &polyapprox.RandomFitsChart,
			polyapprox.SampleDataID: &polyapprox.SampleDataChart,
			polyapprox.RandomMSEID:  &polyapprox.RandomMSEChart,
			polyapprox.SampleMSEID:  &polyapprox.SampleMSEChart,
		},
	}
}
