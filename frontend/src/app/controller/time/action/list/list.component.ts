import { Component, OnInit, Input, ViewChild, Inject, forwardRef } from '@angular/core';
import { Comparator, State, SortOrder}  from "clarity-angular";

import { PageSet }       from '../../../../model/base/basic';
import { Task }          from '../../../../model/time/task';
import { Resource }      from '../../../../model/time/resource';
import { Area }          from '../../../../model/time/area';
import { Action }        from '../../../../model/time/action';
import { TaskService }   from '../../../../service/time/task.service';
import { ActionService } from '../../../../service/time/action.service';

import { CustomComparator }                             from '../../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../../shared/utils';
import { PageSize }                                     from '../../../../shared/const';

import { ShellComponent } from '../../../../base/shell/shell.component';

@Component({
  selector: 'time-action-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ActionListComponent implements OnInit {
  actions: Action[] = [];
  filterAction: Action;

  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private taskService: TaskService,
    private actionService: ActionService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Small;
    this.pageSet.Current = 1;
  }

  NewWithFilter(filterAction: Action) {
    this.filterAction = filterAction;
    this.refresh();
  }

  load(state: any): void {
    if (this.filterAction != undefined) {
      let action = this.filterAction;
      action = loadPageFilterSort<Action>(action, state);
      this.pageSet.Current = action.Page.Current;
      this.currentState = state;
      this.actionService.Count(action).subscribe(count => {
        this.pageSet.Count = count;
        this.actionService.List(action).subscribe(res => {
          this.actions = res;
        })
      })
    }
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }
}
