import { PanelConfig } from './panel-config';
import * as _ from 'lodash';

type rawData = any[];
const shareTable: string = "demo";
const index_open: string = "open";
const index_close: string = "close";
const index_high: string = "high";
const index_low: string = "low";
const index_volume: string = "volume";
const index_dif: string = "dif";
const index_dea: string = "dea";
const index_macd: string = "macd";

export class Mapper {
  private _panelConfig: PanelConfig;

  constructor(panelConfig: PanelConfig) {
    this._panelConfig = panelConfig;
  }

  mapShare(seriesList): any[] {
    if(seriesList.length == 0) {
      throw new Error('Expecting list of keys: got less than one timeseries');
    }
    let seriesMap = new Map<string, any[]>();
    let seriesLength = seriesList[0].datapoints.length;
    for (let i = 0; i < seriesList.length; i++) {
      let targetArr = seriesList[i].target.split(".")
      let columnName = targetArr[1]
      switch (columnName) {
        case index_open:
        case index_close:
        case index_low:
        case index_high:
        case index_volume:
        case index_dif:
        case index_dea:
        case index_macd:
          seriesMap.set(columnName, seriesList[i].datapoints)
      }
    }

    let ret = new Array();
    for (let i = 0; i < seriesLength; i++) {
      ret[i] = new Array();
      let date = this._getDateStr(seriesMap.get(index_close)![i][1]);
      ret[i][0] = date;
      ret[i][1] = seriesMap.get(index_open)![i][0].toFixed(2);
      ret[i][2] = seriesMap.get(index_close)![i][0].toFixed(2);
      ret[i][3] = seriesMap.get(index_low)![i][0].toFixed(2);
      ret[i][4] = seriesMap.get(index_high)![i][0].toFixed(2);
      ret[i][5] = seriesMap.get(index_volume)![i][0].toFixed(2);
      ret[i][6] = seriesMap.get(index_dif)![i][0].toFixed(2);
      ret[i][7] = seriesMap.get(index_dea)![i][0].toFixed(2);
      ret[i][8] = seriesMap.get(index_macd)![i][0].toFixed(2);
    }
    return ret;
  }

	_getDateStr(nS): string {
    let date = new Date(parseInt(nS));
    let year    = date.getFullYear();
    let month   = date.getMonth();
    let day     = date.getDay();
    let hour    = date.getHours();
    let minute  = date.getMinutes();
    let seconds = date.getSeconds();  
    let dateStr = year + '-' + month + '-' + day + ' ' + hour + ':' + minute;
    return dateStr
	}
}
