import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Area }           from '../../../../model/time/area';
import { AreaService }    from '../../../../service/time/area.service';
import { Entity }         from '../../../../model/time/entity';
import { EntityService  } from '../../../../service/time/entity.service';
import { AreaType  } from '../../../../shared/shared.const';

@Component({
  selector: 'time-save-entity',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveEntityComponent implements OnInit {
  entity: Entity = new Entity;
  modelOpened: boolean = false;

  areaMaps: Map<number, string> = new Map();
  @Output() save = new EventEmitter<boolean>();

  constructor(
    private entityService: EntityService,
    private areaService: AreaService,
  ) { }

  ngOnInit() {
    let area = new Area();
    area.Type = AreaType.Leaf; 
    this.areaService.ListAreaMap(area).subscribe(res => {
      this.areaMaps = res;
    })
  }

  New(entity: Entity): void {
    if (entity.Id) {
      this.entityService.Get(entity.Id).subscribe(res => {
        this.entity = res;
        this.modelOpened = true;
      })
    } else {
      this.entity = entity;
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.entity.Id == null) {
      this.entityService.Add(this.entity).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.entityService.Update(this.entity).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
