import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { ShareSettings }            from '../../shared/settings';
import { BaseService  }             from '../base/base.service';
import { Country }                  from '../../model/time/country';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class CountryService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/country';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
    protected shareSettings: ShareSettings,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Country;
  }

  List(): Observable<Country[]> {
    this.operation = this.shareSettings.System.Process.List;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        let ret:Country[] = []; 
        if (res && res.code == 0) {
          res.data.map(
            one => {
              ret.push(new Country(one));
            }
          )
        }
        return ret; 
      }),
    )
  }

  Get(id: number): Observable<Country> {
    this.operation = this.shareSettings.System.Process.Get;
    return this.http.get<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Country(res.data);
        } else {
          return new Country();
        }
      }),
    )
  }

  Add(country: Country): Observable<Country> {
    this.operation = this.shareSettings.System.Process.Create;
    return this.http.post<Resp>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(country)).pipe(
      tap(res => this.log(res)),
      catchError(this.handleError<Resp>()),
      map(res => {
        if (res && res.code == 0) {
          return new Country(res.data);
        } else {
          return new Country();
        }
      }),
    );
  }

  Update(country: Country): Observable<Resp> {
    this.operation = this.shareSettings.System.Process.Update;
    return this.http.put<Resp>(AppConfig.settings.apiServer.endpoint + this.uri + `/${country.Id}`, JSON.stringify(country)).pipe(
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
