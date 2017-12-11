import { PanelConfig } from '../panel-config';
import { coreModule } from 'grafana/app/core/core';

var directiveInited = false;
export function initWrap(panelConfig: PanelConfig, directiveName: string = "wrap") {
  if(directiveInited) {
    return;
  }
  directiveInited = true;

  coreModule.directive(directiveName, function() {
    return {
      templateUrl: panelConfig.pluginDirName + 'wrap/wrap.html',
      restrict: 'E',
      scope: {
        item: "="
      }
    };
  });
}
