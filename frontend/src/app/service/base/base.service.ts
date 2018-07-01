import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from './config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { Response }                 from '../../model/base/response.model';
import { WsStatus }                 from '../../shared/shared.const';

@Injectable()
export class BaseService {
  protected resource: string;
  protected operation: string;

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService ,
  ) { 
  }

  protected handleError<T> (result?: T) {
    return (error: any): Observable<T> => {
      this.messageHandlerService.handleError(this.resource, this.operation, error);
      return of(result as T);
    }
  }

  protected log(res: Response, message?: string) {
    if (res.code != 0) {
      this.messageHandlerService.showWarning(this.resource, this.operation, res.desc);
    } else {   
      this.messageHandlerService.showSuccess(this.resource, this.operation, message);
    }   
  }

  protected logWs(res: Response): boolean {
    if (!res || res.code != 0) {
      this.messageHandlerService.showWarning(this.resource, this.operation, res.desc);
      return false;
    }

    if (res.status == WsStatus.Error) {
      this.messageHandlerService.showWarning(this.resource, this.operation, res.desc);
    } else if (res.status == WsStatus.Message) {
      this.messageHandlerService.showInfo(this.resource, this.operation, res.data);
    } 
    return true
  }

  protected async delay(ms: number) {
    await new Promise(resolve => setTimeout(
        () => resolve(), 1000
      )
    ).then(
      () => {
        // console.log("trigger")
      }
    );
  }
}
