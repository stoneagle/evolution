import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Pool } from '../../model/pool';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { MessageHandlerService  } from '../base/message-handler.service';

@Injectable()
export class PoolService {
  private poolUrl = '/pool';

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
  ) { }

  getPools(): Observable<Pool[]> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.poolUrl).pipe(
      tap(res => this.log('PROCESS.ROLLBACK', res)),
      catchError(this.handleError<Response>('getPools')),
      map(res => {
        let ret:Pool[] = []; 
        if (res) {
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

  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      // this.messageHandlerService.showError(`${operation} failed: ${error.message}`, '');
      this.messageHandlerService.handleError(error);
      return of(result as T);
    }
  }

  private log(message: string, res: Response) {
    if (res.code != 0) { 
      this.messageHandlerService.showWarning(res.desc)
    } else {   
      this.messageHandlerService.showSuccess(message)
    }   
  }
}
