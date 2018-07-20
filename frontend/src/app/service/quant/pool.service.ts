import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { Pool }                     from '../../model/quant/pool';
import { Item }                     from '../../model/quant/item';
import { Resp }                     from '../../model/base/resp';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';

@Injectable()
export class PoolService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) { 
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Quant.Resource.Pool;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/pool';
  }

  ListItem(pool: Pool): Observable<Item[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(this.appSettings.apiServer.endpoint + this.uri + `/list/item`, JSON.stringify(pool)).pipe(
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

  AddItems(pool: Pool): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(this.appSettings.apiServer.endpoint + this.uri + "/items/add", JSON.stringify(pool)).pipe(
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

  DeleteItems(pool: Pool): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Delete;
    return this.http.post<Resp>(this.appSettings.apiServer.endpoint + this.uri + `/items/delete`, JSON.stringify(pool)).pipe(
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

  List(pool: Pool): Observable<Pool[]> {
    return this.BaseList<Pool>(pool, Pool, this.uri + `/list`).pipe(map(pools => {
      return pools;
    }))
  }
  
  Get(id: number): Observable<Pool> {
    return this.BaseGet<Pool>(Pool, this.uri + `/get/${id}`).pipe(map(pool => {
      return pool;
    }))
  }
  
  Add(pool: Pool): Observable<Pool> {
    return this.BaseAdd<Pool>(pool, Pool, this.uri).pipe(map(pool => {
      return pool;
    }))
  }
  
  Update(pool: Pool): Observable<Pool> {
    return this.BaseUpdate<Pool>(pool, Pool, this.uri + `/${pool.Id}`).pipe(map(pool => {
      return pool;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Pool>(Pool, this.uri + `/${id}`).pipe(map(pool => {
      return pool;
    }))
  }
}
