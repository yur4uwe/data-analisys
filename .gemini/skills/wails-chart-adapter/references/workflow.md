# 🛠 Workflow: Adding a New Lab

The architecture uses a `GenericProvider` to eliminate boilerplate. Follow these steps:

## 1. Create the Lab Package
Create a new directory (e.g., `labs/my_lab/`). Inside, create a `config.go` or similar.

## 2. Define the Config and Chart
```go
package my_lab

import "labs/charting"

const LabID = "my-lab-id"
var MyChart = charting.Chart{
    ID: "main-chart",
    Title: "Visualization",
    Type: charting.GraphTypeLine,
    Datasets: map[string]charting.Dataset{
        "ds_1": &charting.GridDataset{ ... },
    },
}

var Config = charting.NewLabConfig(LabID, "Lab Name", map[string]*charting.Chart{
    MyChart.ID: &MyChart,
})
```

## 3. Implement and Assign logic
Define your render function and register it in an `init()` block:
```go
func Render(req *charting.RenderRequest) *charting.RenderResponse {
    res := charting.NewRenderResponse()
    copyChart := charting.CopyChart(MyChart)
    // Update data logic
    res.AddChart(MyChart.ID, &copyChart)
    return res
}

func init() {
    MyChart.RenderFunc = Render
}
```

## 4. Register in `app.go`
Simply register using the `GenericProvider` constructor:
```go
a.registry[my_lab.LabID] = charting.NewProvider(my_lab.Config)
```

## 🧪 Best Practices
- **Immutable Templates**: Always use `charting.CopyChart(Template)` inside `RenderFunc`.
- **Error Handling**: Use `res.NewErrorf("message", err)` to report backend failures to the UI.
