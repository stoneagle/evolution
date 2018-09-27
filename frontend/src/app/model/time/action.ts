import { Basic }    from '../base/basic';
import { Task }     from './task';
import { Area }     from './area';
import { Resource } from './resource';

export class Action extends Basic {
  Name: string;
  TaskId: number;
  UserId: number;
  StartDate: Date;
  EndDate: Date;
  Time: number;
  Task: Task;
  Area: Area;
  Resource: Resource;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Task = new Task();
    }
  }
}
