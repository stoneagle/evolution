import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Country }                  from '../../model/time/country';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class CountryService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Country;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/country';
  }

  List(): Observable<Country[]> {
    return this.BaseList<Country>(Country, this.uri + `/list`).pipe(map(countries => {
      return countries;
    }))
  }

  ListWithCondition(country: Country): Observable<Country[]> {
    return this.BaseListWithCondition<Country>(country, Country, this.uri + `/list`).pipe(map(countries => {
      return countries;
    }))
  }

  Get(id: number): Observable<Country> {
    return this.BaseGet<Country>(Country, this.uri + `/get/${id}`).pipe(map(country => {
      return country;
    }))
  }

  Add(country: Country): Observable<Country> {
    return this.BaseAdd<Country>(country, Country, this.uri).pipe(map(country => {
      return country;
    }))
  }

  Update(country: Country): Observable<Country> {
    return this.BaseUpdate<Country>(country, Country, this.uri + `/${country.Id}`).pipe(map(country => {
      return country;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Country>(Country, this.uri + `/${id}`).pipe(map(country => {
      return country;
    }))
  }
}
