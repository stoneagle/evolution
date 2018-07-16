import { Injectable }                            from '@angular/core';
import { HttpClient, HttpHeaders, HttpRequest  } from '@angular/common/http';
import { Observable }                            from 'rxjs';
import { of }                                    from 'rxjs/observable/of';
import { catchError, map, tap  }                 from 'rxjs/operators';
import { CookieService, CookieOptions }          from 'ngx-cookie';
import { MessageHandlerService  }                from '../base/message-handler.service';
import { BaseService  }                          from '../base/base.service';
import { SessionUser }                           from '../../model/base/sign';
import { Resp, RespObject }                      from '../../model/base/resp';
import { AuthType }                              from '../../shared/const';

@Injectable()
export class SignService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService ,
    private cookieService: CookieService,
  ) { 
    super(http, messageHandlerService);
    this.resource = this.shareSettings.System.Resource.User;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.system + '/sign';
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
    switch (this.appSettings.apiServer.auth.type) {
      case AuthType.BasicAuthJwt:
        token = this.cookieService.get(this.appSettings.apiServer.auth.token);
        break;
    }
    return token;
  }

  setBasicAuthHeader(headersMap: any): any {
    headersMap["Authorization"] = "Basic " + btoa(`${this.username}:${this.password}`);
    return headersMap;
  }

  login(username: string, password: string): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Signin;
    this.username = username;
    this.password = password;
    let headers: HttpHeaders 
    switch (this.appSettings.apiServer.auth.type) {
      case AuthType.BasicAuth:
      case AuthType.BasicAuthJwt:
        headers = new HttpHeaders({
            "Authorization": "Basic " + btoa(`${this.username}:${this.password}`),
        });
        break;
    }
    return this.http.get<Response>(this.uri + `/login`, {
          headers: headers,
          observe: "response",
          responseType: "json",
      }).pipe(
      catchError(this.handleError<Response>()), 
      map(res => {
        if (res == undefined || res.body == undefined) {
          let msg = this.errorInfo.Server.Exception.NoResponse;
          this.messageHandlerService.showError(this.resource, this.operation, msg);
          throw(this.errorInfo.Server.Exception.NoResponse);
        }
        let resp = new RespObject(res.body);
        this.handleResponse(resp)
        this.currentUser = new SessionUser(resp.data);
        switch (this.appSettings.apiServer.auth.type) {
          case AuthType.BasicAuthJwt:
            if (res.headers == undefined) {
              let msg = this.errorInfo.Server.Exception.NoResponse;
              this.messageHandlerService.showError(this.resource, this.operation, msg);
              throw(this.errorInfo.Server.Exception.NoResponse);
            }
            let token = res.headers.get(this.appSettings.apiServer.auth.jwt);
            let expires: number = 3600 * 24 * 1000;
            let date = new Date(Date.now() + expires);
            let cookieOptions: CookieOptions = {
                path: "/",
                expires: date
            };
            this.cookieService.put(this.appSettings.apiServer.auth.token, token, cookieOptions);
            break;
        }
        return true;
      }),
    );
  }

  logout(): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Signout;
    return this.http.get<Resp>(this.uri + `/logout`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        this.handleResponse(res)
        switch (this.appSettings.apiServer.auth.type) {
          case AuthType.BasicAuthJwt:
            this.cookieService.remove(this.appSettings.apiServer.auth.token);
          case AuthType.BasicAuth:
            this.cookieService.remove(this.appSettings.apiServer.auth.session);
        }
        this.resetUser();
        return true;
      }),
    );
  }

  current(): Observable<SessionUser> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(this.uri + `/current`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        this.handleResponse(res)
        this.currentUser = new SessionUser(res.data);
        return this.currentUser;
      }),
    );
  }
}
