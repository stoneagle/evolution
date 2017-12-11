import { PanelConfig } from './panel-config';
import * as _ from 'lodash';

type KeyValue = [string, number];

export class Mapper {

  private _panelConfig: PanelConfig;

  constructor(panelConfig: PanelConfig) {
    this._panelConfig = panelConfig;
  }

  _mapKeysTotal(seriesList): KeyValue[] {
    if(seriesList.length !== 1) {
      throw new Error('Expecting list of keys: got more than one timeseries');
    }
    var kv = {};
    var datapointsLength = seriesList[0].datapoints.length;
    for(let i = 0; i < datapointsLength; i++) {
      let k = seriesList[0].datapoints[i][0].toString();
      if(kv[k] === undefined) {
        kv[k] = 0;
      }
      kv[k]++;
    }

    var res: KeyValue[] = [];
    for(let k in kv) {
      res.push([k, kv[k]]);
    }

    return res;

  }

  _mapNumeric(seriesList): KeyValue[] {
    if(seriesList.length != 2) {
      throw new Error('Expecting timeseries in format (key, value). You can use keys only in total mode');
    }
    if(seriesList[0].datapoints.length !==  seriesList[1].datapoints.length) {
      throw new Error('Timeseries has different length');
    }

    var kv = {};
    var datapointsLength = seriesList[0].datapoints.length;

    var nullMapping = this._panelConfig.getValue('nullMapping');

    for(let i = 0; i < datapointsLength; i++) {
      let k = seriesList[0].datapoints[i][0].toString();
      let v = seriesList[1].datapoints[i][0];
      let vn = parseFloat(v);
      if(v === null) {
        if(nullMapping === undefined || nullMapping === null) {
          throw new Error('Got null value. You set null value mapping in Options -> Mapping -> Null');
        }
        console.log('nullMapping ->' + nullMapping);
        vn = nullMapping;
      }
      if(isNaN(vn)) {
        throw new Error('Got non-numberic value: ' + v);
      }
      if(kv[k] === undefined) {
        kv[k] = [];
      }
      kv[k].push(vn);
    }

    var res: KeyValue[] = [];
    for(let k in kv) {
      res.push([k, this._flatSeries(kv[k])]);
    }

    return res;

  }

  _flatSeries(values: number[]): number {

    if(values.length === 0) {
      return 0;
    }

    var t = this._panelConfig.getValue('statNameOptionValue');

    if(t === 'total') {
      return _.sum(values);
    }

    if(t === 'max') {
      return _.max(values) as number;
    }

    if(t === 'min') {
      return _.min(values) as number;
    }

    if(t === 'current') {
      return values[values.length - 1];
    }

    return 0;
  }

}
