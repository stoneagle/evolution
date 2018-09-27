import { Component, OnInit, Input, ViewChild } from '@angular/core';
import { Inject, forwardRef }                  from '@angular/core';
import { Comparator, State, SortOrder}         from "clarity-angular";

import { PageSet }     from '../../../../model/base/basic';
import { Task }        from '../../../../model/time/task';
import { Resource }    from '../../../../model/time/resource';
import { SessionUser } from '../../../../model/base/sign';

import { TaskService } from '../../../../service/time/task.service';
import { SignService } from '../../../../service/system/sign.service';

import { ShellComponent } from '../../../../base/shell/shell.component';

import { CustomComparator }                             from '../../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../../shared/utils';
import { PageSize }                                     from '../../../../shared/const';

@Component({
  selector: 'time-task-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class TaskListComponent implements OnInit {

  filterResource: Resource = new Resource();
  currentUser: SessionUser = new SessionUser();
  tasks: Task[];
  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private signService: SignService,
    private taskService: TaskService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: State): void {
    let task = new Task();
    task = loadPageFilterSort<Task>(task, state);
    this.pageSet.Current = task.Page.Current;
    this.currentState = state;
    this.taskService.Count(task).subscribe(count => {
      this.pageSet.Count = count;
      this.taskService.List(task).subscribe(res => {
        this.tasks = res;
      })
    })
  }

  changeFilterResource(resource: Resource): void {
    this.filterResource = resource;
    this.refresh();
  }
}
