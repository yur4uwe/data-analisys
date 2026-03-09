package charting

type LabProvider interface {
	GetMetadata() LabMetadata
	GetConfig() LabConfig
	Render(req *RenderRequest) *RenderResponse
}
