import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestEntity }              from '../../model/time/quest';
import { Resp }                 from '../../model/base/resp';

@Injectable()
export class QuestEntityService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/quest-entity';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.QUEST-ENTITY.CONCEPT';
  }

  List(): Observable<QuestEntity[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:QuestEntity[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestEntity(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(questEntity: QuestEntity): Observable<QuestEntity[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(questEntity)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:QuestEntity[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestEntity(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<QuestEntity> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestEntity(res.data);
        } else {
          return new QuestEntity();
        }
      }),
    )
  }

  Add(questEntity: QuestEntity): Observable<QuestEntity> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(questEntity)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestEntity(res.data);
        } else {
          return new QuestEntity();
        }
      }),
    );
  }

  Update(questEntity: QuestEntity): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${questEntity.Id}`, JSON.stringify(questEntity)).pipe(
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
