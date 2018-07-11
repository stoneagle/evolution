import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Classify }                 from '../../model/quant/classify';
import { Resp }                 from '../../model/base/resp';
import { AssetSource }              from '../../model/quant/config';

@Injectable()
export class ClassifyService extends BaseService {
  private uri = '/classify';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'FLOW.RESOURCE.CLASSIFY.CONCEPT';
  }

  List(): Observable<Classify[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Classify[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Classify(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListByAssetSource(assetSource: AssetSource): Observable<Classify[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(assetSource)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Classify[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Classify(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Classify> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Classify(res.data);
        } else {
          return new Classify();
        }
      }),
    )
  }

  Sync(assetSource: AssetSource): Observable<Boolean> {
    this.operation = 'SYSTEM.PROCESS.SYNC';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + '/sync', JSON.stringify(assetSource)).pipe(
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

  Delete(id: number): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
