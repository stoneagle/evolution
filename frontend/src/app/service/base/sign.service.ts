import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';

import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { BaseService  } from '../base/base.service';
import { MessageHandlerService  } from '../base/message-handler.service';

@Injectable()
export class SignService extends BaseService {

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService ,
  ) { 
    super(http, messageHandlerService)
  }

  currentUser: string = null;

  clear(): void {
    this.currentUser = null;
  }

  getCurrentUser(): string {
    return this.currentUser;
  }

  login(username: string, password: string): Observable<Response> {
    this.operation = "SYSTEM.USER.SIGNIN";

    let httpOptions = {
      headers: new HttpHeaders({ 
        'Authorization': "Basic " + btoa(`${username}:${password}`) , 
        'Content-Type': 'application/x-www-form-urlencoded',
        'Access-Control-Allow-Origin': '*'
      })
    };

    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + `/login`, httpOptions).pipe(
      catchError(this.handleError<Response>())
    );
  }

  logout(): Observable<Response> {
    this.operation = "SYSTEM.USER.SIGNOUT";

    let httpOptions = {
      headers: new HttpHeaders({ 
        'Content-Type': 'application/x-www-form-urlencoded',
        'Access-Control-Allow-Origin': '*'
      })
    };

    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + `/logout`, httpOptions).pipe(
      catchError(this.handleError<Response>())
    );
  }
}
