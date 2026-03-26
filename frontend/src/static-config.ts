import { charting } from "../wailsjs/go/models";
import { getDataLabels } from "./chart-render";
import { Dataset, isHeatmapDataset } from "./types";

export const defaultChartOptions = (title: string, chartType?: string) => ({
	responsive: true,
	maintainAspectRatio: false,
	resizeDelay: 100,
	animation: false,
	layout: {
		padding: {
			bottom: (chartType === "heatmap" || chartType === "multi-heatmap") ? 60 : 10,
		},
	},
	plugins: {
		title: {
			display: !!title,
			text: title || "",
			color: "#000000",
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
			display: (ctx: any) => {
				const chartType = ctx.chart.config.type;
				return chartType !== "heatmap" && chartType !== "multi-heatmap";
			},
			labels: {
				color: "#000000",
				font: {
					size: 13,
				},
				padding: 15,
				usePointStyle: true,
			},
		},
		tooltip: {
			enabled: true,
			backgroundColor: "rgba(0, 0, 0, 0.9)",
			titleColor: "#ffffff",
			bodyColor: "#ffffff",
			borderColor: "#ffffff",
			borderWidth: 1,
			padding: 12,
			displayColors: true,
			callbacks: {
				label: (ctx: any) => {
					if (ctx.dataset.type === 'matrix') {
						const d = ctx.dataset.data[ctx.dataIndex];
						if (d) {
							return `MSE: ${d.v.toPrecision(6)} (α: ${d.x.toFixed(2)}, β: ${d.y.toFixed(2)})`;
						}
					}
					return ctx.formattedValue;
				}
			}
		},
		zoom: {
			zoom: {
				wheel: {
					enabled: true,
					speed: 0.02,
					modifierKey: "ctrl",
				},
				pinch: { enabled: true },
				mode: "xy",
			},
			pan: {
				enabled: true,
				mode: "xy",
			},
		},
	},
})

// Helper for color interpolation in heatmaps
function interpolateColor(value: number, colors: string[]): string {
	if (colors.length < 2) {
		// Default: Blue to Red
		colors = ["#1d4ed8", "#b91c1c"];
	}

	// Normalize value to 0..1 range is handled by the caller.
	// Find the segment in the gradient
	const segmentCount = colors.length - 1;
	const segmentIndex = Math.min(Math.floor(value * segmentCount), segmentCount - 1);
	const segmentT = (value * segmentCount) - segmentIndex;

	const c1 = hexToRgb(colors[segmentIndex]);
	const c2 = hexToRgb(colors[segmentIndex + 1]);

	const r = Math.round(c1.r + (c2.r - c1.r) * segmentT);
	const g = Math.round(c1.g + (c2.g - c1.g) * segmentT);
	const b = Math.round(c1.b + (c2.b - c1.b) * segmentT);

	return `rgb(${r},${g},${b})`;
}

function hexToRgb(hex: string) {
	hex = hex.replace("#", "");
	if (hex.length === 3) {
		hex = hex.split("").map(s => s + s).join("");
	}
	const num = parseInt(hex, 16);
	return {
		r: (num >> 16) & 255,
		g: (num >> 8) & 255,
		b: num & 255
	};
}

export const processDataset = (chartType: string) => (dataset: Dataset) => {
	if (!dataset) return { label: "unknown", data: [] };

	let data: any;
	const isHeatmap = chartType === "heatmap" || chartType === "multi-heatmap";

	if (isHeatmap && isHeatmapDataset(dataset)) {
		// Use pointData for heatmaps
		data = dataset.pointData.map((p) => ({
			x: p.x,
			y: p.y,
			v: p.v,
		}));

		const valuesOnly = data.map((d: any) => d.v);
		const min = Math.min(...valuesOnly);
		const max = Math.max(...valuesOnly);
		const range = max - min || 1;

		const colors = (dataset.backgroundColor && dataset.backgroundColor.length > 0)
			? dataset.backgroundColor
			: ["#1d4ed8", "#b91c1c"];

		// Calculate unique coordinates to determine grid size
		const uniqueX = new Set(data.map((p: any) => p.x)).size;
		const uniqueY = new Set(data.map((p: any) => p.y)).size;

		return {
			type: "matrix",
			label: dataset.label || "Heatmap",
			data: data,
			backgroundColorList: colors, // Store for the custom legend plugin
			width: ({ chart }: any) => {
				const area = chart.chartArea;
				if (!area) return 1;
				return (area.right - area.left) / (uniqueX || 1);
			},
			height: ({ chart }: any) => {
				const area = chart.chartArea;
				if (!area) return 1;
				return (area.bottom - area.top) / (uniqueY || 1);
			},
			backgroundColor: (ctx: any) => {
				const val = ctx.dataset.data[ctx.dataIndex]?.v ?? 0;
				let normalized = (val - min) / range;

				// Apply Power Scaling (Gamma) for extreme sensitivity to small changes
				// A power < 1 stretches the "cool" end of the spectrum
				normalized = Math.pow(normalized, 0.3);

				return interpolateColor(normalized, colors);
			},
			borderColor: "rgba(255,255,255,0.1)",
			borderWidth: 1,
			datalabels: { display: false }
		};
	}

	// For standard datasets (Grid or Categorical)
	if ("data" in dataset) {
		data = dataset.data;
	} else {
		console.warn(`Empty data in dataset ${dataset.label}`);
		data = [];
	}

	const datalabels = getDataLabels(dataset.dataLabels, chartType as any);

	return {
		label: dataset.label || "Unnamed dataset",
		data: data,
		borderColor: dataset.borderColor || "#000000",
		backgroundColor: (dataset as any).backgroundColor ?? dataset.borderColor ?? "#000000",
		tension: (dataset as any).tension ?? 0,
		fill: (dataset as any).fill ?? false,
		hidden: dataset.hidden ?? false,
		pointRadius: (dataset as any).pointRadius ?? 0,
		borderWidth: dataset.borderWidth ?? 2,
		showLine: !(dataset as any).hideLine !== false,
		togglable: dataset.togglable !== false,
		pointStyle: (dataset as any).pointStyle ?? undefined,
		datalabels: datalabels,
	};
}

export function newScales(chartConfig: charting.Chart, hasContinuousAxes: boolean) {
	const xAxisType = chartConfig.xAxisConfig || (hasContinuousAxes ? "linear" : "category");
	const yAxisType = chartConfig.yAxisConfig || (hasContinuousAxes ? "linear" : "linear");

	return {
		x: {
			type: xAxisType as any,
			border: {
				display: !hasContinuousAxes,
			},
			title: {
				display: !!chartConfig.xAxisLabel,
				text: chartConfig.xAxisLabel ?? "",
				color: "#000000",
				font: {
					size: 14,
					weight: "bold",
				},
			},
			ticks: {
				color: "#000000",
				font: {
					size: 12,
				},
			},
			grid: {
				color: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0
						? "#000000"
						: "rgba(0,0,0,0.1)",
				lineWidth: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0 ? 2 : 1,
			},
		},
		y: {
			type: yAxisType as any,
			border: {
				display: !hasContinuousAxes,
			},
			title: {
				display: !!chartConfig.yAxisLabel,
				text: chartConfig.yAxisLabel ?? "",
				color: "#000000",
				font: {
					size: 14,
					weight: "bold",
				},
			},
			ticks: {
				color: "#000000",
				font: {
					size: 12,
				},
			},
			grid: {
				color: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0
						? "#000000"
						: "rgba(0,0,0,0.1)",
				lineWidth: (ctx: any) =>
					hasContinuousAxes && ctx.tick?.value === 0 ? 2 : 1,
			},
		},
	};
}
