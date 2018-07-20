import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { TreeModel }                from 'ng2-tree';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Area }                     from '../../model/time/area';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class AreaService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Area;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/area';
  }

  List(area: Area): Observable<Area[]> {
    return this.BaseList<Area>(area, Area, this.uri + `/list`).pipe(map(areas => {
      return areas;
    }))
  }

  ListAllTree(): Observable<Map<number, TreeModel>> {
    return this.BaseGet<Resp>(Resp, this.uri + `/list/tree/all`).pipe(map(res => {
      let ret:Map<number, TreeModel> = new Map(); 
      for (let key in res.data) {
        ret.set(+key, res.data[key]);
      }
      return ret; 
    }))
  }

  ListOneTree(fieldId: number): Observable<TreeModel> {
    return this.BaseGet<Resp>(Resp, this.uri + `/list/tree/one/${fieldId}`).pipe(map(res => {
      let ret:TreeModel; 
      ret = res.data;
      return ret; 
    }))
  }

  ListParent(fieldId: number): Observable<Area[]> {
    return this.BaseList<Area>(null, Area, this.uri + `/list/parent/${fieldId}`).pipe(map(areas => {
      return areas; 
    }))
  }

  ListChildren(parentId: number): Observable<Area[]> {
    return this.BaseList<Area>(null, Area, this.uri + `/list/children/${parentId}`).pipe(map(areas => {
      return areas; 
    }))
  }

  Get(id: number): Observable<Area> {
    return this.BaseGet<Area>(Area, this.uri + `/get/${id}`).pipe(map(area => {
      return area;
    }))
  }

  Add(area: Area): Observable<Area> {
    return this.BaseAdd<Area>(area, Area, this.uri).pipe(map(area => {
      return area;
    }))
  }

  Update(area: Area): Observable<Area> {
    return this.BaseUpdate<Area>(area, Area, this.uri + `/${area.Id}`).pipe(map(area => {
      return area;
    }))
  }

  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Area>(Area, this.uri + `/${id}`).pipe(map(area => {
      return area;
    }))
  }
}
