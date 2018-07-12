import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Resource }              from '../../../../model/time/resource';
import { Area }                  from '../../../../model/time/area';
import { AreaService }           from '../../../../service/time/area.service';
import { ResourceService }       from '../../../../service/time/resource.service';
import { ResourceSaveComponent } from './../save/save.component';

@Component({
  selector: 'time-resource-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ResourceListComponent implements OnInit {
  @ViewChild(ResourceSaveComponent)
  saveResource: ResourceSaveComponent;

  filterAreaId: number;

  resources: Resource[];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private resourceService: ResourceService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
  }

  setFilterAreaId(areaId: number) {
    this.filterAreaId = areaId;
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    let resource = new Resource();
    if (id) {
      resource.Id = id;
    }
    this.saveResource.New(resource);
  }

  delete(resource: Resource): void {
    this.resourceService.Delete(resource.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    if (this.filterAreaId == undefined) {
      this.resourceService.List().subscribe(res => {
        this.totalCount = res.length;
        this.resources = res.slice(from, to);
      })
    } else {
      let resource = new Resource();
      resource.Area = new Area(); 
      if (this.filterAreaId != undefined) {
        resource.Area.Id = this.filterAreaId;
        resource.WithSub = true;
      }
      this.resourceService.ListWithCondition(resource).subscribe(res => {
        this.totalCount = res.length;
        this.resources = res.slice(from, to);
      })
    }
  }
}
