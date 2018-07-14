import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Action }                   from '../../model/time/action';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ActionService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/action';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.ACTION.CONCEPT';
  }

  List(): Observable<Action[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
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
    this.operation = 'SYSTEM.PROCESS.LIST';
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
    this.operation = 'SYSTEM.PROCESS.GET';
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
    this.operation = 'SYSTEM.PROCESS.CREATE';
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
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${action.Id}`, JSON.stringify(action)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
    );
  }

  Delete(id: number): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
