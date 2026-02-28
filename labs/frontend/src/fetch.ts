import { Render } from "../wailsjs/go/main/App";
import { common } from "../wailsjs/go/models";
import { registry } from "./registry";

export function fetchChartData(): void {
  const labId = window.activeLabId!;
  const chartId = window.activeChartId!;

  if (!labId || !chartId) {
    console.error("Lab or chart ID not found");
    console.error("Lab ID:", labId);
    console.error("Chart ID:", chartId);
    return;
  }

  console.log("Fetching lab data for id:", labId);

  const req = new common.RenderRequest({
    LabID: labId,
    ChartID: chartId,
    GraphVariables: registry.getGraphVariables(labId, chartId),
    ChartVariables: registry.getChartVariables(labId, chartId),
    DatasetVisibility: null,
  });

  // Call Render (now async via goroutine)
  Render(req);
}
