import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { SignService  }             from '../system/sign.service';
import { Kanban }                     from '../../model/time/syncfusion';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class SyncfusionService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/syncfusion';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected signService: SignService,
  ) {
    super(http, messageHandlerService);
  }

  ListKanban(): Observable<Kanban[]> {
    this.resource = 'TIME.RESOURCE.TASK.CONCEPT';
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list/kanban`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Kanban[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Kanban(one));
            }
          )
        }
        return ret; 
      }),
    )
  }
}
