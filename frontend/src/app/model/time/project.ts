import { Basic }              from '../base/basic';
import { Quest, QuestTarget } from './quest';
import { Area }               from './area';

export class Project extends Basic {
  Name: string;
  StartDate: Date;
  Duration: number;
  QuestTarget: QuestTarget;
  Quest: Quest;
  Area: Area;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Quest = new Quest();
      this.Area = new Area();
      this.QuestTarget = new QuestTarget();
    }
  }
}
