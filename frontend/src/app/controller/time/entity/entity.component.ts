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
  @Input() currentFieldId: number;
  currentAreaId: number;

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
		if (this.currentFieldId != undefined) {
      area.FieldId = this.currentFieldId;
    }
    this.areaService.ListAreaMap(area).subscribe(res => {
      this.areaMaps = res;
    })
  }

  setCurrentAreaId(areaId: number) {
    this.currentAreaId = areaId;
  }

  initCurrentField(fieldId: number) {
    this.currentFieldId = fieldId;
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
    if (state && state.page && this.currentFieldId != undefined) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    if (this.currentFieldId == undefined) {
      this.entityService.List().subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    } else {
      let entity = new Entity();
      entity.Area = new Area(); 
      if (this.currentAreaId != undefined) {
        entity.Area.Id = this.currentAreaId;
        entity.WithSub = true;
      }
      entity.Area.FieldId = this.currentFieldId;
      this.entityService.ListWithCondition(entity).subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    }
  }
}
