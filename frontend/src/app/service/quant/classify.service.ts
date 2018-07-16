import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Classify }                 from '../../model/quant/classify';
import { Resp }                     from '../../model/base/resp';
import { AssetSource }              from '../../model/quant/config';

@Injectable()
export class ClassifyService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Quant.Resource.Classify;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/classify';
  }

  ListByAssetSource(assetSource: AssetSource): Observable<Classify[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(this.appSettings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(assetSource)).pipe(
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

  Sync(assetSource: AssetSource): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Sync;
    return this.http.post<Resp>(this.appSettings.apiServer.endpoint + this.uri + '/sync', JSON.stringify(assetSource)).pipe(
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

  List(): Observable<Classify[]> {
    return this.BaseList<Classify>(Classify, this.uri + `/list`).pipe(map(classifies => {
      return classifies;
    }))
  }
  
  ListWithCondition(classify: Classify): Observable<Classify[]> {
    return this.BaseListWithCondition<Classify>(classify, Classify, this.uri + `/list`).pipe(map(classifies => {
      return classifies;
    }))
  }
  
  Get(id: number): Observable<Classify> {
    return this.BaseGet<Classify>(Classify, this.uri + `/get/${id}`).pipe(map(classify => {
      return classify;
    }))
  }
  
  Add(classify: Classify): Observable<Classify> {
    return this.BaseAdd<Classify>(classify, Classify, this.uri).pipe(map(classify => {
      return classify;
    }))
  }
  
  Update(classify: Classify): Observable<Classify> {
    return this.BaseUpdate<Classify>(classify, Classify, this.uri + `/${classify.Id}`).pipe(map(classify => {
      return classify;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Classify>(Classify, this.uri + `/${id}`).pipe(map(classify => {
      return classify;
    }))
  }
}
