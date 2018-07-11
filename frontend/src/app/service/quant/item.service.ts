import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { WebsocketService  }        from '../base/websocket.service';
import { BaseService  }             from '../base/base.service';
import { Item }                     from '../../model/quant/item';
import { Resp }                 from '../../model/base/resp';
import { Classify }                 from '../../model/quant/classify';
import { AssetSource }              from '../../model/quant/config';
import { Event }                    from '../../model/base/socket';

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
    this.resource = 'FLOW.RESOURCE.ITEM.CONCEPT';
  }

  List(item?: Item): Observable<Item[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(item)).pipe(
      catchError(this.handleError<Resp>()),
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
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          let item = new Item(res.data);
          if (item.Classify == null) {
            item.Classify = [];
          }
          return item;
        } else {
          return new Item();
        }
      }),
    )
  }

  SyncClassify(classify: Classify): Observable<Boolean> {
    this.operation = 'SYSTEM.PROCESS.SYNC';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/sync/classify`, JSON.stringify(classify)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return true;
        } else {
          return false;
        }
      }),
    );
  }

  SyncPoint(id: number): Observable<Boolean> {
    this.operation = 'SYSTEM.PROCESS.SYNC';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/point/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
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
    this.operation = 'SYSTEM.PROCESS.SYNC';

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

  Delete(id: number): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
