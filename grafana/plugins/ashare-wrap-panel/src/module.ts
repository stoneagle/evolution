import { PanelConfig } from './panel-config';
import { MetricsPanelCtrl, loadPluginCss } from 'grafana/app/plugins/sdk';
import { Mapper } from './mapper';
import { EchartsInit } from './echarts/echarts';
import { EchartsUtil } from './echarts/util';
import * as _ from 'lodash';
import * as echarts from 'echarts';

const defaults = {
  echartsType: 'share',
  // https://github.com/grafana/grafana/blob/v4.1.1/public/app/plugins/panel/singlestat/module.ts#L57
  nullMapping: undefined,
};

export class PanelCtrl extends MetricsPanelCtrl {
  static templateUrl = "partials/template.html";

  public util: EchartsUtil;
  public mapper: Mapper;
  private echartsTypeOptions = ["share", "basic"];
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

    EchartsInit(this._panelConfig, 'asharePluginEcharts');
    this.util = new EchartsUtil();
    this.mapper = new Mapper(this._panelConfig);

		this.events.on('render', this.onRender.bind(this));
		this.events.on('data-error', this.onDataError.bind(this));
		this.events.on('init-edit-mode', this.onInitEditMode.bind(this));
    this.events.on('data-received', this.onDataReceived.bind(this));
		this.events.on('refresh', this.onRefresh.bind(this));
		// this.events.on('data-snapshot-load', this.onDataReceived.bind(this));
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
    // 渲染样式
    this._setElementHeight();
  }

  onRefresh() {
  }

	_setElementHeight() {
    // 调整容器高度
    this.$panelContainer.find('.ashare-panel').css('min-height', this.$panelContoller.height + 'px');
    this.minHeight = this.$panelContoller.height-10;
    this.$panelContainer.find('.echarts-plugin').css('min-height', this.minHeight + 'px');
    this.$panelContainer.find('.echarts').css('min-height', this.minHeight + 'px');
	}

  onDataReceived(seriesList: any) {
    // influxdb数据更新时,绘制图表
    this.onRender();
    if (this._panelConfig.getValue('echartsType') === 'share') {
      // 展示k线类别数据 
      if (this.echartsInitFlag == false) {
        var container = <HTMLDivElement> document.getElementById("echarts");
        this.wrapChart = echarts.init(container);
        this.echartsInitFlag = true
      }
      var rawData: any[] = this.mapper.mapShare(seriesList);
      this.util.setData(rawData);
      this.wrapChart.setOption(this.util.getKLineOption());
    } 
  }

  onInitEditMode() {
    // 编辑模式初始化
    var thisPartialPath = this._panelConfig.pluginDirName + 'partials/';
    this.addEditorTab('Options', thisPartialPath + 'options.html', 2);
  }

  onDataError(err) {
    // 错误异常时的处理
    this.$scope.data = [];
    this.$scope.dataError = err;
  }
}
