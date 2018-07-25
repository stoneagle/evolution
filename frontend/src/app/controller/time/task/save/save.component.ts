import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { QuestTeam }          from '../../../../model/time/quest';
import { Task, TaskSettings } from '../../../../model/time/task';
import { Project }            from '../../../../model/time/project';
import { Resource }           from '../../../../model/time/resource';
import { User }               from '../../../../model/system/user';
import { TaskService  }       from '../../../../service/time/task.service';
import { ProjectService  }    from '../../../../service/time/project.service';
import { ResourceService  }   from '../../../../service/time/resource.service';
import { QuestTeamService  }  from '../../../../service/time/quest-team.service';
import { UserService  }       from '../../../../service/system/user.service';

@Component({
  selector: 'time-task-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class TaskSaveComponent implements OnInit {
  task: Task;
  modelOpened: boolean = false;
  resourceMaps: Map<number, Resource> = new Map();
  userMaps: Map<number, User> = new Map();

  @Output() save = new EventEmitter<Task>();

  constructor(
    private taskService: TaskService,
    private taskSettings: TaskSettings,
    private projectService: ProjectService,
    private resourceService: ResourceService,
    private questTeamService: QuestTeamService,
    private userService: UserService,
  ) { }

  ngOnInit() {
    this.task = new Task();
  }

  New(projectId: number, id?: number): void {
    if (id) {
      this.taskService.Get(id).subscribe(res => {
        this.task = res;
        this.modelOpened = true;
        let questTeam = new QuestTeam();
        questTeam.QuestId = this.task.QuestTarget.QuestId;
        this.questTeamService.List(questTeam).subscribe(teams => {
          let user = new User();
          user.Ids = [];
          teams.forEach((one, k) => {
            user.Ids.push(one.UserId);
          });
          this.userService.List(user).subscribe(res => {
            this.userMaps = new Map();
            res.forEach((u, k) => {
              this.userMaps.set(u.Id, u);
            }) 
          });
        });
        let resource = new Resource();
        resource.Area.Id = this.task.QuestTarget.AreaId;
        this.resourceService.List(resource).subscribe(resources => {
          this.resourceMaps = new Map();
          resources.forEach((one, k) => {
            this.resourceMaps.set(one.Id, one);
          })
        });
      })
    } else {
      this.projectService.Get(projectId).subscribe(project => {
        let questTeam = new QuestTeam();
        questTeam.QuestId = project.QuestTarget.QuestId;
        this.questTeamService.List(questTeam).subscribe(teams => {
          let user = new User();
          user.Ids = [];
          teams.forEach((one, k) => {
            user.Ids.push(one.UserId);
          });
          this.userService.List(user).subscribe(res => {
            this.userMaps = new Map();
            res.forEach((u, k) => {
              this.userMaps.set(u.Id, u);
            }) 
          });
        });
        let resource = new Resource();
        resource.Area.Id = project.QuestTarget.AreaId;
        this.resourceService.List(resource).subscribe(resources => {
          this.resourceMaps = new Map();
          resources.forEach((one, k) => {
            this.resourceMaps.set(one.Id, one);
          })
        });
        this.task = new Task();
        this.task.Project = project;
        this.task.Area = project.Area;
        this.task.ProjectId = this.task.Project.Id;
        this.modelOpened = true;
      })
    }
  }            

  Submit(): void {
    if (this.task.Id == null) {
      this.task.StartDate = new Date();
      this.task.Status = this.taskSettings.Status.Backlog;
      this.taskService.Add(this.task).subscribe(res => {
        this.task.Resource = this.resourceMaps.get(this.task.ResourceId);
        this.task.Id = res.Id;
        this.task.UuidNumber = res.UuidNumber;
        this.task.NewFlag = true;
        this.save.emit(this.task);
        this.modelOpened = false;
      })
    } else {
      this.taskService.Update(this.task).subscribe(res => {
        this.task.Resource = this.resourceMaps.get(this.task.ResourceId);
        this.task.UuidNumber = res.UuidNumber;
        this.task.NewFlag = false;
        this.save.emit(this.task);
        this.modelOpened = false;
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
