import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Pool } from '../model/pool';
import { Response } from '../model/base/response.model';
import { Setting } from '../constants/app.setting';

@Injectable()
export class PoolService {
  private poolUrl = '/pool';

  constructor(
    private http: HttpClient,
  ) { }

  getPools(): Observable<Pool[]> {
    return this.http.get<Response>(Setting.BACKEND_ENDPOINT + this.poolUrl).pipe(
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
      tap(pools => this.log(`fetched pools`)),
      catchError(this.handleError('getPools', []))
    )
  }

  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error);
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    }
  }

  private log(message: string) {
    // this.messageService.add('PoolService:' + message);
  }
}
