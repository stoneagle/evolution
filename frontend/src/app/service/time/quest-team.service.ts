import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestTeam }                from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestTeamService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.QuestTeam;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/quest-team';
  }

  List(questTeam: QuestTeam): Observable<QuestTeam[]> {
    return this.BaseList<QuestTeam>(questTeam, QuestTeam, this.uri + `/list`).pipe(map(questTeams => {
      return questTeams;
    }))
  }

  Get(id: number): Observable<QuestTeam> {
    return this.BaseGet<QuestTeam>(QuestTeam, this.uri + `/get/${id}`).pipe(map(questTeam => {
      return questTeam;
    }))
  }

  Add(questTeam: QuestTeam): Observable<QuestTeam> {
    return this.BaseAdd<QuestTeam>(questTeam, QuestTeam, this.uri).pipe(map(questTeam => {
      return questTeam;
    }))
  }

  Update(questTeam: QuestTeam): Observable<QuestTeam> {
    return this.BaseUpdate<QuestTeam>(questTeam, QuestTeam, this.uri + `/${questTeam.Id}`).pipe(map(questTeam => {
      return questTeam;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<QuestTeam>(QuestTeam, this.uri + `/${id}`).pipe(map(questTeam => {
      return questTeam;
    }))
  }
}
