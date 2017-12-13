import { PanelConfig } from './panel-config';
import * as _ from 'lodash';

type rawData = any[];
const shareTable: string = "demo";

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
        case "close":
        case "high":
        case "low":
        case "open":
        case "volume":
          seriesMap.set(columnName, seriesList[i].datapoints)
      }
    }

    let closeList = seriesMap.get("close")
    let highList = seriesMap.get("high")
    let lowList = seriesMap.get("low")
    let openList = seriesMap.get("open")
    let volumeList = seriesMap.get("volume")
    console.log(volumeList);
    let ret = new Array();
    for (let i = 0; i < seriesLength; i++) {
      ret[i] = new Array();
      let date = this._getDateStr(closeList![i][1]);
      ret[i][0] = date;
      ret[i][1] = openList![i][0].toFixed(2);
      ret[i][2] = closeList![i][0].toFixed(2);
      ret[i][3] = lowList![i][0].toFixed(2);
      ret[i][4] = highList![i][0].toFixed(2);
      ret[i][5] = volumeList![i][0].toFixed(2);
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
