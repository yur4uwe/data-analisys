import { Chart, ChartTypeRegistry, ScaleChartOptions } from "chart.js";
import { MatrixController, MatrixElement } from "chartjs-chart-matrix";
import ChartDataLabels from "chartjs-plugin-datalabels";
import ZoomPlugin from "chartjs-plugin-zoom";
import { charting } from "../wailsjs/go/models";
import { defaultChartOptions, newScales, processDataset } from "./static-config";

Chart.register(MatrixController, MatrixElement, ChartDataLabels, ZoomPlugin);

// Helper for color conversion in legend
function hexToRgb(hex: string) {
  hex = hex.replace("#", "");
  if (hex.length === 3) {
    hex = hex.split("").map((s) => s + s).join("");
  }
  const num = parseInt(hex, 16);
  return {
    r: (num >> 16) & 255,
    g: (num >> 8) & 255,
    b: num & 255,
  };
}

// Custom Heatmap Legend Plugin
const heatmapLegendPlugin = {
  id: "heatmapLegend",
  afterDraw: (chart: any) => {
    const chartType = chart.config.type;
    if (chartType !== "heatmap" && chartType !== "multi-heatmap") return;

    const ctx = chart.ctx;
    const area = chart.chartArea;
    const dataset = chart.data.datasets[0];
    if (!dataset || !dataset.data.length) return;

    // Retrieve min/max values and colors
    const values = dataset.data.map((d: any) => d.v);
    const min = Math.min(...values);
    const max = Math.max(...values);
    
    // Find the colors used for the heatmap (need to reach back to the processDataset logic or re-calculate)
    // For simplicity, we'll assume the interpolateColor logic is accessible or we use the colors from the dataset config if we can store them.
    // Since we can't easily import from static-config here without circular deps, 
    // we'll use a standard approach: pull from the dataset if we stored it there.
    const colors = dataset.backgroundColorList || ["#1d4ed8", "#b91c1c"];

    const legendWidth = 200;
    const legendHeight = 12;
    const x = area.left + (area.right - area.left - legendWidth) / 2;
    const y = chart.height - 40; // Relative to canvas bottom

    // Draw gradient bar
    const gradient = ctx.createLinearGradient(x, 0, x + legendWidth, 0);
    const gamma = 0.3; // Match the sensitivity in static-config.ts
    
    // Add more stops to simulate the power curve
    const numStops = 10;
    for (let i = 0; i <= numStops; i++) {
      const t = i / numStops;
      // We need to "invert" the gamma for the legend so the colors align
      // The cells use: normalizedColor = Math.pow(normalizedValue, gamma)
      // So at position 't' on the legend, we want the color for value 't^gamma'
      // Wait, let's just draw the colors at their actual mapped positions
      const mappedT = Math.pow(t, gamma);
      
      // Since our interpolateColor logic is in static-config.ts and hard to access here,
      // we interpolate between the first and last colors directly for the legend.
      // If we have more than 2 colors, we'd need a more complex loop.
      const startColor = hexToRgb(colors[0]);
      const endColor = hexToRgb(colors[colors.length - 1]);
      
      const r = Math.round(startColor.r + (endColor.r - startColor.r) * mappedT);
      const g = Math.round(startColor.g + (endColor.g - startColor.g) * mappedT);
      const b = Math.round(startColor.b + (endColor.b - startColor.b) * mappedT);
      
      gradient.addColorStop(t, `rgb(${r},${g},${b})`);
    }

    ctx.fillStyle = gradient;
    ctx.fillRect(x, y, legendWidth, legendHeight);

    // Draw labels
    ctx.fillStyle = "#000000";
    ctx.font = "bold 11px Arial";
    ctx.textAlign = "center";
    const formatValue = (v: number) => {
      if (Math.abs(v) < 0.0001 || Math.abs(v) > 10000) {
        return v.toExponential(2);
      }
      return v.toFixed(4);
    };
    ctx.fillText(formatValue(min), x, y + legendHeight + 14);
    ctx.fillText(formatValue(max), x + legendWidth, y + legendHeight + 14);
  }
};

Chart.register(heatmapLegendPlugin);

// Store chart instances globally for access
declare global {
  interface Window {
    chartInstances: Map<string, Chart>;
  }
}

if (!window.chartInstances) {
  window.chartInstances = new Map();
}

export function getDataLabels(
  pointLabels: string[] | undefined,
  chartType: keyof ChartTypeRegistry,
) {
  if (Array.isArray(pointLabels) && pointLabels.length > 0) {
    return {
      display: true,
      formatter: (_: any, ctx: any) => {
        if (!pointLabels || ctx.dataIndex >= pointLabels.length) return "";
        return pointLabels[ctx.dataIndex] ?? "";
      },
      anchor: "end" as const,
      align: "top" as const,
      color: "#000000",
      font: { weight: "bold" as const, size: 11 },
      backgroundColor: "rgba(255,255,255,0.85)",
      borderRadius: 3,
      padding: 4,
    };
  }

  switch (chartType) {
    case "pie":
    case "doughnut":
      return {
        display: true,
        color: "#ffffff",
        font: { weight: "bold" as const, size: 14 },
        formatter: (value: number, ctx: any) => {
          const data = ctx.dataset.data as number[];
          if (!data || data.length === 0) return value;
          const total = data.reduce((a, b) => a + b, 0);
          if (total === 0) return value;
          const pct = ((value / total) * 100).toFixed(1);
          return `${value}\n(${pct}%)`;
        },
      };

    case "bar":
      return {
        display: false,
        color: "#ffffff",
        font: { weight: "bold" as const, size: 14 },
      };
    default:
      return {
        display: false,
      };
  }
}

export function renderChartInto(chartConfig: charting.Chart, container: HTMLElement) {
  if (!chartConfig) {
    console.error("renderChartInto: chartConfig is null or undefined!");
    return;
  }

  // Clear previous content
  container.innerHTML = "";

  const canvas = document.createElement("canvas");
  canvas.id = `chart-${chartConfig.id}`;
  container.appendChild(canvas);

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    console.error("renderChartInto: Canvas context is null!");
    return;
  }

  const chartType = (chartConfig.type as keyof ChartTypeRegistry) || "line";
  console.log("Rendering chart ID:", chartConfig.id, "Type:", chartType);

  const hasScales = !["pie", "doughnut", "polarArea"].includes(chartType);
  const hasContinuousAxes = ["scatter", "line", "bubble"].includes(chartType);

  let chartLabels: string[] = Array.isArray(chartConfig.labels) ? chartConfig.labels : [];

  // Process datasets based on chart type
  const processedDatasets = Object.values(chartConfig.datasets || {}).map(processDataset(hasScales, chartType));

  let maxDataLength = 0

  for (const dataset of processedDatasets) {
    maxDataLength = Math.max(dataset.data.length, maxDataLength);
  }

  if (maxDataLength < chartLabels.length) {
    chartLabels = chartLabels.slice(0, maxDataLength);
  }

  const chartOptions: any = defaultChartOptions(chartConfig.title || "", chartType);

  // Only add scales for charts that use them
  if (hasScales) {
    chartOptions.scales = newScales(
      chartConfig,
      hasContinuousAxes,
    ) as unknown as ScaleChartOptions;
  }

  console.log(`Initializing Chart.js for ${chartConfig.id} with ${chartLabels.length} labels`);

  try {
    const chart = new Chart(ctx, {
      type: chartType,
      data: {
        labels: chartLabels,
        datasets: processedDatasets as any[],
      },
      options: chartOptions,
    });
    window.chartInstances.set(chartConfig.id, chart);
  } catch (err) {
    console.error(`Failed to initialize Chart.js for ${chartConfig.id}:`, err);
  }

  const resetZoom = document.getElementById("reset-zoom-btn");
  if (resetZoom) {
    const resetZoomCallback = () => {
      const instance = window.chartInstances.get(chartConfig.id);
      if (instance) {
        (instance as any).resetZoom();
      }
    };
    resetZoom.onclick = resetZoomCallback;
  }
}

export function renderMultiChart(chartConfig: charting.Chart) {
  if (!chartConfig || !chartConfig.datasets) {
    console.error("renderMultiChart: chartConfig or datasets is missing");
    return;
  }

  console.log("Rendering multiple charts for:", chartConfig.id);
  const container = document.getElementById("chart-container")!;
  if (!container) {
    console.error("renderMultiChart: chart-container not found");
    return;
  }

  container.innerHTML = "";
  container.style.display = "grid";
  container.style.gridTemplateColumns = "repeat(auto-fit, minmax(320px, 1fr))";
  container.style.gap = "16px";

  const singleType = (chartConfig.type || "").replace("multi-", "") as keyof ChartTypeRegistry;
  const labels = Array.isArray(chartConfig.labels) ? chartConfig.labels : [];

  Object.entries(chartConfig.datasets).forEach(([datasetId, dataset]) => {
    if (!dataset) return;

    const wrapper = document.createElement("div");
    wrapper.className = "chart-wrapper";
    wrapper.style.position = "relative";
    wrapper.style.minHeight = "400px";
    wrapper.style.height = "400px";
    wrapper.style.width = "100%";
    container.appendChild(wrapper);

    // Each cluster should ideally have its own labels matching its data length
    let clusterLabels = labels;
    if (!labels && dataset.pointLabels && dataset.pointLabels.length > 0) {
      clusterLabels = dataset.pointLabels;
    }

    // Synthetic single-dataset chart reusing all config from parent
    const syntheticChart: charting.Chart = charting.Chart.createFrom({
      ...chartConfig,
      labels: clusterLabels,
      type: singleType,
      id: `${chartConfig.id}-${datasetId}`,
      title: dataset.label || datasetId,
      datasets: { [datasetId]: dataset },
    });

    renderChartInto(syntheticChart, wrapper);
  });
}

// Helper functions for dataset control
export function toggleDatasetVisibility(chartId: string, datasetIndex: number) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets[datasetIndex].hidden =
      !chart.data.datasets[datasetIndex].hidden;

    chart.update();
  }
}

export function updateDatasets(chartId: string, newDatasets: any[]) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets = newDatasets;
    chart.update();
  }
}

export function addDataset(chartId: string, newDataset: any) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets.push(newDataset);
    chart.update();
  }
}

export function removeDataset(chartId: string, datasetIndex: number) {
  const chart = window.chartInstances.get(chartId);
  if (chart) {
    chart.data.datasets.splice(datasetIndex, 1);
    chart.update();
  }
}
