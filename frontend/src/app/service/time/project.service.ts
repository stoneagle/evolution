import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { Project }                  from '../../model/time/project';
import { Response }                 from '../../model/base/response.model';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';

@Injectable()
export class ProjectService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/project';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) { 
    super(http, messageHandlerService);
    this.resource = 'FLOW.RESOURCE.POOL.CONCEPT';
  }

  List(): Observable<Project[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Project[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Project(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Project> {
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          let project = new Project(res.data);
          return project;
        } else {
          return new Project();
        }
      }),
    )
  }

  Add(project: Project): Observable<Project> {
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(project)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Project(res.data);
        } else {
          return new Project();
        }
      }),
    );
  }

  Update(project: Project): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${project.id}`, JSON.stringify(project)).pipe(
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
