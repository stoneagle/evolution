import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Item } from '../../model/business/item';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { Classify } from '../../model/business/classify';
import { AssetSource } from '../../model/config/config';
import { MessageHandlerService  } from '../base/message-handler.service';
import { WebsocketService  } from '../base/websocket.service';
import { Event } from '../../model/base/socket';
import { BaseService  } from '../base/base.service';

@Injectable()
export class ItemService extends BaseService {
  private uri = '/item';
  private connection;

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    private websocketService: WebsocketService,
  ) { 
    super(http, messageHandlerService);
    this.resource = 'RESOURCE.ITEM.CONCEPT';
  }

  List(): Observable<Item[]> {
    this.operation = 'PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
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

  Get(id: number): Observable<Item> {
    this.operation = 'PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Item(res.data);
        } else {
          return new Item();
        }
      }),
    )
  }

  SyncClassify(classify: Classify): Observable<Boolean> {
    this.operation = 'PROCESS.SYNC';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/sync/classify`, JSON.stringify(classify)).pipe(
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

  WsSyncSource(classifyBatch: Classify[]): Observable<Boolean> {
    this.operation = 'PROCESS.SYNC';

    let url = AppConfig.settings.apiServer.websocket + this.uri + `/sync/classify/ws`;
    this.connection = this.websocketService.connect(url);

    let index = 0;
    let sub = this.connection.asObservable().map((response: MessageEvent): Boolean => {
      let res = JSON.parse(response.data);
      if (!this.logWs(res)) {
        this.connection.complete();
        this.connection = null;
        return false;
      }

      if (index < classifyBatch.length) {
        this.delay(2000).then(any=>{
          this.connection.next(classifyBatch[index]);
          index++;
        });
        return false;
      } else {
        this.log(res);
        return true;
      }
    });
    return sub;
  }

  Delete(id: number): Observable<Response> {
    this.operation = 'PROCESS.DELETE';
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>())
    );
  }
}
