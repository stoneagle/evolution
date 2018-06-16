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

@Injectable()
export class ItemService {
  private uri = '/item';

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
  ) { }

  List(): Observable<Item[]> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri).pipe(
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
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
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

  SyncSource(assetSource: AssetSource): Observable<Boolean> {
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/sync/source`, JSON.stringify(assetSource)).pipe(
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
