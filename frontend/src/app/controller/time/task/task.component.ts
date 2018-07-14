import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Task }                from '../../../model/time/task';
import { TaskService }         from '../../../service/time/task.service';
import { TaskKanbanComponent } from './kanban/kanban.component';

@Component({
  selector: 'time-task',
  templateUrl: './task.component.html',
  styleUrls: ['./task.component.css']
})
export class TaskComponent implements OnInit {
  @ViewChild(TaskKanbanComponent)
  taskKanbanComponent: TaskKanbanComponent;

  tasks: Task[];

  constructor(
    private taskService: TaskService,
  ) { }

  ngOnInit() {
  }
}
