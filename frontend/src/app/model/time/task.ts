import { Basic }   from '../base/basic';
import { Project } from './project';

export class Task extends Basic {
  ProjectId: number;
  ResourceId: number;
  UserId: number;
  Name: string;
  StartDate: Date;
  EndDate: Date;
  Duration: number;
  Status: number;
  Project: Project;
}
