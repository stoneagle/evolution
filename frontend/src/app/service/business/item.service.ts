import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Item } from '../../model/business/item';
import { Response } from '../../model/base/response.model';
import { WsStatus } from '../../shared/shared.const';
import { AppConfig } from '../base/config.service';
import { Classify } from '../../model/business/classify';
import { AssetSource } from '../../model/config/config';
import { MessageHandlerService  } from '../base/message-handler.service';
import { WebsocketService  } from '../base/websocket.service';
import { Event } from '../../model/base/socket';
import { SocketService  } from '../base/socket.service';

@Injectable()
export class ItemService {
  private uri = '/item';
  public connection;

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
    private socketService: SocketService,
    private websocketService: WebsocketService,
  ) { }

  List(): Observable<Item[]> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>('PROCESS.LIST')),
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
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>('PROCESS.GET')),
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
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/sync/classify`, JSON.stringify(classify)).pipe(
      tap(res => this.log('PROCESS.SYNC', res)),
      catchError(this.handleError<Response>('PROCESS.SYNC')),
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
    let url = AppConfig.settings.apiServer.websocket + this.uri + `/sync/classify/ws`;
    this.connection = this.websocketService.connect(url);

    let index = 0;
    let sub = this.connection.asObservable().map((response: MessageEvent): Boolean => {
      let res = JSON.parse(response.data);
      console.log(res);

      if (!res || res.code != 0) {
        this.messageHandlerService.showWarning(res.desc);
        this.connection.complete();
        return false;
      }

      if (res.status != WsStatus.Message) {
        this.messageHandlerService.showWarning(res.desc);
      } else if (res.status != WsStatus.Message) {
        this.messageHandlerService.showSuccess(res.data);
      }
      if (index < classifyBatch.length) {
        delay(2000).then(any=>{
          this.connection.next(classifyBatch[index]);
        });
        index++;
        return false;
      } else {
        return true;
      }
    });
    return sub;
  }

  Delete(id: number): Observable<Response> {
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log('PROCESS.DELETE', res)),
      catchError(this.handleError<Response>('PROCESS.DELETE'))
    );
  }

  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      this.messageHandlerService.handleError(error);
      return of(result as T);
    }
  }

  private log(message: string, res: Response) {
    if (res.code != 0) {
      this.messageHandlerService.showWarning(res.desc);
    } else {   
      this.messageHandlerService.showSuccess(message);
    }   
  }
}

async function delay(ms: number) {
  await new Promise(resolve => setTimeout(
      () => resolve(), 1000
    )
  ).then(
    () => console.log("fired")
  );
}

