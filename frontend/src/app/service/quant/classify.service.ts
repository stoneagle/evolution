import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Classify }                 from '../../model/quant/classify';
import { Response }                 from '../../model/base/response.model';
import { AssetSource }              from '../../model/quant/config';

@Injectable()
export class ClassifyService extends BaseService {
  private uri = '/classify';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'RESOURCE.CLASSIFY.CONCEPT';
  }

  List(): Observable<Classify[]> {
    this.operation = 'PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
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
    this.operation = 'PROCESS.LIST';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(assetSource)).pipe(
      catchError(this.handleError<Response>()),
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
    this.operation = 'PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
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
    this.operation = 'PROCESS.SYNC';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + '/sync', JSON.stringify(assetSource)).pipe(
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

  Delete(id: number): Observable<Response> {
    this.operation = 'PROCESS.DELETE';
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>())
    );
  }
}
