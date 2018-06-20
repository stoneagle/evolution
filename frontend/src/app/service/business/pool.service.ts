import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Pool } from '../../model/business/pool';
import { Item } from '../../model/business/item';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { MessageHandlerService  } from '../base/message-handler.service';
import { BaseService  } from '../base/base.service';

@Injectable()
export class PoolService extends BaseService {
  private uri = '/pool';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) { 
    super(http, messageHandlerService);
    this.resource = 'RESOURCE.POOL.CONCEPT';
  }

  List(): Observable<Pool[]> {
    this.operation = 'PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Pool[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Pool(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListItem(pool: Pool): Observable<Item[]> {
    this.operation = 'PROCESS.LIST';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/item`, JSON.stringify(pool)).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Item[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Item(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Pool> {
    this.operation = 'PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Pool(res.data);
        } else {
          return new Pool();
        }
      }),
    )
  }

  Add(pool: Pool): Observable<Pool> {
    this.operation = 'PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(pool)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Pool(res.data);
        } else {
          return new Pool();
        }
      }),
    );
  }

  AddItems(pool: Pool): Observable<Boolean> {
    this.operation = 'PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + "/items/add", JSON.stringify(pool)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return true;
        } else {
          return false;
        }
      }),
    );
  }

  DeleteItems(pool: Pool): Observable<Boolean> {
    this.operation = 'PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/items/delete`, JSON.stringify(pool)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return true;
        } else {
          return false;
        }
      }),
    );
  }

  Update(pool: Pool): Observable<Response> {
    this.operation = 'PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${pool.Id}`, JSON.stringify(pool)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
    );
  }

  Delete(id: number): Observable<Response> {
    this.operation = 'PROCESS.DELETE';
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>())
    );
  }
}
