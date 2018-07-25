import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { UserResource }             from '../../model/time/user-resource';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class UserResourceService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.UserResource;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/user-resource';
  }

  List(userResource: UserResource): Observable<UserResource[]> {
    return this.BaseList<UserResource>(userResource, UserResource, this.uri + `/list`).pipe(map(userResources => {
      return userResources;
    }))
  }
  
  AreaSumTime(userResource: UserResource): Observable<number> {
    return this.BasePostResp(userResource, this.uri + `/list/time`).pipe(map(resp => {
      return resp.data;
    }))
  }

  Get(id: number): Observable<UserResource> {
    return this.BaseGet<UserResource>(UserResource, this.uri + `/get/${id}`).pipe(map(userResource => {
      return userResource;
    }))
  }
  
  Add(userResource: UserResource): Observable<UserResource> {
    return this.BaseAdd<UserResource>(userResource, UserResource, this.uri).pipe(map(userResource => {
      return userResource;
    }))
  }
  
  Update(userResource: UserResource): Observable<UserResource> {
    return this.BaseUpdate<UserResource>(userResource, UserResource, this.uri + `/${userResource.Id}`).pipe(map(userResource => {
      return userResource;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<UserResource>(UserResource, this.uri + `/${id}`).pipe(map(userResource => {
      return userResource;
    }))
  }
}
