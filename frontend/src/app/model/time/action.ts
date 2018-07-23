import { Basic } from '../base/basic';
import { Task } from './task';

export class Action extends Basic {
  Name: string;
  TaskId: number;
  UserId: number;
  StartDate: Date;
  EndDate: Date;
  Time: number;
  Task: Task;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Task = new Task();
    }
  }
}
