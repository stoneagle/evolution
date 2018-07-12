import { Injectable }                            from '@angular/core';
import { HttpClient, HttpHeaders, HttpRequest  } from '@angular/common/http';
import { Observable }                            from 'rxjs';
import { of }                                    from 'rxjs/observable/of';
import { catchError, map, tap  }                 from 'rxjs/operators';
import { CookieService, CookieOptions }          from 'ngx-cookie';
import { AppConfig }                             from '../base/config.service';
import { MessageHandlerService  }                from '../base/message-handler.service';
import { BaseService  }                          from '../base/base.service';
import { SessionUser }                           from '../../model/base/sign';
import { Resp, RespObject }                      from '../../model/base/resp';
import { AuthType }                              from '../../shared/shared.const';

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

  resetUser(): void {
    this.currentUser = new SessionUser();
  }

  getCurrentUser(): SessionUser {
    return this.currentUser;
  }

  getAuthToken(): string {
    let token: string = "";
    switch (AppConfig.settings.apiServer.auth.type) {
      case AuthType.BasicAuthJwt:
        token = this.cookieService.get(AppConfig.settings.apiServer.auth.token);
        break;
    }
    return token;
  }

  setBasicAuthHeader(headersMap: any): any {
    headersMap["Authorization"] = "Basic " + btoa(`${this.username}:${this.password}`);
    return headersMap;
  }

  login(username: string, password: string): Observable<Boolean> {
    this.operation = "SYSTEM.PROCESS.SIGNIN";
    this.username = username;
    this.password = password;
    let headers: HttpHeaders 
    switch (AppConfig.settings.apiServer.auth.type) {
      case AuthType.BasicAuth:
      case AuthType.BasicAuthJwt:
        headers = new HttpHeaders({
            "Authorization": "Basic " + btoa(`${this.username}:${this.password}`),
        });
        break;
    }
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/login`, {
          headers: headers,
          observe: "response",
          responseType: "json",
      }).pipe(
      catchError(this.handleError<Response>()), 
      map(res => {
        if (res == undefined || res.body == undefined) {
          return false;
        }
        let resp = new RespObject(res.body);
        if (resp == undefined) {
          return false;
        }
        if (resp.code != 0) {
          return false;
        }

        this.currentUser = new SessionUser(resp.data);
        switch (AppConfig.settings.apiServer.auth.type) {
          case AuthType.BasicAuthJwt:
            if (res.headers == undefined) {
              return false
            }
            let token = res.headers.get(AppConfig.settings.apiServer.auth.jwt);
            let expires: number = 3600 * 24 * 1000;
            let date = new Date(Date.now() + expires);
            let cookieOptions: CookieOptions = {
                path: "/",
                expires: date
            };
            this.cookieService.put(AppConfig.settings.apiServer.auth.token, token, cookieOptions);
            break;
        }
        return true;
      }),
    );
  }

  logout(): Observable<Boolean> {
    this.operation = "SYSTEM.PROCESS.SIGNOUT";
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/logout`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res.code != 0) {
          return false;
        } else {
          switch (AppConfig.settings.apiServer.auth.type) {
            case AuthType.BasicAuthJwt:
              this.cookieService.remove(AppConfig.settings.apiServer.auth.token);
            case AuthType.BasicAuth:
              this.cookieService.remove(AppConfig.settings.apiServer.auth.session);
          }
          this.resetUser();
          return true;
        }
      }),
    );
  }

  current(): Observable<SessionUser> {
    this.operation = "SYSTEM.PROCESS.GET";
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/current`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res == undefined) {
          return new SessionUser();
        }

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
