import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Pool } from '../model/pool';
import { Setting } from '../constants/app.setting';

const httpOptions = {
  headers: new HttpHeaders({ 
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*'
  })
};

@Injectable()
export class PoolService {
  private poolUrl = '/pool';

  constructor(
    private http: HttpClient,
  ) { }

  getPools(): Observable<Pool[]> {
    return this.http.get<Pool[]>(Setting.BACKEND_ENDPOINT + this.poolUrl, httpOptions).pipe(
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
