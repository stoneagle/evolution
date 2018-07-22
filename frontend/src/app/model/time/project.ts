import { InternationalConfig as N18 } from '../../service/base/international.service';

import { Basic }              from '../base/basic';
import { Quest, QuestTarget } from './quest';
import { Area }               from './area';

export class Project extends Basic {
  Name: string;
  StartDate: Date;
  Status: number;
  QuestTargetId: number;
  QuestTarget: QuestTarget;
  Quest: Quest;
  Area: Area;
  NewFlag: boolean;
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

export class ProjectSettings {
  public Status = {
    Wait: 1,
    Finish: 2,
  };
  public StatusInfo = {
    1: N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.WAIT,
    2: N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.FINISH,
  }
}

