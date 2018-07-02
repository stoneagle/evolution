import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { TreeModel }                from 'ng2-tree';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Area }                     from '../../model/time/area';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class AreaService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/area';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.AREA.CONCEPT';
  }

  List(): Observable<Map<number, TreeModel>> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Map<number, TreeModel> = new Map(); 
        if (res && res.code == 0) {
          for (let key in res.data) {
            ret.set(+key, res.data[key]);
          }
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Area> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Area(res.data);
        } else {
          return new Area();
        }
      }),
    )
  }

  Add(area: Area): Observable<Area> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(area)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Area(res.data);
        } else {
          return new Area();
        }
      }),
    );
  }

  Update(area: Area): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${area.Id}`, JSON.stringify(area)).pipe(
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
