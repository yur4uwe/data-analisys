import { Chart } from "chart.js";
import { common } from "../wailsjs/go/models";

// Store chart instances globally for access
declare global {
  interface Window {
    chartInstances: Map<string, Chart>;
  }
}

if (!window.chartInstances) {
  window.chartInstances = new Map();
}

export function renderChart(chartConfig: common.Chart) {
  let container = document.getElementById("chart-container");
  if (!container) {
    container = document.createElement("div");
    container.id = "chart-container";
    document.body.appendChild(container);
  }

  // Clear previous content
  container.innerHTML = "";

  if (!chartConfig) {
    console.log("Charts data is null or undefined!");
    return;
  }

  const canvasWrapper = document.createElement("div");
  canvasWrapper.className = "chart-wrapper";

  const canvas = document.createElement("canvas");
  canvas.id = `chart-${chartConfig.id}`;

  canvasWrapper.appendChild(canvas);
  container!.appendChild(canvasWrapper);

  // Initialize Chart.js
  const ctx = canvas.getContext("2d");
  if (!ctx) {
    console.log("Canvas context is null or undefined!");
    return;
  }
  const chart = new Chart(ctx, {
    type: "line",
    data: {
      labels: chartConfig.datasets["orig-data"].pointData
        ? chartConfig.datasets["orig-data"].pointData.map(
            ({ x }: { x: number; y: number }) => x.toFixed(2),
          )
        : [],
      datasets: Object.values(chartConfig.datasets).map((dataset) => ({
        label: dataset.label,
        data: dataset.pointData || [],
        borderColor: dataset.borderColor ?? "#666",
        backgroundColor: dataset.backgroundColor ?? "transparent",
        tension: dataset.tension ?? 0,
        fill: dataset.fill ?? false,
        hidden: dataset.hidden ?? false,
        pointRadius: dataset.pointRadius ?? 0,
        borderWidth: dataset.borderWidth ?? 2,
        showLine: dataset.showLine !== false,
      })),
    },
    options: {
      responsive: true,
      plugins: {
        title: {
          display: true,
          text: chartConfig.title,
          color: "#ffffff",
          font: {
            size: 18,
            weight: "bold",
          },
          padding: {
            top: 10,
            bottom: 20,
          },
        },
        legend: {
          labels: {
            color: "#ffffff",
            font: {
              size: 13,
            },
            padding: 15,
            usePointStyle: true,
          },
        },
        tooltip: {
          backgroundColor: "rgba(0, 0, 0, 0.9)",
          titleColor: "#ffffff",
          bodyColor: "#ffffff",
          borderColor: "#ffffff",
          borderWidth: 1,
          padding: 12,
          displayColors: true,
        },
      },
      scales: {
        x: {
          title: {
            display: true,
            text: chartConfig.xAxisLabel,
            color: "#ffffff",
            font: {
              size: 14,
              weight: "bold",
            },
          },
          ticks: {
            color: "#ffffff",
            font: {
              size: 12,
            },
          },
          grid: {
            color: "rgba(255, 255, 255, 0.2)",
          },
        },
        y: {
          title: {
            display: true,
            text: chartConfig.yAxisLabel,
            color: "#ffffff",
            font: {
              size: 14,
              weight: "bold",
            },
          },
          ticks: {
            color: "#ffffff",
            font: {
              size: 12,
            },
          },
          grid: {
            color: "rgba(255, 255, 255, 0.2)",
          },
        },
      },
    },
  });

  // Store chart instance for later access
  window.chartInstances.set(chartConfig.id, chart);
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
