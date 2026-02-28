import { Chart } from "chart.js";
import { fetchChartData } from "./fetch";
import { registry } from "./registry";
import { InitializeChart, InitializeLab } from "./lab-init";

export class ActiveLabChangeEvent extends Event {
  constructor(public labId: string) {
    super("activeLabChange");
    console.log("ActiveLabChangeEvent created with id: ", labId);
    this.labId = labId;
  }
}

export class ActiveChartChangeEvent extends Event {
  constructor(public chartId: string) {
    super("activeChartChange");
    this.chartId = chartId;
  }
}

export class RerenderEvent extends Event {
  constructor(public chartId: string) {
    super("rerender");
    this.chartId = chartId;
  }
}

// Extend Window interface to recognize custom event
declare global {
  interface Window {
    activeChartId: string | null;
    activeLabId: string | null;
    chartInstances: Map<string, Chart>;
  }
  interface WindowEventMap {
    activeLabChange: ActiveLabChangeEvent;
    activeChartChange: ActiveChartChangeEvent;
    rerender: RerenderEvent;
  }
}

if (!window.chartInstances) {
  window.chartInstances = new Map<string, Chart>();
}

window.addEventListener(
  "activeChartChange",
  (event: ActiveChartChangeEvent) => {
    window.activeChartId = event.chartId;
    const error = document.getElementById("error-container");
    if (error) {
      error.innerHTML = "";
    }
    InitializeChart(event.chartId);
    fetchChartData();
  },
);

window.addEventListener("activeLabChange", (event: ActiveLabChangeEvent) => {
  window.activeLabId = event.labId;

  console.log("Active lab change event:", event.labId);

  // Initialize the lab (create chart tabs, etc.)
  InitializeLab(event.labId);

  // Get the lab from registry and select the first chart automatically
  const lab = registry.getLab(event.labId);
  if (lab && Object.keys(lab.Charts).length > 0) {
    const firstChartId = Object.keys(lab.Charts)[0];
    window.activeChartId = firstChartId;
    console.log(
      "Active lab changed to:",
      window.activeLabId,
      "| First chart selected:",
      firstChartId,
    );
    window.dispatchEvent(new ActiveChartChangeEvent(firstChartId));
  } else {
    window.activeChartId = null;
    console.log(
      "Active lab changed to:",
      window.activeLabId,
      "| No charts available",
    );
  }
});

window.addEventListener("rerender", () => {
  console.log("Rerendering...");
  fetchChartData();
});
