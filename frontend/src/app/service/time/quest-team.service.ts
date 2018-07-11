import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestTeam }                from '../../model/time/quest';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class QuestTeamService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/quest-team';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.QUEST-TEAM.CONCEPT';
  }

  List(): Observable<QuestTeam[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:QuestTeam[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestTeam(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(questTeam: QuestTeam): Observable<QuestTeam[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(questTeam)).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:QuestTeam[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new QuestTeam(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<QuestTeam> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestTeam(res.data);
        } else {
          return new QuestTeam();
        }
      }),
    )
  }

  Add(questTeam: QuestTeam): Observable<QuestTeam> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(questTeam)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new QuestTeam(res.data);
        } else {
          return new QuestTeam();
        }
      }),
    );
  }

  Update(questTeam: QuestTeam): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${questTeam.Id}`, JSON.stringify(questTeam)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
    );
  }

  Delete(id: number): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>())
    );
  }
}
