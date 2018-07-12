import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Task }              from '../../../model/time/task';
import { TaskService }       from '../../../service/time/task.service';
import { TaskSaveComponent } from './save/save.component';

@Component({
  selector: 'time-task',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.css']
})
export class TaskComponent implements OnInit {
  @ViewChild(TaskSaveComponent)
  taskSaveComponent: TaskSaveComponent;

  tasks: Task[];

  constructor(
    private taskService: TaskService,
  ) { }

  ngOnInit() {
  }

  saved(saved: boolean): void {
    if (saved) {
    }
  }

  openSaveModel(id?: number): void {
    this.taskSaveComponent.New(id);
  }
}
