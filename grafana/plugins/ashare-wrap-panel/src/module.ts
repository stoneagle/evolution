import { PanelConfig } from './panel-config';
import { MetricsPanelCtrl, loadPluginCss } from 'grafana/app/plugins/sdk';
import { initWrap } from './wrap/wrap';
import { WrapUtil } from './wrap/util';
import * as _ from 'lodash';
import * as echarts from 'echarts';

const defaults = {
  statNameOptionValue: 'current',
  statProgressType: 'shared',
  statProgressMaxValue: null,
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
    this._initStyles();

    initWrap(this._panelConfig, 'asharePluginWrap');
    this.wraper = new WrapUtil();

		this.events.on('render', this.onRender.bind(this));
		// this.events.on('refresh', this.postRefresh.bind(this));
		// this.events.on('data-error', this.onDataError.bind(this));
		// this.events.on('data-received', this.onDataReceived.bind(this));
		// this.events.on('data-snapshot-load', this.onDataReceived.bind(this));
		// this.events.on('init-edit-mode', this.onInitEditMode.bind(this));
    // this.events.on('data-received', this._onDataReceived.bind(this));
  }

  link(scope, element, attrs, ctrl) {
		this.$panelContainer = element.find('.panel-container');
		this.$panelContoller = ctrl;
  }

  _initStyles() {
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
    this.$panelContainer.find('.wrap').css('min-height', this.minHeight + 'px');
	}


  _onDataReceived(seriesList: any) {
    // influxdb数据更新时,绘制图表
    if (this.echartsInitFlag == false) {
      var container = <HTMLDivElement> document.getElementById("wrap");
      this.wrapChart = echarts.init(container);
      this.echartsInitFlag = true
    }
    var rawData: any[] = [
        ['2013/5/8', 2242.39,2246.3,2235.42,2255.21],
        ['2013/5/9', 2246.96,2232.97,2221.38,2247.86],
        ['2013/5/10', 2228.82,2246.83,2225.81,2247.67],
        ['2013/5/13', 2247.68,2241.92,2231.36,2250.85],
        ['2013/5/14', 2238.9,2217.01,2205.87,2239.93],
        ['2013/5/15', 2217.09,2224.8,2213.58,2225.19],
        ['2013/5/16', 2221.34,2251.81,2210.77,2252.87],
        ['2013/5/17', 2249.81,2282.87,2248.41,2288.09],
        ['2013/5/20', 2286.33,2299.99,2281.9,2309.39],
        ['2013/5/21', 2297.11,2305.11,2290.12,2305.3],
        ['2013/5/22', 2303.75,2302.4,2292.43,2314.18],
        ['2013/5/23', 2293.81,2275.67,2274.1,2304.95],
        ['2013/5/24', 2281.45,2288.53,2270.25,2292.59],
        ['2013/5/27', 2286.66,2293.08,2283.94,2301.7],
        ['2013/5/28', 2293.4,2321.32,2281.47,2322.1],
        ['2013/5/29', 2323.54,2324.02,2321.17,2334.33],
        ['2013/5/30', 2316.25,2317.75,2310.49,2325.72],
        ['2013/5/31', 2320.74,2300.59,2299.37,2325.53],
        ['2013/6/3', 2300.21,2299.25,2294.11,2313.43],
        ['2013/6/4', 2297.1,2272.42,2264.76,2297.1],
        ['2013/6/5', 2270.71,2270.93,2260.87,2276.86],
        ['2013/6/6', 2264.43,2242.11,2240.07,2266.69],
        ['2013/6/7', 2242.26,2210.9,2205.07,2250.63],
        ['2013/6/13', 2190.1,2148.35,2126.22,2190.1]
    ];
    this.wraper.setData(rawData);
		this.wrapChart.setOption(this.wraper.getOption());
    this.onRender();
  }

  _onInitEditMode() {
    // 编辑模式初始化
    var thisPartialPath = this._panelConfig.pluginDirName + 'partials/';
    this.addEditorTab('Options', thisPartialPath + 'options.html', 2);
  }

  invertColorOrder() {
    var tmp = this.panel.colors[0];
    this.panel.colors[0] = this.panel.colors[2];
    this.panel.colors[2] = tmp;
    this.onRender();
  }

  addColorKeyMapping() {
    this.panel.colorKeyMappings.push({
      key: 'KEY_NAME',
      color: "rgba(50, 172, 45, 0.97)"
    });
  }

  removeColorKeyMapping(index) {
    this.panel.colorKeyMappings.splice(index, 1);
    this.onRender();
  }

  _dataError(err) {
    this.$scope.data = [];
    this.$scope.dataError = err;
  }
}
