import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Project, ProjectSettings } from '../../../../model/time/project';
import { Quest, QuestTarget }       from '../../../../model/time/quest';
import { Area }                     from '../../../../model/time/area';
import { Task }                     from '../../../../model/time/task';
import { Resource }                 from '../../../../model/time/resource';
import { TaskService  }             from '../../../../service/time/task.service';
import { ProjectService  }          from '../../../../service/time/project.service';
import { QuestService  }            from '../../../../service/time/quest.service';
import { QuestTargetService  }      from '../../../../service/time/quest-target.service';

@Component({
  selector: 'time-project-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class ProjectSaveComponent implements OnInit {
  project: Project;
  modelOpened: boolean = false;
  questTargetMaps: Map<number, QuestTarget> = new Map();
  tasks: Task[] = [];

  @Output() save = new EventEmitter<Project>();

  constructor(
    private taskService: TaskService,
    private projectService: ProjectService,
    private questService: QuestService,
    private questTargetService: QuestTargetService,
    private projectSettings: ProjectSettings,
  ) { }

  ngOnInit() {
    this.project = new Project();
  }

  New(questId: number, id?: number): void {
    this.questService.Get(questId).subscribe(quest => {
      let questTarget = new QuestTarget();
      questTarget.QuestId = questId;
      this.questTargetService.List(questTarget).subscribe(targets => {
        targets.forEach((one, k) => {
          this.questTargetMaps.set(one.Id, one);
        })
      })

      if (id) {
        let task = new Task();
        task.ProjectId = id;
        this.taskService.List(task).subscribe(tasks => {
          this.tasks = tasks;
        });
        this.projectService.Get(id).subscribe(res => {
          this.project = res;
          this.project.Quest.UuidNumber = quest.UuidNumber;
          this.modelOpened = true;
        })
      } else {
        this.project = new Project();
        this.project.Quest = quest;
        this.tasks = [];
        this.modelOpened = true;
      }
    })
  }            

  Submit(): void {
    this.project.StartDate = new Date(this.project.StartDate);
    if (this.project.Id == null) {
      this.project.Status = this.projectSettings.Status.Wait;
      this.projectService.Add(this.project).subscribe(res => {
        this.project.NewFlag = true;
        this.project.Id = res.Id;
        this.project.QuestTarget = this.questTargetMaps.get(this.project.QuestTargetId)
        this.project.UuidNumber = res.UuidNumber
        this.save.emit(this.project);
        this.modelOpened = false;
      })
    } else {
      this.projectService.Update(this.project).subscribe(res => {
        this.project.NewFlag = false;
        this.project.UuidNumber = res.UuidNumber
        this.project.QuestTarget = this.questTargetMaps.get(this.project.QuestTargetId)
        this.save.emit(this.project);
        this.modelOpened = false;
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
