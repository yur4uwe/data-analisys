# 📡 API Contract Specification

This document defines the data structures for the Wails-Chart.js bridge.

## RenderRequest (Go -> JSON)
```go
type RenderRequest struct {
    LabID             string                        `json:"LabID"`
    ChartID           string                        `json:"ChartID"`
    ChartVariables    map[string]float64            `json:"ChartVariables"`
    GraphVariables    map[string]map[string]float64 `json:"GraphVariables"`
    DatasetVisibility map[string]bool               `json:"DatasetVisibility"`
}
```
### Variable Resolution
- `GetChartVariable(chartId, varId)`: Key is `ChartID-VariableID`.
- `GetGraphVariable(chartId, graphId, varId)`: Key is `DatasetID-VariableID`.

## RenderResponse (Go -> JSON)
```go
type RenderResponse struct {
    Charts map[string]Chart `json:"charts"`
    Error  error            `json:"error,omitempty"`
}
```
### Response Methods
- `AddChart(id string, chart *Chart)`: Adds a chart to the response (Use only with copied and populated chart)
- `NewError(message string)`: Creates a new error response
- `NewErrorf(format string, args ...interface{})`: Creates a new formatted error response

## Chart Structure
```go
type Chart struct {
    ID             string             `json:"id"`
    Title          string             `json:"title"`
    Type           GraphType          `json:"type"`
    XAxisLabel     string             `json:"xAxisLabel"`
    YAxisLabel     string             `json:"yAxisLabel"`
    Datasets       map[string]Dataset `json:"datasets"`
    Labels         []string           `json:"labels,omitempty"`
    ChartVariables []MutableField     `json:"chartVariables"`
}
```

## Dataset Interface
```go
type Dataset interface {
	UpdateData([]any) // panics if data is not of the correct type
	UpdateLabel(string) // sets the label of the dataset
	UpdateVariableLabel(int, string) // first argument is the index of the variable in MutableFields slice
	GetData() []any // should be cast to the correct type before use
	GetFields() []MutableField
	Copy() Dataset
	GetBase() *BaseDataset
    GetType() GraphType
}
```
