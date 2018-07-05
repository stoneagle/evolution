import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Entity }           from '../../../../model/time/entity';
import { EntityService  }   from '../../../../service/time/entity.service';

@Component({
  selector: 'time-save-entity',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveEntityComponent implements OnInit {
  entity: Entity = new Entity;
  modelOpened: boolean = false;

  @Input() areaMaps: Map<number, string> = new Map();
  @Output() save = new EventEmitter<boolean>();

  constructor(
    private entityService: EntityService,
  ) { }

  ngOnInit() {
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
