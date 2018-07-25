import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestResource }            from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestResourceService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.QuestResource;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/quest-resource';
  }

  List(questResource: QuestResource): Observable<QuestResource[]> {
    return this.BaseList<QuestResource>(questResource, QuestResource, this.uri + `/list`).pipe(map(questResources => {
      return questResources;
    }))
  }

  Get(id: number): Observable<QuestResource> {
    return this.BaseGet<QuestResource>(QuestResource, this.uri + `/get/${id}`).pipe(map(questResource => {
      return questResource;
    }))
  }

  Add(questResource: QuestResource): Observable<QuestResource> {
    return this.BaseAdd<QuestResource>(questResource, QuestResource, this.uri).pipe(map(questResource => {
      return questResource;
    }))
  }

  Update(questResource: QuestResource): Observable<QuestResource> {
    return this.BaseUpdate<QuestResource>(questResource, QuestResource, this.uri + `/${questResource.Id}`).pipe(map(questResource => {
      return questResource;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<QuestResource>(QuestResource, this.uri + `/${id}`).pipe(map(questResource => {
      return questResource;
    }))
  }
}
