import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Resource }                 from '../../model/time/resource';
import { Resp }                 from '../../model/base/resp';

@Injectable()
export class ResourceService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/resource';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.RESOURCE.CONCEPT';
  }

  List(): Observable<Resource[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Resource[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Resource(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(resource: Resource): Observable<Resource[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(resource)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Resource[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Resource(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Resource> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Resource(res.data);
        } else {
          return new Resource();
        }
      }),
    )
  }

  Add(resource: Resource): Observable<Resource> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(resource)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Resource(res.data);
        } else {
          return new Resource();
        }
      }),
    );
  }

  Update(resource: Resource): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${resource.Id}`, JSON.stringify(resource)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
    );
  }

  Delete(id: number): Observable<Resp> {
    this.operation = 'SYSTEM.PROCESS.DELETE';
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
