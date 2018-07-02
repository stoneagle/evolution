import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { EntityLife }           from '../../../../../model/time/entity';
import { EntityLifeService  }   from '../../../../../service/time/entity-life.service';

@Component({
  selector: 'time-save-entity-life',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveEntityLifeComponent implements OnInit {
  entityLife: EntityLife = new EntityLife;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private entityLifeService: EntityLifeService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.entityLifeService.Get(id).subscribe(res => {
        this.entityLife = res;
        this.modelOpened = true;
      })
    } else {
      this.entityLife = new EntityLife();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.entityLife.Id == null) {
      this.entityLifeService.Add(this.entityLife).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.entityLifeService.Update(this.entityLife).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
