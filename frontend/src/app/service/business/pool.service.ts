import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Pool } from '../../model/business/pool';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { MessageHandlerService  } from '../base/message-handler.service';

@Injectable()
export class PoolService {
  private uri = '/pool';

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
  ) { }

  List(): Observable<Pool[]> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>('PROCESS.LIST')),
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

  Get(id: number): Observable<Pool> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>('PROCESS.GET')),
      map(res => {
        if (res && res.code == 0) {
          return new Pool(res.data);
        } else {
          return new Pool();
        }
      }),
    )
  }

  Add(pool: Pool): Observable<Pool> {
      return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(pool)).pipe(
      tap(res => this.log('PROCESS.ADD', res)),
      catchError(this.handleError<Response>('PROCESS.ADD')),
        map(res => {
          if (res && res.code == 0) {
            return new Pool(res.data);
          } else {
            return new Pool();
          }
        }),
    );
  }

  Update(pool: Pool): Observable<Response> {
      return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${pool.Id}`, JSON.stringify(pool)).pipe(
      tap(res => this.log('PROCESS.UPDATE', res)),
      catchError(this.handleError<Response>('PROCESS.UPDATE')),
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
