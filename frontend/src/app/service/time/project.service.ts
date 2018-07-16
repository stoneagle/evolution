import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Project }                  from '../../model/time/project';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ProjectService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Project;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/project';
  }

  List(): Observable<Project[]> {
    return this.BaseList<Project>(Project, this.uri + `/list`).pipe(map(projects => {
      return projects;
    }))
  }

  ListWithCondition(project: Project): Observable<Project[]> {
    return this.BaseListWithCondition<Project>(project, Project, this.uri + `/list`).pipe(map(projects => {
      return projects;
    }))
  }

  Get(id: number): Observable<Project> {
    return this.BaseGet<Project>(Project, this.uri + `/get/${id}`).pipe(map(project => {
      return project;
    }))
  }

  Add(project: Project): Observable<Project> {
    return this.BaseAdd<Project>(project, Project, this.uri).pipe(map(project => {
      return project;
    }))
  }

  Update(project: Project): Observable<Project> {
    return this.BaseUpdate<Project>(project, Project, this.uri + `/${project.Id}`).pipe(map(project => {
      return project;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Project>(Project, this.uri + `/${id}`).pipe(map(project => {
      return project;
    }))
  }
}
