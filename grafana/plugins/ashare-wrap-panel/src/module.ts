import { PanelConfig } from './panel-config';
import { MetricsPanelCtrl, loadPluginCss } from 'grafana/app/plugins/sdk';
import { Mapper } from './mapper';
import { initWrap } from './wrap/wrap';
import { WrapUtil } from './wrap/util';
import * as _ from 'lodash';
import * as echarts from 'echarts';

const defaults = {
  coloringType: 'none',
  sortingOrder: 'none',
  valueLabelType: 'percentage',
  prefix: '',
  postfix: '',
  thresholds: '10, 30',
  // https://github.com/grafana/grafana/blob/v4.1.1/public/app/plugins/panel/singlestat/module.ts#L57
  colors: ["rgba(245, 54, 54, 0.9)", "rgba(237, 129, 40, 0.89)", "rgba(50, 172, 45, 0.97)"],
  colorsKeyMappingDefault: "rgba(245, 54, 54, 0.9)",
  colorKeyMappings: [],
  nullMapping: undefined,
};

export class PanelCtrl extends MetricsPanelCtrl {
  static templateUrl = "partials/template.html";

  public wraper: WrapUtil;
  public mapper: Mapper;
  private $panelContainer: any;
  private $panelContoller: any;
  private minHeight: number;
  private _panelConfig: PanelConfig;

  private echartsInitFlag = false;
  private wrapChart;

  constructor($scope: any, $injector) {
    super($scope, $injector);

    _.defaults(this.panel, defaults);

    this._panelConfig = new PanelConfig(this.panel);
    this.initStyles();

    initWrap(this._panelConfig, 'asharePluginWrap');
    this.wraper = new WrapUtil();
    this.mapper = new Mapper(this._panelConfig);

		this.events.on('render', this.onRender.bind(this));
		// this.events.on('refresh', this.postRefresh.bind(this));
		// this.events.on('data-error', this.onDataError.bind(this));
		// this.events.on('data-snapshot-load', this.onDataReceived.bind(this));
		// this.events.on('init-edit-mode', this.onInitEditMode.bind(this));
    this.events.on('data-received', this.onDataReceived.bind(this));
  }

  link(scope, element, attrs, ctrl) {
		this.$panelContainer = element.find('.panel-container');
		this.$panelContoller = ctrl;
  }

  initStyles() {
    // 读取grafana基础样式
    loadPluginCss({
      light: this._panelConfig.pluginDirName + 'css/panel.base.css',
      dark: this._panelConfig.pluginDirName + 'css/panel.base.css'
    });
    loadPluginCss({
      light: this._panelConfig.pluginDirName + 'css/panel.light.css',
      dark: this._panelConfig.pluginDirName + 'css/panel.dark.css'
    });
  }

  onRender() {
    this.setElementHeight();
  }

	setElementHeight() {
    this.$panelContainer.find('.ashare-panel').css('min-height', this.$panelContoller.height + 'px');
    this.minHeight = this.$panelContoller.height-10;
    this.$panelContainer.find('.wrap-plugin').css('min-height', this.minHeight + 'px');
    this.$panelContainer.find('.wrap').css('min-height', this.minHeight + 'px');
	}


  onDataReceived(seriesList: any) {
    // influxdb数据更新时,绘制图表
    this.onRender();
    if (this.echartsInitFlag == false) {
      var container = <HTMLDivElement> document.getElementById("wrap");
      this.wrapChart = echarts.init(container);
      this.echartsInitFlag = true
    }
    var rawData: any[] = this.mapper.mapShare(seriesList);
    this.wraper.setData(rawData);
    this.wrapChart.setOption(this.wraper.getOption());
  }

  onInitEditMode() {
    // 编辑模式初始化
    var thisPartialPath = this._panelConfig.pluginDirName + 'partials/';
    this.addEditorTab('Options', thisPartialPath + 'options.html', 2);
  }

  _dataError(err) {
    this.$scope.data = [];
    this.$scope.dataError = err;
  }
}
