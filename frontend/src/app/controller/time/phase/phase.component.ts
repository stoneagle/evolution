import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Phase }             from '../../../model/time/phase';
import { Field }             from '../../../model/time/field';
import { PhaseService }      from '../../../service/time/phase.service';
import { SavePhaseComponent }            from './save/save.component';

@Component({
  selector: 'time-phase',
  templateUrl: './phase.component.html',
  styleUrls: ['./phase.component.css']
})
export class PhaseComponent implements OnInit {
  @ViewChild(SavePhaseComponent)
  savePhase: SavePhaseComponent;

  @Input() currentField: Field = new Field();

  phases: Phase[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private phaseService: PhaseService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.savePhase.New(id);
  }

  delete(phase: Phase): void {
    this.phaseService.Delete(phase.Id).subscribe(res => {
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
    let phase = new Phase();
    phase.FieldId = this.currentField.Id;
    this.phaseService.ListWithCondition(phase).subscribe(res => {
      this.totalCount = res.length;
      this.phases = res.slice(from, to);
    })
  }
}
