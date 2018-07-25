import { Component, OnInit, Input, ViewChild } from '@angular/core';
import { Comparator, State, SortOrder}         from "clarity-angular";

import { PageSet }            from '../../../model/base/basic';
import { Phase }              from '../../../model/time/phase';
import { Field }              from '../../../model/time/field';

import { PhaseService }       from '../../../service/time/phase.service';
import { PhaseSaveComponent } from './save/save.component';

import { CustomComparator }                             from '../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../shared/utils';
import { PageSize }                                     from '../../../shared/const';

@Component({
  selector: 'time-phase',
  templateUrl: './phase.component.html',
  styleUrls: ['./phase.component.css']
})
export class PhaseComponent implements OnInit {
  @ViewChild(PhaseSaveComponent)
  phaseSaveComponent: PhaseSaveComponent;

  @Input() currentField: Field = new Field();

  preSorted = SortOrder.Desc;
  phases: Phase[];

  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private phaseService: PhaseService,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: any): void {
    let phase = new Phase();
    phase.FieldId = this.currentField.Id;
    phase = loadPageFilterSort<Phase>(phase, state);
    this.pageSet.Current = phase.Page.Current;
    this.currentState = state;
    this.phaseService.Count(phase).subscribe(count => {
      this.pageSet.Count = count;
      this.phaseService.List(phase).subscribe(res => {
        this.phases = res;
      })
    })
  }

  saved(savedPhase: Phase): void {
    if (savedPhase.Id) {
      this.load(this.currentState);
    } else {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.phaseSaveComponent.New(id);
  }

  delete(phase: Phase): void {
    this.phaseService.Delete(phase.Id).subscribe(res => {
      let state = deleteState(this.pageSet, this.currentState, 1);
      this.load(state);
    })
  }
}
