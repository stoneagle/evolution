import { Component, OnInit, ViewChild, Input } from '@angular/core';
import { Comparator, State, SortOrder}         from "clarity-angular";

import { PageSet }               from '../../../../model/base/basic';
import { Resource }              from '../../../../model/time/resource';
import { Area }                  from '../../../../model/time/area';
import { AreaService }           from '../../../../service/time/area.service';
import { ResourceService }       from '../../../../service/time/resource.service';
import { ResourceSaveComponent } from './../save/save.component';

import { CustomComparator }                             from '../../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../../shared/utils';
import { PageSize }                                     from '../../../../shared/const';

@Component({
  selector: 'time-resource-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ResourceListComponent implements OnInit {
  @ViewChild(ResourceSaveComponent)
  saveResourceComponent: ResourceSaveComponent;

  @Input() initAllFlag: boolean = false;
  filterAreaId: number;
  filterAreaIds: number[] = [];

  resources: Resource[];
  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private areaService: AreaService,
    private resourceService: ResourceService,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: any): void {
    if ((this.initAllFlag) || (this.filterAreaId)) {
      let resource = new Resource();
      if (this.filterAreaIds.length > 0) {
        resource.Area.Ids = this.filterAreaIds;
      }
      resource = loadPageFilterSort<Resource>(resource, state);
      this.pageSet.Current = resource.Page.Current;
      this.currentState = state;
      this.resourceService.Count(resource).subscribe(count => {
        this.pageSet.Count = count;
        this.resourceService.List(resource).subscribe(res => {
          this.resources = res;
        })
      })
    }
  }

  setFilterAreaId(areaId: number) {
    this.filterAreaId = areaId;
    this.areaService.ListAllLeaf(this.filterAreaId).subscribe(areaIds => {
      areaIds.push(this.filterAreaId);
      this.filterAreaIds = areaIds;
    })
  }

  saved(savedResource: Resource): void {
    if (savedResource.Id) {
      this.load(this.currentState);
    } else {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveResourceComponent.New(this.filterAreaId, id);
  }

  delete(resource: Resource): void {
    this.resourceService.Delete(resource.Id).subscribe(res => {
      let state = deleteState(this.pageSet, this.currentState, 1);
      this.load(state);
    })
  }
}
