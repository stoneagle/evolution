import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestTarget }              from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestTargetService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.QuestTarget;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/quest-target';
  }

  BatchSave(questTargets: QuestTarget[]): Observable<Boolean> {
    return this.BaseAdd<Resp>(questTargets, Resp, this.uri + `/batch`).pipe(map(res => {
      return true;
    }))
  }

  List(questTarget: QuestTarget): Observable<QuestTarget[]> {
    return this.BaseList<QuestTarget>(questTarget, QuestTarget, this.uri + `/list`).pipe(map(questTargets => {
      return questTargets;
    }))
  }

  Get(id: number): Observable<QuestTarget> {
    return this.BaseGet<QuestTarget>(QuestTarget, this.uri + `/get/${id}`).pipe(map(questTarget => {
      return questTarget;
    }))
  }

  Add(questTarget: QuestTarget): Observable<QuestTarget> {
    return this.BaseAdd<QuestTarget>(questTarget, QuestTarget, this.uri).pipe(map(questTarget => {
      return questTarget;
    }))
  }

  Update(questTarget: QuestTarget): Observable<QuestTarget> {
    return this.BaseUpdate<QuestTarget>(questTarget, QuestTarget, this.uri + `/${questTarget.Id}`).pipe(map(questTarget => {
      return questTarget;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<QuestTarget>(QuestTarget, this.uri + `/${id}`).pipe(map(questTarget => {
      return questTarget;
    }))
  }
}
