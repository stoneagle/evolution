import { Component, OnInit, ViewChild } from '@angular/core';

import { EntityLife }             from '../../../../model/time/entity';
import { EntityLifeService }      from '../../../../service/time/entity-life.service';
import { SaveEntityLifeComponent }            from './save/save.component';

@Component({
  selector: 'time-entity-life',
  templateUrl: './life.component.html',
  styleUrls: ['./life.component.css']
})
export class EntityLifeComponent implements OnInit {
  @ViewChild(SaveEntityLifeComponent)
  saveEntityLife: SaveEntityLifeComponent;

  lifes: EntityLife[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private entityLifeService: EntityLifeService,
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
    this.saveEntityLife.New(id);
  }

  delete(entityLife: EntityLife): void {
    this.entityLifeService.Delete(entityLife.Id).subscribe(res => {
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
    this.entityLifeService.List().subscribe(res => {
      this.totalCount = res.length;
      this.lifes = res.slice(from, to);
    })
  }
}
