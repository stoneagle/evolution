import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Resource }            from '../../../../model/time/resource';
import { Area }                from '../../../../model/time/area';
import { Phase }               from '../../../../model/time/phase';
import { UserResource }        from '../../../../model/time/user-resource';
import { AreaService }         from '../../../../service/time/area.service';
import { ResourceService }     from '../../../../service/time/resource.service';
import { PhaseService }        from '../../../../service/time/phase.service';
import { UserResourceService } from '../../../../service/time/user-resource.service';
import { SignService }         from '../../../../service/system/sign.service';

@Component({
  selector: 'time-user-resource-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class UserResourceListComponent implements OnInit {
  @Input() currentField: number;
  @Input() areaMaps: Map<number, string>;

  resources: Resource[];
  userResourceMaps: Map<number, UserResource> = new Map();
  phaseMaps: Map<number, Phase> = new Map();
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private resourceService: ResourceService,
    private phaseService: PhaseService,
    private userResourceService: UserResourceService,
    private signService: SignService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.phaseService.List().subscribe(res => {
      this.phaseMaps = new Map();
      res.forEach((one, k) => {
        this.phaseMaps.set(one.Id, one);
      })
    })
  }

  initCurrentField(fieldId: number) {
    this.currentField = fieldId;
    this.refresh();
  }

  load(state: any): void {
    if (state && state.page && this.currentField != undefined) {
      this.refreshAll(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshAll(0, 10);
  }

  refreshAll(from: number, to: number): void {
    if (this.currentField == undefined) {
      return;
    }
    let resource = new Resource();
    resource.Area = new Area(); 
    resource.Area.FieldId = this.currentField;
    this.resourceService.ListWithCondition(resource).subscribe(res => {
      this.totalCount = res.length;
      this.resources = res.slice(from, to);

      this.refreshUserResource();
    })
  }

  refreshUserResource(): void {
    // 更新userResourceMap
    let userResource = new UserResource();
    userResource.UserId = this.signService.getCurrentUser().Id;
    this.userResourceService.ListWithCondition(userResource).subscribe(res => {
      this.userResourceMaps = new Map();
      res.forEach((one, k) => {
        this.userResourceMaps.set(one.ResourceId, one);
      })
    })
  }

  createUserResource(resource: Resource): void {
    let userResource = new UserResource();
    userResource.ResourceId = resource.Id;
    userResource.Time = 0;
    userResource.UserId = this.signService.getCurrentUser().Id;
    this.userResourceService.Add(userResource).subscribe(res => {
      this.refreshUserResource();
    })
  }

  execUserResource(resource: Resource): void {
  }

  deleteUserResource(resource: Resource): void {
    let userResource = this.userResourceMaps.get(resource.Id);
    this.userResourceService.Delete(userResource.Id).subscribe(res => {
      this.refreshUserResource();
    })
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
