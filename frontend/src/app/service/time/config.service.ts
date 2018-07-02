import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { TreeModel }                from 'ng2-tree';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class ConfigService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/config';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.CONFIG.CONCEPT';
  }

  FieldMap(): Observable<Map<number, string>> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/field`).pipe(
      catchError(this.handleError<Response>()),
      map(res => {
        let ret:Map<number, string> = new Map(); 
        if (res && res.code == 0) {
          for (let key in res.data) {
            ret.set(+key, res.data[key]);
          }
        }
        return ret; 
      }),
    )
  }
}
