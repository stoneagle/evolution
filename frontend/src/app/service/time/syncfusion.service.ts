import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { ShareSettings }            from '../../shared/settings';
import { SignService  }             from '../system/sign.service';
import { Kanban, Gantt }            from '../../model/time/syncfusion';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class SyncfusionService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/syncfusion';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected signService: SignService,
    protected shareSettings: ShareSettings,
  ) {
    super(http, messageHandlerService);
  }

  GetScheduleManager(): any {
    return new ej.DataManager({
      url: AppConfig.settings.apiServer.endpoint + this.uri + `/list/schedule`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  GetTreeGridManager(fieldId: number): any {
    return new ej.DataManager({
      url: AppConfig.settings.apiServer.endpoint + this.uri + `/list/treegrid/${fieldId}`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  ListKanban(): Observable<Kanban[]> {
    this.resource = this.shareSettings.Time.Resource.Task;
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/kanban`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Kanban[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Kanban(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListGantt(): Observable<Gantt[]> {
    this.resource = this.shareSettings.Time.Resource.Quest;
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/gantt`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Gantt[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Gantt(one));
            }
          )
        }
        return ret; 
      }),
    )
  }
}
