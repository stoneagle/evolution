import { Injectable }               from '@angular/core';
import { HttpClient, HttpHeaders  } from '@angular/common/http';
import { Observable }               from 'rxjs';
import { of }                       from 'rxjs/observable/of';
import { catchError, map, tap  }    from 'rxjs/operators';
import { MessageHandlerService  }   from '../base/message-handler.service';
import { BaseService  }             from '../base/base.service';
import { Area }                     from '../../model/time/area';
import { Resource }                 from '../../model/time/resource';
import { Resp }                     from '../../model/base/resp';

@Injectable()
export class ResourceService extends BaseService {
  constructor(
    protected http: HttpClient,
    protected messageHandlerService: MessageHandlerService,
  ) {
    super(http, messageHandlerService);
    this.resource = this.shareSettings.Time.Resource.Resource;
    this.uri = this.appSettings.apiServer.endpoint + this.appSettings.apiServer.prefix.time + '/resource';
  }

  ListAreas(id: number): Observable<Area[]> {
    return this.BaseResp(this.uri + `/list/areas/${id}`).pipe(map(res => {
      let ret: Area[] = []; 
      res.data.map(
        one => {
          ret.push(new Area(one));
        }
      )
      return ret;
    }))
  }

  ListGroupByLeaf(resource: Resource): Observable<Area[]> {
    return this.BaseList<Area>(resource, Area, this.uri + `/list/leaf`).pipe(map(areas => {
      return areas;
    }))
  }

  List(resource: Resource): Observable<Resource[]> {
    return this.BaseList<Resource>(resource, Resource, this.uri + `/list`).pipe(map(resources => {
      return resources;
    }))
  }
  
  Get(id: number): Observable<Resource> {
    return this.BaseGet<Resource>(Resource, this.uri + `/get/${id}`).pipe(map(resource => {
      return resource;
    }))
  }
  
  Add(resource: Resource): Observable<Resource> {
    return this.BaseAdd<Resource>(resource, Resource, this.uri).pipe(map(resource => {
      return resource;
    }))
  }
  
  Update(resource: Resource): Observable<Resource> {
    return this.BaseUpdate<Resource>(resource, Resource, this.uri + `/${resource.Id}`).pipe(map(resource => {
      return resource;
    }))
  }
  
  Delete(id: number): Observable<Boolean> {
    return this.BaseDelete<Resource>(Resource, this.uri + `/${id}`).pipe(map(resource => {
      return resource;
    }))
  }
}
