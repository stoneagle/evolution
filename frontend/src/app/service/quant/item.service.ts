import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { WebsocketService  }        from '../base/websocket.service';
import { BaseService  }             from '../base/base.service';
import { Item }                     from '../../model/quant/item';
import { Resp }                     from '../../model/base/resp';
import { Classify }                 from '../../model/quant/classify';
import { AssetSource }              from '../../model/quant/config';
import { Event }                    from '../../model/base/socket';

@Injectable()
export class ItemService extends BaseService {
  private connection;
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
    private websocketService: WebsocketService,
  ) { 
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Quant.Resource.Item;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/item';
  }

  SyncClassify(classify: Classify): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Sync;
    return this.http.post<Resp>(this.uri + `/sync/classify`, JSON.stringify(classify)).pipe(
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
    this.operation = this.shareSettings.System.Process.Sync;
    return this.http.get<Resp>(this.uri + `/point/${id}`).pipe(
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
    this.operation = this.shareSettings.System.Process.Sync;

    let url = this.appSettings.apiServer.websocket + this.uri + `/sync/classify/ws`;
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

  List(): Observable<Item[]> {
    return this.BaseList<Item>(Item, this.uri + `/list`).pipe(map(items => {
      return items;
    }))
  }
  
  ListWithCondition(item: Item): Observable<Item[]> {
    return this.BaseListWithCondition<Item>(item, Item, this.uri + `/list`).pipe(map(items => {
      return items;
    }))
  }
  
  Get(id: number): Observable<Item> {
    return this.BaseGet<Item>(Item, this.uri + `/get/${id}`).pipe(map(item => {
      return item;
    }))
  }
  
  Add(item: Item): Observable<Item> {
    return this.BaseAdd<Item>(item, Item, this.uri).pipe(map(item => {
      return item;
    }))
  }
  
  Update(item: Item): Observable<Item> {
    return this.BaseUpdate<Item>(item, Item, this.uri + `/${item.Id}`).pipe(map(item => {
      return item;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Item>(Item, this.uri + `/${id}`).pipe(map(item => {
      return item;
    }))
  }
}
