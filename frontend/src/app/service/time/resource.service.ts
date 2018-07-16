import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';
import { Area }                     from '../../model/time/area';
import { Resource }                 from '../../model/time/resource';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ResourceService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/resource';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Resource;
  }

  List(): Observable<Resource[]> {
    this.operation = this.shareSettings.System.Process.List;
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

  ListAreas(id: number): Observable<Area[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/areas/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Area[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Area(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  ListWithCondition(resource: Resource): Observable<Resource[]> {
    this.operation = this.shareSettings.System.Process.List;
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

  ListGroupByLeaf(resource: Resource): Observable<Area[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/leaf`, JSON.stringify(resource)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Area[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Area(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Resource> {
    this.operation = this.shareSettings.System.Process.Get;
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
    this.operation = this.shareSettings.System.Process.Create;
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
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${resource.Id}`, JSON.stringify(resource)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
    );
  }

  Delete(id: number): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Delete;
    return this.http.delete<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${id}`).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>())
    );
  }
}
