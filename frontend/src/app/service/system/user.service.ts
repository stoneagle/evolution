import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { User }                     from '../../model/system/user';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class UserService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.System.Resource.User;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.system + '/user';
  }

  List(user: User): Observable<User[]> {
    return this.BaseList<User>(user, User, this.uri + `/list`).pipe(map(users => {
      return users;
    }))
  }
  
  Get(id: number): Observable<User> {
    return this.BaseGet<User>(User, this.uri + `/get/${id}`).pipe(map(user => {
      return user;
    }))
  }
  
  Add(user: User): Observable<User> {
    return this.BaseAdd<User>(user, User, this.uri).pipe(map(user => {
      return user;
    }))
  }
  
  Update(user: User): Observable<User> {
    return this.BaseUpdate<User>(user, User, this.uri + `/${user.Id}`).pipe(map(user => {
      return user;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<User>(User, this.uri + `/${id}`).pipe(map(user => {
      return user;
    }))
  }
}
