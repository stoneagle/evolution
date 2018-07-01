import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { AppConfig }                from '../base/config.service';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Country }                 from '../../model/time/country';
import { Response }                 from '../../model/base/response.model';

@Injectable()
export class CountryService extends BaseService {
  private uri = AppConfig.settings.apiServer.prefix.time + '/country';

  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = 'TIME.RESOURCE.COUNTRY.CONCEPT';
  }

  List(): Observable<Country[]> {
    this.operation = 'SYSTEM.PROCESS.LIST';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/list`).pipe(
      catchError(this.handleError<Response>()),
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
    this.operation = 'SYSTEM.PROCESS.GET';
    return this.http.get<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/get/${id}`).pipe(
      catchError(this.handleError<Response>()),
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
    this.operation = 'SYSTEM.PROCESS.CREATE';
    return this.http.post<Response>(AppConfig.settings.apiServer.endpoint + this.uri, JSON.stringify(country)).pipe(
    tap(res => this.log(res)),
    catchError(this.handleError<Response>()),
      map(res => {
        if (res && res.code == 0) {
          return new Country(res.data);
        } else {
          return new Country();
        }
      }),
    );
  }

  Update(country: Country): Observable<Response> {
    this.operation = 'SYSTEM.PROCESS.UPDATE';
    return this.http.put<Response>(AppConfig.settings.apiServer.endpoint + this.uri + `/${country.Id}`, JSON.stringify(country)).pipe(
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
