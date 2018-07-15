import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Task }        from '../../../../model/time/task';
import { Resource }    from '../../../../model/time/resource';
import { SessionUser } from '../../../../model/base/sign';
import { TaskService } from '../../../../service/time/task.service';
import { SignService } from '../../../../service/system/sign.service';

@Component({
  selector: 'time-task-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class TaskListComponent implements OnInit {

  tasks: Task[];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;
  filterResource: Resource = new Resource();
  currentUser: SessionUser = new SessionUser();

  constructor(
    private signService: SignService,
    private taskService: TaskService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.currentUser = this.signService.getCurrentUser();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  changeFilterResource(resource: Resource): void {
    this.filterResource = resource;
    this.refresh();
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
    let task = new Task();
    task.UserId = this.currentUser.Id
    task.ResourceId = this.filterResource.Id
    this.taskService.ListWithCondition(task).subscribe(res => {
      this.totalCount = res.length;
      this.tasks = res.slice(from, to);
    })
  }
}
