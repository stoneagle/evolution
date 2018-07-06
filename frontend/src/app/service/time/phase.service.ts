import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Phase }                 from '../../model/time/phase';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class PhaseService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/phase';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.PHASE.CONCEPT';
  }

  List(): Observable<Phase[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Phase[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Phase(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(phase: Phase): Observable<Phase[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(phase)).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Phase[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Phase(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Phase> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Phase(res.data);
        } else {
          return new Phase();
        }
      }),
    )
  }

  Add(phase: Phase): Observable<Phase> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(phase)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Phase(res.data);
        } else {
          return new Phase();
        }
      }),
    );
  }

  Update(phase: Phase): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${phase.Id}`, JSON.stringify(phase)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
    );
  }

  Delete(id: number): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>())
    );
  }
}
