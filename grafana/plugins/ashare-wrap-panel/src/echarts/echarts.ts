import { PanelConfig } from '../panel-config';
import { coreModule } from 'grafana/app/core/core';

var directiveInited = false;
export function EchartsInit(panelConfig: PanelConfig, directiveName: string = "echarts") {
  if(directiveInited) {
    return;
  }
  directiveInited = true;

  coreModule.directive(directiveName, function() {
    return {
      templateUrl: panelConfig.pluginDirName + 'echarts/echarts.html',
      restrict: 'E',
      scope: {
        item: "="
      }
    };
  });
}
