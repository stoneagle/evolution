import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Action }                   from '../../model/time/action';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ActionService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Action;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/action';
  }

  List(): Observable<Action[]> {
    return this.BaseList<Action>(Action, this.uri + `/list`).pipe(map(actions => {
      return actions;
    }))
  }

  ListWithCondition(action: Action): Observable<Action[]> {
    return this.BaseListWithCondition<Action>(action, Action, this.uri + `/list`).pipe(map(actions => {
      return actions;
    }))
  }

  Get(id: number): Observable<Action> {
    return this.BaseGet<Action>(Action, this.uri + `/get/${id}`).pipe(map(action => {
      return action;
    }))
  }

  Add(action: Action): Observable<Action> {
    return this.BaseAdd<Action>(action, Action, this.uri).pipe(map(action => {
      return action;
    }))
  }

  Update(action: Action): Observable<Action> {
    return this.BaseUpdate<Action>(action, Action, this.uri + `/${action.Id}`).pipe(map(action => {
      return action;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Action>(Action, this.uri + `/${id}`).pipe(map(action => {
      return action;
    }))
  }
}
