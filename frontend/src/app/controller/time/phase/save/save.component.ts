import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Phase }           from '../../../../model/time/phase';
import { PhaseService  }   from '../../../../service/time/phase.service';

@Component({
  selector: 'time-save-phase',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SavePhaseComponent implements OnInit {
  phase: Phase = new Phase;
  modelOpened: boolean = false;

  @Input() currentList: Phase[] = [];
  @Input() currentFieldId: number;
  @Output() save = new EventEmitter<boolean>();

  constructor(
    private phaseService: PhaseService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.phaseService.Get(id).subscribe(res => {
        this.phase = res;
        this.modelOpened = true;
      })
    } else {
      this.phase = new Phase();
      this.phase.Level = this.currentList.length +1;
      this.phase.FieldId = this.currentFieldId;
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.phase.Id == null) {
      this.phaseService.Add(this.phase).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.phaseService.Update(this.phase).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
