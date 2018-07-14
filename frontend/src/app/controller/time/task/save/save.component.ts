import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { QuestTeam }         from '../../../../model/time/quest';
import { Task }              from '../../../../model/time/task';
import { Project }           from '../../../../model/time/project';
import { Resource }          from '../../../../model/time/resource';
import { User }              from '../../../../model/system/user';
import { TaskService  }      from '../../../../service/time/task.service';
import { ProjectService  }   from '../../../../service/time/project.service';
import { ResourceService  }  from '../../../../service/time/resource.service';
import { QuestTeamService  } from '../../../../service/time/quest-team.service';
import { UserService  }      from '../../../../service/system/user.service';

@Component({
  selector: 'time-task-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class TaskSaveComponent implements OnInit {
  task: Task = new Task;
  modelOpened: boolean = false;
  resourceMaps: Map<number, Resource> = new Map();
  userMaps: Map<number, User> = new Map();

  @Output() save = new EventEmitter<Task>();

  constructor(
    private taskService: TaskService,
    private projectService: ProjectService,
    private resourceService: ResourceService,
    private questTeamService: QuestTeamService,
    private userService: UserService,
  ) { }

  ngOnInit() {
    this.task.Project = new Project();
  }

  New(projectId: number, id?: number): void {
    this.projectService.Get(projectId).subscribe(project => {
      let questTeam = new QuestTeam();
      questTeam.QuestId = project.QuestId;
      this.questTeamService.ListWithCondition(questTeam).subscribe(teams => {
        let user = new User();
        user.Ids = [];
        teams.forEach((one, k) => {
          user.Ids.push(one.UserId);
        });
        this.userService.ListWithCondition(user).subscribe(res => {
          this.userMaps = new Map();
          res.forEach((u, k) => {
            this.userMaps.set(u.Id, u);
          }) 
        });
      });

      let resource = new Resource();
      resource.Area.Id = project.AreaId;
      this.resourceService.ListWithCondition(resource).subscribe(resources => {
        this.resourceMaps = new Map();
        resources.forEach((one, k) => {
          this.resourceMaps.set(one.Id, one);
        })
      });

      if (id) {
        this.taskService.Get(id).subscribe(res => {
          this.task = res;
          this.task.Project = project;
          this.task.ProjectId = this.task.Project.Id;
          this.modelOpened = true;
        })
      } else {
        this.task = new Task();
        this.task.Project = project;
        this.task.ProjectId = this.task.Project.Id;
        this.modelOpened = true;
      }
    })
  }            

  Submit(): void {
    // this.task.StartDate = new Date(this.task.StartDate);
    if (this.task.Id == null) {
      this.taskService.Add(this.task).subscribe(res => {
        this.task.Resource = this.resourceMaps.get(this.task.ResourceId);
        this.save.emit(this.task);
        this.modelOpened = false;
      })
    } else {
      this.taskService.Update(this.task).subscribe(res => {
        this.task.Resource = this.resourceMaps.get(this.task.ResourceId);
        this.save.emit(this.task);
        this.modelOpened = false;
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
