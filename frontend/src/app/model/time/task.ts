import { InternationalConfig as N18 } from '../../service/base/international.service';

import { Basic }    from '../base/basic';
import { Project }  from './project';
import { Resource } from './resource';


export class Task extends Basic {
  ProjectId: number;
  ResourceId: number;
  UserId: number;
  Name: string;
  Desc: string;
  StartDate: Date;
  EndDate: Date;
  Duration: number;
  Status: number;
  Project: Project;
  Resource: Resource;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Project = new Project();
      this.Resource = new Resource();
    }
  }

  StartDateReset: boolean;
  EndDateReset: boolean;
}

export class TaskSettings {
  public Status = {
    Backlog: 1,
    Todo: 2,
    Progress: 3,
    Done: 4
  };
  public StatusName = {
    1: "backlog",
    2: "todo",
    3: "progress",
    4: "done",
  };
  public StatusInfo = {
    1: N18.settings.TIME.RESOURCE.TASK.STATUS_NAME.BACKLOG,
    2: N18.settings.TIME.RESOURCE.TASK.STATUS_NAME.TODO,
    3: N18.settings.TIME.RESOURCE.TASK.STATUS_NAME.PROGRESS,
    4: N18.settings.TIME.RESOURCE.TASK.STATUS_NAME.DONE,
  }
}

