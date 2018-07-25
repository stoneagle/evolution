import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Field }           from '../../../../model/time/field';
import { FieldService  }   from '../../../../service/time/field.service';

@Component({
  selector: 'time-field-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class FieldSaveComponent implements OnInit {
  field: Field = new Field;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<Field>();

  constructor(
    private fieldService: FieldService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.fieldService.Get(id).subscribe(res => {
        this.field = res;
        this.modelOpened = true;
      })
    } else {
      this.field = new Field();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.field.Id == null) {
      this.fieldService.Add(this.field).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(this.field);
      })
    } else {
      this.fieldService.Update(this.field).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(this.field);
      })
    }
  }
}
