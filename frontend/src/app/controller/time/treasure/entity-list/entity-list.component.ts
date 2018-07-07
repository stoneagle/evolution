import { Component, OnInit, ViewChild, Input } from '@angular/core';

import { Entity }          from '../../../../model/time/entity';
import { Area }            from '../../../../model/time/area';
import { Phase }           from '../../../../model/time/phase';
import { Treasure }        from '../../../../model/time/treasure';
import { AreaService }     from '../../../../service/time/area.service';
import { EntityService }   from '../../../../service/time/entity.service';
import { PhaseService }    from '../../../../service/time/phase.service';
import { TreasureService } from '../../../../service/time/treasure.service';
import { SignService }     from '../../../../service/system/sign.service';

@Component({
  selector: 'time-treasure-entity',
  templateUrl: './entity-list.component.html',
  styleUrls: ['./entity-list.component.css']
})
export class TreasureEntityComponent implements OnInit {
  @Input() currentField: number;
  @Input() areaMaps: Map<number, string>;

  entities: Entity[];
  treasureMaps: Map<number, Treasure> = new Map();
  phaseMaps: Map<number, Phase> = new Map();
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private entityService: EntityService,
    private phaseService: PhaseService,
    private treasureService: TreasureService,
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

      this.refreshTreasure();
    })
  }

  refreshTreasure(): void {
    // 更新treasureMap
    let treasure = new Treasure();
    treasure.UserId = this.signService.getCurrentUser().Id;
    this.treasureService.ListWithCondition(treasure).subscribe(res => {
      this.treasureMaps = new Map();
      res.forEach((one, k) => {
        this.treasureMaps.set(one.EntityId, one);
      })
    })
  }

  createTreasure(entity: Entity): void {
    let treasure = new Treasure();
    treasure.EntityId = entity.Id;
    treasure.Status = 1;
    treasure.PhaseId = this.phaseMaps.keys().next().value;
    treasure.SumTime = 0;
    treasure.UserId = this.signService.getCurrentUser().Id;
    this.treasureService.Add(treasure).subscribe(res => {
      this.refreshTreasure();
    })
  }

  execTreasure(entity: Entity): void {
  }

  deleteTreasure(entity: Entity): void {
    let treasure = this.treasureMaps.get(entity.Id);
    this.treasureService.Delete(treasure.Id).subscribe(res => {
      this.refreshTreasure();
    })
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
