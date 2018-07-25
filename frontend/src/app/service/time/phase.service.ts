import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Phase }                    from '../../model/time/phase';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class PhaseService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Phase;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/phase';
  }

  List(phase: Phase): Observable<Phase[]> {
    return this.BaseList<Phase>(phase, Phase, this.uri + `/list`).pipe(map(phases => {
      return phases;
    }))
  }

  Get(id: number): Observable<Phase> {
    return this.BaseGet<Phase>(Phase, this.uri + `/get/${id}`).pipe(map(phase => {
      return phase;
    }))
  }

  Add(phase: Phase): Observable<Phase> {
    return this.BaseAdd<Phase>(phase, Phase, this.uri).pipe(map(phase => {
      return phase;
    }))
  }

  Update(phase: Phase): Observable<Phase> {
    return this.BaseUpdate<Phase>(phase, Phase, this.uri + `/${phase.Id}`).pipe(map(phase => {
      return phase;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Phase>(Phase, this.uri + `/${id}`).pipe(map(phase => {
      return phase;
    }))
  }
}
