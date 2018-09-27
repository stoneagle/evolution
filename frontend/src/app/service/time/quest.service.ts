import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Quest }                    from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Quest;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/quest';
  }

  List(quest: Quest): Observable<Quest[]> {
    return this.BaseList<Quest>(quest, Quest, this.uri + `/list`).pipe(map(quests => {
      return quests;
    }))
  }

  Get(id: number): Observable<Quest> {
    return this.BaseGet<Quest>(Quest, this.uri + `/get/${id}`).pipe(map(quest => {
      return quest;
    }))
  }

  Add(quest: Quest): Observable<Quest> {
    return this.BaseAdd<Quest>(quest, Quest, this.uri).pipe(map(quest => {
      return quest;
    }))
  }

  Update(quest: Quest): Observable<Quest> {
    return this.BaseUpdate<Quest>(quest, Quest, this.uri + `/${quest.Id}`).pipe(map(quest => {
      return quest;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Quest>(Quest, this.uri + `/${id}`).pipe(map(quest => {
      return quest;
    }))
  }
}
