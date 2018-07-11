import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Entity }              from '../../../../model/time/entity';
import { Area }                from '../../../../model/time/area';
import { AreaService }         from '../../../../service/time/area.service';
import { EntityService }       from '../../../../service/time/entity.service';
import { EntitySaveComponent } from './../save/save.component';

@Component({
  selector: 'time-entity-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class EntityListComponent implements OnInit {
  @ViewChild(EntitySaveComponent)
  saveEntity: EntitySaveComponent;

  filterAreaId: number;

  entities: Entity[];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private entityService: EntityService,
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
    let entity = new Entity();
    if (id) {
      entity.Id = id;
    }
    this.saveEntity.New(entity);
  }

  delete(entity: Entity): void {
    this.entityService.Delete(entity.Id).subscribe(res => {
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
      this.entityService.List().subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    } else {
      let entity = new Entity();
      entity.Area = new Area(); 
      if (this.filterAreaId != undefined) {
        entity.Area.Id = this.filterAreaId;
        entity.WithSub = true;
      }
      this.entityService.ListWithCondition(entity).subscribe(res => {
        this.totalCount = res.length;
        this.entities = res.slice(from, to);
      })
    }
  }
}
