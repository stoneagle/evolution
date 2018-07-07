import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders, HttpRequest  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { CookieService }            from 'ngx-cookie';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { SessionUser }              from '../../model/base/sign';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class SignService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.system + '/sign';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService ,
    private cookieService: CookieService,
  ) { 
    super(http, messageHandlerService);
    this.resource = 'SYSTEM.RESOURCE.USER.CONCEPT';
  }

  currentUser: SessionUser = new SessionUser();
  username: string;
  password: string;

  clear(): void {
    this.cookieService.remove(AppConfig.settings.apiServer.auth.session);
    this.currentUser = new SessionUser();
  }

  getCurrentUser(): SessionUser {
    return this.currentUser;
  }

  setAuthHeader(headersMap: any): any {
    switch (AppConfig.settings.apiServer.auth.type) {
      case "basic-auth":
        headersMap["Authorization"] = "Basic " + btoa(`${this.username}:${this.password}`);
        break;
    }
    return headersMap;
  }

  login(username: string, password: string): Observable<Boolean> {
    this.operation = "SYSTEM.PROCESS.SIGNIN";
    this.username = username;
    this.password = password;
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/login`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()), 
      map(res => {
        if (res == undefined) {
          return false;
        }
        if (res.code != 0) {
          return false;
        } else {
          this.currentUser = new SessionUser(res.data);
          return true;
          // 目前由服务端控制
          // let expires: number = 10 * 3600 * 24 * 1000;
          // let date = new Date(Date.now() + expires);
          // let cookieptions: CookieOptions = {
          //     path: "/",
          //     expires: date
          // };
          // this.cookieService.put(AppConfig.settings.apiServer.auth.session, res.data, cookieptions);
        }
      }),
    );
  }

  logout(): Observable<Boolean> {
    this.operation = "SYSTEM.PROCESS.SIGNOUT";
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/logout`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
      map(res => {
        if (res.code != 0) {
          return false;
        } else {
          this.clear();
          return true;
        }
      }),
    );
  }

  current(): Observable<SessionUser> {
    this.operation = "SYSTEM.PROCESS.GET";
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/current`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:SessionUser
        if (res.code != 0) {
          ret = new SessionUser();
        } else {
          ret = new SessionUser(res.data);
        }
        this.currentUser = ret;
        return ret;
      }),
    );
  }
}
