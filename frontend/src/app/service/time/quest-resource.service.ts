import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';
import { QuestResource }            from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestResourceService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/quest-resource';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.QuestResource;
  }

  List(): Observable<QuestResource[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:QuestResource[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestResource(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(questResource: QuestResource): Observable<QuestResource[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(questResource)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:QuestResource[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestResource(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<QuestResource> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestResource(res.data);
        } else {
          return new QuestResource();
        }
      }),
    )
  }

  Add(questResource: QuestResource): Observable<QuestResource> {
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(questResource)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestResource(res.data);
        } else {
          return new QuestResource();
        }
      }),
    );
  }

  Update(questResource: QuestResource): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${questResource.Id}`, JSON.stringify(questResource)).pipe(
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
