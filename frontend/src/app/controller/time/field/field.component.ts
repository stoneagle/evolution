import { Component, OnInit, ViewChild } from '@angular/core';
import { Comparator, State, SortOrder}  from "clarity-angular";

import { PageSet }      from '../../../model/base/basic';
import { Field }        from '../../../model/time/field';
import { FieldService } from '../../../service/time/field.service';

import { BaseComponent }      from '../base/base.component';
import { FieldSaveComponent } from './save/save.component';
import { PhaseComponent }     from '../phase/phase.component';

import { CustomComparator }                             from '../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../shared/utils';
import { PageSize }                                     from '../../../shared/const';

@Component({
  selector: 'time-field',
  templateUrl: './field.component.html',
  styleUrls: ['./field.component.css']
})
export class FieldComponent implements OnInit {
  @ViewChild(FieldSaveComponent)
  saveComponent: FieldSaveComponent;
  @ViewChild(PhaseComponent)
  phaseComponent: PhaseComponent;

  idComparator: Comparator<Field> = new CustomComparator<Field>("id", "number");
  createdComparator: Comparator<Field> = new CustomComparator<Field>("created_at", "date");
  fields: Field[];

  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private fieldService: FieldService,
  ) { 
  }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: State): void {
    let field = new Field();
    field = loadPageFilterSort<Field>(field, state);
    this.pageSet.Current = field.Page.Current;
    this.currentState = state;
    this.fieldService.Count(field).subscribe(count => {
      this.pageSet.Count = count;
      this.fieldService.List(field).subscribe(res => {
        this.fields = res;
      })
    })
  }

  saved(savedField: Field): void {
    if (savedField.Id) {
      this.load(this.currentState);
    } else {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveComponent.New(id);
  }

  openPhaseSaveModel(): void {
    if (this.phaseComponent != undefined) {
      this.phaseComponent.openSaveModel();
    }
  }

  delete(field: Field): void {
    this.fieldService.Delete(field.Id).subscribe(res => {
      let state = deleteState(this.pageSet, this.currentState, 1);
      this.load(state);
    })
  }
}
