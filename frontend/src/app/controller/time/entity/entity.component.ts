import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Entity }             from '../../../model/time/entity';
import { Area }             from '../../../model/time/area';
import { AreaService }      from '../../../service/time/area.service';
import { EntityService }      from '../../../service/time/entity.service';
import { SaveEntityComponent }            from './save/save.component';

@Component({
  selector: 'time-entity',
  templateUrl: './entity.component.html',
  styleUrls: ['./entity.component.css']
})
export class EntityComponent implements OnInit {
  @ViewChild(SaveEntityComponent)
  saveEntity: SaveEntityComponent;
  @Input() currentField: number;

  entities: Entity[];
  areaMaps: Map<number, string> = new Map();
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private entityService: EntityService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    let area = new Area();
		if (this.currentField != undefined) {
      area.FieldId = this.currentField;
    }
    this.areaService.ListAreaMap(area).subscribe(res => {
      this.areaMaps = res;
    })
  }

  initCurrentField(fieldId: number) {
    this.currentField = fieldId;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    let entity = new Entity();
    if (id) {
      entity.Id = id;
    }
    this.saveEntity.New(entity);
  }

  openSaveModelWithArea(areaId: number): void {
    let entity = new Entity();
    entity.AreaId = areaId;
    this.saveEntity.New(entity);
  }

  delete(entity: Entity): void {
    this.entityService.Delete(entity.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page && this.currentField != undefined) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    if (this.currentField == undefined) {
      this.entityService.List().subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    } else {
      let entity = new Entity();
      entity.Area = new Area(); 
      entity.Area.FieldId = this.currentField;
      this.entityService.ListWithCondition(entity).subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    }
  }
}
