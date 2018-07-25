import { Component, OnInit, ViewChild } from '@angular/core';

import { Project }                 from '../../../model/time/project';
import { Action }                  from '../../../model/time/action';
import { Task }                  from '../../../model/time/task';
import { SyncfusionSettings }      from '../../../model/time/syncfusion';
import { QuestService }            from '../../../service/time/quest.service';
import { QuestTargetService }      from '../../../service/time/quest-target.service';
import { ProjectService }          from '../../../service/time/project.service';
import { TaskService }             from '../../../service/time/task.service';
import { ActionService }           from '../../../service/time/action.service';
import { ProjectGanttComponent }   from '../project/gantt/gantt.component';
import { TaskKanbanComponent }     from '../task/kanban/kanban.component';
import { ActionScheduleComponent } from '../action/schedule/schedule.component';

@Component({
  selector: 'time-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  @ViewChild(ProjectGanttComponent)
  projectGanttComponent: ProjectGanttComponent;
  @ViewChild(TaskKanbanComponent)
  taskKanbanComponent: TaskKanbanComponent;
  @ViewChild(ActionScheduleComponent)
  actionScheduleComponent: ActionScheduleComponent;

  scheduleHeight: string;
  scheduleViews: string[];
  scheduleCurrentView: string;
  ganttHeight: string;
  ganttLevel: string;

  constructor(
    private projectService: ProjectService,
    private taskService: TaskService,
    private syncfusionSettings: SyncfusionSettings,
    private questService: QuestService,
    private questTargetService: QuestTargetService,
  ) { 
  }

  ngOnInit() {
    let upRowHeight = "400px";
    this.scheduleHeight = upRowHeight;
    this.scheduleCurrentView = "day";
    this.scheduleViews = ["Day"];

    this.ganttHeight = upRowHeight;
    this.ganttLevel = this.syncfusionSettings.GanttLevel.Project;
  }

  ganttQuestChange(questId: number): void {
    if (questId != undefined) {
      this.taskKanbanComponent.refreshData(questId, 0);
    }
  }

  ganttProjectChange(projectId: number): void {
    if (projectId != undefined) {
      this.taskKanbanComponent.refreshData(0, projectId);
    }
  }

  ganttTaskSaved($event: Task): void {
    if ($event != undefined) {
      this.taskKanbanComponent.refreshData(0, $event.Project.Id);
      this.projectGanttComponent.projectProgressUpdate($event.ProjectId);
    }
  }

  kanbanActionSaved($event: Action): void {
    if ($event != undefined) {
      this.actionScheduleComponent.actionSaved($event);
    }
  }

  kanbanTaskDeleted($event: number): void {
    if ($event != undefined) {
      this.taskService.Get($event).subscribe(task => {
        this.projectGanttComponent.projectProgressUpdate(task.ProjectId);
      })
    }
  }

  kanbanTaskUpdated($event: Task): void {
    if ($event != undefined) {
      this.projectGanttComponent.projectProgressUpdate($event.ProjectId);
    }
  }
}
