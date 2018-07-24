import { Component, OnInit, Input, ViewChild, Inject, forwardRef } from '@angular/core';

import { Task }          from '../../../../model/time/task';
import { Resource }      from '../../../../model/time/resource';
import { Area }          from '../../../../model/time/area';
import { Action }        from '../../../../model/time/action';
import { TaskService }   from '../../../../service/time/task.service';
import { ActionService } from '../../../../service/time/action.service';

import { ShellComponent } from '../../../../base/shell/shell.component';

@Component({
  selector: 'time-action-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ActionListComponent implements OnInit {
  actions: Action[] = [];
  filterAction: Action;
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private taskService: TaskService,
    private actionService: ActionService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
  }

  NewWithFilter(filterAction: Action) {
    this.filterAction = filterAction;
    this.refreshClassify(0, 10);
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
    if (this.filterAction != undefined) {
      this.actionService.List(this.filterAction).subscribe(res => {
        this.totalCount = res.length;
        this.actions = res.slice(from, to);
      })
    }
  }
}
