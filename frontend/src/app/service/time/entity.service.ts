import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Entity }                 from '../../model/time/entity';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class EntityService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/entity';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.ENTITY.CONCEPT';
  }

  List(): Observable<Entity[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Entity[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Entity(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(entity: Entity): Observable<Entity[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`, JSON.stringify(entity)).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Entity[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Entity(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Entity> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Entity(res.data);
        } else {
          return new Entity();
        }
      }),
    )
  }

  Add(entity: Entity): Observable<Entity> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(entity)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Entity(res.data);
        } else {
          return new Entity();
        }
      }),
    );
  }

  Update(entity: Entity): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${entity.Id}`, JSON.stringify(entity)).pipe(
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
