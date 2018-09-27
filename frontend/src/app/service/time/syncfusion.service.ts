import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { SignService  }             from '../system/sign.service';
import { Kanban, Gantt }            from '../../model/time/syncfusion';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class SyncfusionService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected signService: SignService,
  ) {
    super(http, messageHandlerService);
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/syncfusion';
  }

  GetScheduleManager(): any {
    return new ej.DataManager({
      url: this.uri + `/list/schedule`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  GetTreeGridManager(fieldId: number): any {
    return new ej.DataManager({
      url: this.uri + `/list/treegrid/${fieldId}`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  GetGanttManager(level: string, status: string): any {
    return new ej.DataManager({
      url: this.uri + `/list/gantt/${level}/${status}`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  GetKanbanManager(): any {
    return new ej.DataManager({
      url: this.uri + `/list/kanban`,
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
  }

  ListKanban(questId: number, projectId: number): Observable<Kanban[]> {
    this.resource = this.shareSettings.Time.Resource.Task;
    this.operation = this.shareSettings.System.Process.List;
    return this.BaseResp(this.uri + `/list/kanban/${questId}/${projectId}`).pipe(
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

  ListGantt(level: string, status: string): Observable<Gantt[]> {
    this.resource = this.shareSettings.Time.Resource.Quest;
    this.operation = this.shareSettings.System.Process.List;
    return this.BaseResp(this.uri + `/list/gantt/${level}/${status}`).pipe(
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
