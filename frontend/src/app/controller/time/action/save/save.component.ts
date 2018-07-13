import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 
import { EJ_DATETIMEPICKER_COMPONENTS } from 'ej-angular2/src/ej/datetimepicker.component';
import { EJ_SCHEDULE_COMPONENTS } from 'ej-angular2/src/ej/schedule.component';

import { QuestTeam }         from '../../../../model/time/quest';
import { Task }              from '../../../../model/time/task';
import { Action }            from '../../../../model/time/action';
import { Project }           from '../../../../model/time/project';
import { SessionUser }       from '../../../../model/base/sign';
import { ActionService  }    from '../../../../service/time/action.service';
import { TaskService  }      from '../../../../service/time/task.service';
import { ProjectService  }   from '../../../../service/time/project.service';
import { SignService  }      from '../../../../service/system/sign.service';

@Component({
  selector: 'time-action-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class ActionSaveComponent implements OnInit {
  action: Action = new Action;
  modelOpened: boolean = false;
  taskMaps: Map<number, Task> = new Map();
	currentUser: SessionUser = new SessionUser();

  @Output() save = new EventEmitter<Action>();

  constructor(
    private taskService: TaskService,
    private projectService: ProjectService,
    private signService: SignService,
    private actionService: ActionService,
  ) { }

  ngOnInit() {
    this.action.Task = new Task();
    this.action.StartDate = new Date();
    this.action.EndDate = new Date();
    this.currentUser = this.signService.getCurrentUser();
    if (this.currentUser.Id == null) {
      this.signService.current().subscribe(user => {
        this.currentUser = user;
      });
    }
  }

  NewWithDate(startDate: Date, endDate: Date): void {
    this.taskService.ListByUser(this.currentUser.Id).subscribe(tasks => {
      tasks.forEach((one, k) => {
        this.taskMaps.set(one.Id, one)
      })

      this.action = new Action();
      this.action.StartDate = startDate;
      this.action.EndDate = endDate;
      this.action.UserId = this.currentUser.Id;
      this.modelOpened = true;
    })
  }

  New(id?: number): void {
    this.taskService.ListByUser(this.currentUser.Id).subscribe(tasks => {
      tasks.forEach((one, k) => {
        this.taskMaps.set(one.Id, one)
      })

      if (id) {
        this.actionService.Get(id).subscribe(res => {
          this.action = res;
          this.action.UserId = this.currentUser.Id;
          this.action.StartDate = new Date(this.action.StartDate);
          this.action.EndDate = new Date(this.action.EndDate);
          this.modelOpened = true;
        })
      } else {
        this.action = new Action();
        this.action.StartDate = new Date();
        this.action.EndDate = new Date();
        this.action.UserId = this.currentUser.Id;
        this.modelOpened = true;
      }
    })
  }            

  Submit(): void {
    let time = this.action.EndDate.getTime() - this.action.StartDate.getTime();
    this.action.Time = time / 1000 / 60;
    if (this.action.Id == null) {
      this.actionService.Add(this.action).subscribe(res => {
        this.action.Id = res.Id;
        this.save.emit(this.action);
        this.modelOpened = false;
      })
    } else {
      this.actionService.Update(this.action).subscribe(res => {
        this.save.emit(this.action);
        this.modelOpened = false;
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
