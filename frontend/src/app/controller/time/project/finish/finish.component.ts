import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Project, ProjectSettings }   from '../../../../model/time/project';
import { Task, TaskSettings }         from '../../../../model/time/task';
import { QuestTarget, QuestSettings } from '../../../../model/time/quest';
import { ProjectService }             from '../../../../service/time/project.service';
import { TaskService }                from '../../../../service/time/task.service';
import { QuestTargetService }         from '../../../../service/time/quest-target.service';
import { MessageHandlerService  }     from '../../../../service/base/message-handler.service';
import { ShareSettings }              from '../../../../shared/settings';
import { ErrorInfo }                  from '../../../../shared/error';

@Component({
  selector: 'time-project-finish',
  templateUrl: './finish.component.html',
  styleUrls: ['./finish.component.css']
})

export class ProjectFinishComponent implements OnInit {
  project: Project;
  questTarget: QuestTarget;
  finishTarget: boolean = false;
  modelOpened: boolean = false;
  @Output() finish = new EventEmitter<Project>();

  constructor(
    private projectService: ProjectService,
    private projectSettings: ProjectSettings,
    private taskService: TaskService,
    private taskSettings: TaskSettings,
    private questTargetService: QuestTargetService,
    private questSettings: QuestSettings,
    private messageHandlerService: MessageHandlerService,
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
  ) { }

  ngOnInit() {
    this.project = new Project();
  }

  New(id: number): void {
    this.projectService.Get(id).subscribe(project => {
      this.project = project;
      this.questTarget = this.project.QuestTarget;
      this.modelOpened = true;
    })
  }            

  Submit(): void {
    if (this.finishTarget) {
      let project = new Project();
      project.Status = this.projectSettings.Status.Wait; 
      project.QuestTargetId = this.questTarget.Id;
      this.projectService.List(project).subscribe(projects => {
        for (let k in projects) {
          let one = projects[k];
          if (one.Id != this.project.Id) {
            if (this.project.Status != this.projectSettings.Status.Finish) {
              this.messageHandlerService.showWarning(
                this.shareSettings.Time.Resource.Project,
                this.shareSettings.System.Process.Finish,
                this.errorInfo.Time.ProjectNotFinish
              );
              this.modelOpened = false;
              return;
            }
          }
        }
        this.FinishProject();
      })
    } else {
      this.FinishProject();
    }
  }

  FinishProject(): void {
    let task = new Task();
    task.ProjectId = this.project.Id;
    this.taskService.List(task).subscribe(tasks => {
      for (let k in tasks) {
        let one = tasks[k];
        if (one.Status != this.taskSettings.Status.Done) {
          this.messageHandlerService.showWarning(
            this.shareSettings.Time.Resource.Project,
            this.shareSettings.System.Process.Finish,
            this.errorInfo.Time.TaskNotFinish
          );
          this.modelOpened = false;
          return;
        }
      }
      this.project.Status = this.projectSettings.Status.Finish;
      this.projectService.Update(this.project).subscribe(res => {
        if (this.finishTarget) {
          this.questTarget.Status = this.questSettings.TargetStatus.Finish;
          this.questTargetService.Update(this.questTarget).subscribe(questTarget => {
            this.finish.emit(res);
          })
        } else {
          this.finish.emit(res);
        }
      })
      this.modelOpened = false;
    })
  }
}
