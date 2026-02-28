package common

type LabProvider interface {
	GetMetadata() LabMetadata
	GetConfig() LabConfig
	Render(req *RenderRequest) *RenderResponse
}
