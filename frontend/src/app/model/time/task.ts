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

  StartDateReset: boolean;
  EndDateReset: boolean;
}
