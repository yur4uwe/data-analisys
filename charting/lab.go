package charting

// LabMetadata contains information about a single lab
type LabMetadata struct {
	ID     string
	Name   string
	Charts map[string]ChartMetadata
}

// GetLabsResponse returns all available labs for the UI
type GetLabsResponse struct {
	Labs []LabMetadata `json:"labs"`
}

// LabConfig is the complete configuration for a lab
type LabConfig struct {
	Lab    LabMetadata       `json:"lab"`
	Charts map[string]*Chart `json:"charts"`
}

// GenericProvider implements LabProvider for standard lab implementations
type GenericProvider struct {
	Config LabConfig
}

func NewProvider(config LabConfig) *GenericProvider {
	return &GenericProvider{Config: config}
}

func (p *GenericProvider) GetMetadata() LabMetadata {
	return p.Config.Lab
}

func (p *GenericProvider) GetConfig() LabConfig {
	return p.Config
}

func (p *GenericProvider) Render(req *RenderRequest) *RenderResponse {
	res := NewRenderResponse()
	if req == nil {
		return res.NewError("request is nil")
	}

	chart, ok := p.Config.Charts[req.ChartID]
	if !ok {
		return res.NewErrorf("chart with id %q not found", req.ChartID)
	}

	if chart.RenderFunc == nil {
		return res.NewErrorf("chart with id %q has no render function defined", req.ChartID)
	}

	return chart.RenderFunc(req)
}

func NewLabConfig(labID, labName string, charts map[string]*Chart) LabConfig {
	chartsMeta := make(map[string]ChartMetadata, len(charts))
	for id, chart := range charts {
		chartsMeta[id] = chart.Meta()
	}

	return LabConfig{
		Lab: LabMetadata{
			ID:     labID,
			Name:   labName,
			Charts: chartsMeta,
		},
		Charts: charts,
	}
}

func (lc *LabConfig) Meta() LabMetadata {
	return lc.Lab
}
