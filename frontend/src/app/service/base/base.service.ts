import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable, Subject }      from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from './config.service';
import { IConfig }                  from '../../model/base/config';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { Resp }                     from '../../model/base/resp';
import { WsStatus }                 from '../../shared/const';
import { ShareSettings }            from '../../shared/settings';
import { ErrorInfo }                from '../../shared/error';
import { Serializable }             from '../../model/base/serializable';

@Injectable()
export class BaseService {
  protected resource: string;
  protected operation: string;
  protected shareSettings: ShareSettings;
  protected appSettings: IConfig;
  protected uri: string;
  protected errorInfo: ErrorInfo;

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) { 
    this.shareSettings = new ShareSettings();
    this.appSettings = AppConfig.settings;
    this.errorInfo = new ErrorInfo();
  }

  protected BaseGet<T extends Serializable >(c: new (any) => T, url: string): Observable<T> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(url).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new c(res.data);
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BaseResp(url: string): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(url).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return res;
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BasePostResp(params: any, url: string): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.post<Resp>(url, JSON.stringify(params)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return res;
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BaseList<T extends Serializable >(params: any, c: new (any) => T, url: string): Observable<T[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.post<Resp>(url, JSON.stringify(params)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          let ret: T[] = []; 
          res.data.map(
            one => {
              ret.push(new c(one));
            }
          )
          return ret
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  public Count(params: any): Observable<number> {
    let url = this.uri + "/count";
    this.operation = this.shareSettings.System.Process.Count;
    return this.http.post<Resp>(url, JSON.stringify(params)).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return res.data;
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BaseAdd<T extends Serializable >(params: any, c: new (any) => T, url: string): Observable<T> {
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(url, JSON.stringify(params)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new c(res.data);
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BaseUpdate<T extends Serializable >(params: any, c: new (any) => T, url: string): Observable<T> {
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(url, JSON.stringify(params)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new c(res.data);
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected BaseDelete<T extends Serializable >(c: new (any) => T, url: string): Observable<Boolean> {
    this.operation = this.shareSettings.System.Process.Delete;
    return this.http.delete<Resp>(url).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return true;
        } else {
          this.handleResponse(res);
        }
      }),
    )
  }

  protected handleError<T> (result?: T) {
    return (error: any): Observable<T> => {
      this.messageHandlerService.handleError(this.resource, this.operation, error);
      return of(result as T);
    }
  }

  protected handleResponse(res: Resp): void {
    if (res == undefined) {
      let msg = this.errorInfo.Server.Exception.NoResponse;
      this.messageHandlerService.showError(this.resource, this.operation, msg);
      throw(this.errorInfo.Server.Exception.NoResponse);
    }
    if (res.code != 0) {
      let msg = this.errorInfo.Server.Error[res.code] + ":" + res.desc;
      this.messageHandlerService.showError(this.resource, this.operation, msg);
      throw(msg);
    }
    return
  }

  protected log(res: Resp, message?: string) {
    if (res.code != 0) {
      this.messageHandlerService.showWarning(this.resource, this.operation, res.desc);
    } else {   
      this.messageHandlerService.showSuccess(this.resource, this.operation, message);
    }   
  }

  protected logWs(res: Resp): boolean {
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

  protected DateJsonKeepFormat = function() {
    let timezoneOffsetInHours = -(this.getTimezoneOffset() / 60); //UTC minus local time
    let sign = timezoneOffsetInHours >= 0 ? '+' : '-';
    let leadingZero = (timezoneOffsetInHours < 10) ? '0' : '';
    let correctedDate = new Date(this.getFullYear(), this.getMonth(), 
    this.getDate(), this.getHours(), this.getMinutes(), this.getSeconds(), 
    this.getMilliseconds());
    correctedDate.setHours(this.getHours() + timezoneOffsetInHours);
    let iso = correctedDate.toISOString().replace('Z', '');
    return iso + sign + leadingZero + Math.abs(timezoneOffsetInHours).toString() + ':00';
  }
}
