import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { Resp }                     from '../../model/base/resp';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';

@Injectable()
export class ConfigService extends BaseService {
  private uri = '/config';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) { 
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Quant.Resource.Config;
  }

  AssetList(): Observable<Map<string, string>> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + "/asset").pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Map<string, string> = new Map();
        if (res && res.code == 0) {
          for (let key in res.data) {
            ret.set(key, res.data[key]);
          }
        }
        return ret;
      }),
    )
  }

  TypeList(resource: string): Observable<Map<string, Map<string, string[]>>> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/type/${resource}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Map<string, Map<string, string[]>> = new Map();
        if (res && res.code == 0) {
          for (let key in res.data) {
            let typeMap: Map<string, string[]> = new Map();
            for (let one in res.data[key]) {
              typeMap.set(one, res.data[key][one])
            }
            ret.set(key, typeMap);
          }
        }
        return ret;
      }),
    )
  }

  StrategyList(ctype: string): Observable<Map<string, string>> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/strategy/${ctype}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Map<string, string> = new Map();
        if (res && res.code == 0) {
          for (let key in res.data) {
            ret.set(res.data[key], res.data[key]);
          }
        }
        return ret;
      }),
    )
  }
}
