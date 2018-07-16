import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';
import { Action }                   from '../../model/time/action';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ActionService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/action';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Action;
  }

  List(): Observable<Action[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Action[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Action(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(action: Action): Observable<Action[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(action)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Action[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Action(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Action> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Action(res.data);
        } else {
          return new Action();
        }
      }),
    )
  }

  Add(action: Action): Observable<Action> {
    this.operation = this.shareSettings.System.Process.Create;
    // action.StartDate.toJSON = this.DateJsonKeepFormat;
    // action.EndDate.toJSON = this.DateJsonKeepFormat;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(action)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Action(res.data);
        } else {
          return new Action();
        }
      }),
    );
  }

  Update(action: Action): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${action.Id}`, JSON.stringify(action)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
    );
  }

  Delete(id: number): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Delete;
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
