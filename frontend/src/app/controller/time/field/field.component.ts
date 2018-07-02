import { Component, OnInit, ViewChild } from '@angular/core';

import { Field }              from '../../../model/time/field';
import { FieldService }       from '../../../service/time/field.service';
import { SaveFieldComponent } from './save/save.component';

@Component({
  selector: 'time-field',
  templateUrl: './field.component.html',
  styleUrls: ['./field.component.css']
})
export class FieldComponent implements OnInit {
  @ViewChild(SaveFieldComponent)
  saveField: SaveFieldComponent;

  fields: Field[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private fieldService: FieldService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveField.New(id);
  }

  delete(field: Field): void {
    this.fieldService.Delete(field.Id).subscribe(res => {
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
    this.fieldService.List().subscribe(res => {
      this.totalCount = res.length;
      this.fields = res.slice(from, to);
    })
  }
}
