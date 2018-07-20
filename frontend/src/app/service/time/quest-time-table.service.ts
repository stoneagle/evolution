import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { QuestTimeTable }           from '../../model/time/quest';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class QuestTimeTableService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.QuestTimeTable;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/quest-timetable';
  }
  
  List(questTimeTable: QuestTimeTable): Observable<QuestTimeTable[]> {
    return this.BaseList<QuestTimeTable>(questTimeTable, QuestTimeTable, this.uri + `/list`).pipe(map(questTimeTables => {
      return questTimeTables;
    }))
  }
  
  Get(id: number): Observable<QuestTimeTable> {
    return this.BaseGet<QuestTimeTable>(QuestTimeTable, this.uri + `/get/${id}`).pipe(map(questTimeTable => {
      return questTimeTable;
    }))
  }
  
  Add(questTimeTable: QuestTimeTable): Observable<QuestTimeTable> {
    return this.BaseAdd<QuestTimeTable>(questTimeTable, QuestTimeTable, this.uri).pipe(map(questTimeTable => {
      return questTimeTable;
    }))
  }
  
  Update(questTimeTable: QuestTimeTable): Observable<QuestTimeTable> {
    return this.BaseUpdate<QuestTimeTable>(questTimeTable, QuestTimeTable, this.uri + `/${questTimeTable.Id}`).pipe(map(questTimeTable => {
      return questTimeTable;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<QuestTimeTable>(QuestTimeTable, this.uri + `/${id}`).pipe(map(questTimeTable => {
      return questTimeTable;
    }))
  }
}
