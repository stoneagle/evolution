import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Field }                    from '../../model/time/field';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class FieldService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Field;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/field';
  }

  Map(): Observable<Map<number, string>> {
    return this.BaseList<Field>(Field, this.uri + `/list`).pipe(map(fields => {
      let ret:Map<number, string> = new Map(); 
      fields.forEach((one, k) => {
        ret.set(one.Id, one.Name);
      })
      return ret;
    }))
  }

  List(): Observable<Field[]> {
    return this.BaseList<Field>(Field, this.uri + `/list`).pipe(map(fields => {
      return fields;
    }))
  }

  ListWithCondition(field: Field): Observable<Field[]> {
    return this.BaseListWithCondition<Field>(field, Field, this.uri + `/list`).pipe(map(fields => {
      return fields;
    }))
  }

  Get(id: number): Observable<Field> {
    return this.BaseGet<Field>(Field, this.uri + `/get/${id}`).pipe(map(field => {
      return field;
    }))
  }

  Add(field: Field): Observable<Field> {
    return this.BaseAdd<Field>(field, Field, this.uri).pipe(map(field => {
      return field;
    }))
  }

  Update(field: Field): Observable<Field> {
    return this.BaseUpdate<Field>(field, Field, this.uri + `/${field.Id}`).pipe(map(field => {
      return field;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Field>(Field, this.uri + `/${id}`).pipe(map(field => {
      return field;
    }))
  }
}
