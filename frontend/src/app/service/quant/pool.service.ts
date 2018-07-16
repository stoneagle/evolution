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
  private uri = '/pool';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) { 
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Quant.Resource.Pool;
  }

  List(): Observable<Pool[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
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
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/item`, JSON.stringify(pool)).pipe(
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

  Get(id: number): Observable<Pool> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          let pool = new Pool(res.data);
          if (pool.Item == null) {
            pool.Item = [];
          }
          return pool;
        } else {
          return new Pool();
        }
      }),
    )
  }

  Add(pool: Pool): Observable<Pool> {
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(pool)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
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
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + "/items/add", JSON.stringify(pool)).pipe(
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
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/items/delete`, JSON.stringify(pool)).pipe(
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

  Update(pool: Pool): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${pool.Id}`, JSON.stringify(pool)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
    );
  }

  Delete(id: number): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Delete;
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
