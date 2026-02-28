import { common } from "../wailsjs/go/models";
import { buildInputFieldName } from "./lab-init";

class Registry {
  registry: Record<string, common.LabMetadata>;
  constructor() {
    this.registry = {};
  }

  addLab(id: string, metadata: common.LabMetadata) {
    this.registry[id] = metadata;
  }

  getLab(id: string): common.LabMetadata | undefined {
    return this.registry[id];
  }

  getChartVariables(labId: string, chartId: string): Record<string, number> {
    const vars: Record<string, number> = {};
    const lab = this.getLab(labId);
    if (!lab) return vars;
    const chart = lab.Charts[chartId];
    if (!chart) return vars;

    const variableDefs = chart.chartVariables;
    if (!variableDefs) {
      return vars;
    }

    variableDefs.forEach((variableField) => {
      const input = document.getElementById(
        `${chartId}-${variableField.id}`,
      ) as HTMLInputElement;
      if (!input) {
        return;
      }
      vars[buildInputFieldName(chartId, null, variableField.id)] = parseFloat(
        input.value,
      );
    });

    return vars;
  }

  getGraphVariables(
    labId: string,
    chartId: string,
  ): Record<string, Record<string, number>> {
    console.log("getGraphVariables");
    const vars: Record<string, Record<string, number>> = {
      [chartId]: {},
    };
    const lab = this.getLab(labId);
    if (!lab) return vars;
    const chart = lab.Charts[chartId];
    if (!chart) return vars;

    for (const graphId in chart.graphVariables) {
      if (!chart.graphVariables[graphId]) {
        continue;
      }

      for (const field of chart.graphVariables[graphId]) {
        const input = document.getElementById(
          buildInputFieldName(chartId, graphId, field.id),
        ) as HTMLInputElement;
        if (!input) {
          continue;
        }
        vars[chartId][buildInputFieldName(graphId, null, field.id)] =
          parseFloat(input.value);
      }
    }

    console.log("Graph Variables:", vars);

    return vars;
  }
}

export const registry = new Registry();
