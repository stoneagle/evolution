import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Task }                     from '../../model/time/task';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class TaskService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Task;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/task';
  }

  ListByUser(userId: number): Observable<Task[]> {
    return this.BaseList<Task>(null, Task, this.uri + `/list/user/${userId}`).pipe(map(tasks => {
      return tasks;
    }))
  }
  

  List(task: Task): Observable<Task[]> {
    return this.BaseList<Task>(task, Task, this.uri + `/list`).pipe(map(tasks => {
      return tasks;
    }))
  }
  
  Get(id: number): Observable<Task> {
    return this.BaseGet<Task>(Task, this.uri + `/get/${id}`).pipe(map(task => {
      return task;
    }))
  }
  
  Add(task: Task): Observable<Task> {
    return this.BaseAdd<Task>(task, Task, this.uri).pipe(map(task => {
      return task;
    }))
  }
  
  Update(task: Task): Observable<Task> {
    return this.BaseUpdate<Task>(task, Task, this.uri + `/${task.Id}`).pipe(map(task => {
      return task;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Task>(Task, this.uri + `/${id}`).pipe(map(task => {
      return task;
    }))
  }
}
