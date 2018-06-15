import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap  } from 'rxjs/operators';
import { Response } from '../../model/base/response.model';
import { AppConfig } from '../base/config.service';
import { MessageHandlerService  } from '../base/message-handler.service';

@Injectable()
export class ConfigService {
  private uri = '/config';

  constructor(
    private http: HttpClient,
    private messageHandlerService: MessageHandlerService,
  ) { }

  AssetList(): Observable<Map<string, string>> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + "/asset").pipe(
      catchError(this.handleError<Response>('PROCESS.LIST')),
      map(res => {
        let ret:Map<string, string> = new Map();
        if (res && res.code == 0) {
          for (let key in res.data) {
            ret.set(key, res.data[key]);
          }
        }
        return ret;
      }),
    )
  }

  TypeList(resource: string): Observable<Map<string, Map<string, string[]>>> {
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/type/${resource}`).pipe(
      catchError(this.handleError<Response>('PROCESS.LIST')),
      map(res => {
        let ret:Map<string, Map<string, string[]>> = new Map();
        if (res && res.code == 0) {
          for (let key in res.data) {
            let typeMap: Map<string, string[]> = new Map();
            for (let one in res.data[key]) {
              typeMap.set(one, res.data[key][one])
            }
            ret.set(key, typeMap);
          }
        }
        return ret;
      }),
    )
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
