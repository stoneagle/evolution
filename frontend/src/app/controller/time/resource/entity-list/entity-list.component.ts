import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Entity }          from '../../../../model/time/entity';
import { Area }            from '../../../../model/time/area';
import { Phase }           from '../../../../model/time/phase';
import { Resource }        from '../../../../model/time/resource';
import { AreaService }     from '../../../../service/time/area.service';
import { EntityService }   from '../../../../service/time/entity.service';
import { PhaseService }    from '../../../../service/time/phase.service';
import { ResourceService } from '../../../../service/time/resource.service';
import { SignService }     from '../../../../service/system/sign.service';

@Component({
  selector: 'time-resource-entity',
  templateUrl: './entity-list.component.html',
  styleUrls: ['./entity-list.component.css']
})
export class ResourceEntityComponent implements OnInit {
  @Input() currentField: number;
  @Input() areaMaps: Map<number, string>;

  entities: Entity[];
  resourceMaps: Map<number, Resource> = new Map();
  phaseMaps: Map<number, Phase> = new Map();
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private entityService: EntityService,
    private phaseService: PhaseService,
    private resourceService: ResourceService,
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
    let entity = new Entity();
    entity.Area = new Area(); 
    entity.Area.FieldId = this.currentField;
    this.entityService.ListWithCondition(entity).subscribe(res => {
      this.totalCount = res.length;
      this.entities = res.slice(from, to);

      this.refreshResource();
    })
  }

  refreshResource(): void {
    // 更新resourceMap
    let resource = new Resource();
    resource.UserId = this.signService.getCurrentUser().Id;
    this.resourceService.ListWithCondition(resource).subscribe(res => {
      this.resourceMaps = new Map();
      res.forEach((one, k) => {
        this.resourceMaps.set(one.EntityId, one);
      })
    })
  }

  createResource(entity: Entity): void {
    let resource = new Resource();
    resource.EntityId = entity.Id;
    resource.Status = 1;
    resource.PhaseId = this.phaseMaps.keys().next().value;
    resource.SumTime = 0;
    resource.UserId = this.signService.getCurrentUser().Id;
    this.resourceService.Add(resource).subscribe(res => {
      this.refreshResource();
    })
  }

  execResource(entity: Entity): void {
  }

  deleteResource(entity: Entity): void {
    let resource = this.resourceMaps.get(entity.Id);
    this.resourceService.Delete(resource.Id).subscribe(res => {
      this.refreshResource();
    })
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
