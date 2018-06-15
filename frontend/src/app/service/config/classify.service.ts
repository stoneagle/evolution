import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Classify } from '../../model/config/classify';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { MessageHandlerService  } from '../base/message-handler.service';

@Injectable()
export class ClassifyService {
  private uri = '/classify';

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
  ) { }

  List(): Observable<Classify[]> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri).pipe(
      catchError(this.handleError<Response>('PROCESS.LIST')),
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
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      catchError(this.handleError<Response>('PROCESS.GET')),
      map(res => {
        if (res && res.code == 0) {
          return new Classify(res.data);
        } else {
          return new Classify();
        }
      }),
    )
  }

  Sync(classify: Classify): Observable<Boolean> {
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + '/sync', JSON.stringify(classify)).pipe(
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
