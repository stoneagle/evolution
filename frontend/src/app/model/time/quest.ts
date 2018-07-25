import { InternationalConfig as N18 } from '../../service/base/international.service';

import { Basic }    from '../base/basic';
import { Resource } from './resource';
import { Area }     from './area';


export class Quest extends Basic {
  Name: string;
  StartDate: Date;
  EndDate: Date;
  FounderId: number;
  Members: number;
  Constraint: number;
  Status: number;
  Uuid: string;
  UuidNumber: number;
  StartDateReset: boolean;
  EndDateReset: boolean;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.EndDate = new Date();
    }
  }
}

export class QuestTarget extends Basic {
  QuestId: number;
  AreaId: number;
  Desc: string;
  Status: number;
  Area: Area;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Area = new Area();
    }
  }
}

export class QuestTimeTable extends Basic {
  QuestId: number;
  StartTime: string;
  EndTime: string;
  Type: number;
}

export class QuestResource extends Basic {
  QuestId: number;
  ResourceId: number;
  PhaseId: number;
  Desc: string;
  Number: number;
  Status: number;
}

export class QuestTeam extends Basic {
  QuestId: number;
  StartDate: Date;
  EndDate: Date;
  UserId: number;
}

export class QuestSettings {
  public Status = {
    Recruit: 1,
    Exec: 2,
    Finish: 3,
    Fail: 4,
  }

  public StatusInfo = {
    1: N18.settings.TIME.RESOURCE.QUEST.STATUS_NAME.RECRUIT,
    2: N18.settings.TIME.RESOURCE.QUEST.STATUS_NAME.EXEC,
    3: N18.settings.TIME.RESOURCE.QUEST.STATUS_NAME.FINISH,
    4: N18.settings.TIME.RESOURCE.QUEST.STATUS_NAME.FAIL,
  }

  public StatusFilterOptions = [
    {key: 0, value: N18.settings.SYSTEM.RESOURCE.GENERAL.ALL },
    {key: this.Status.Recruit, value:this.StatusInfo[this.Status.Recruit]},
    {key: this.Status.Exec, value:this.StatusInfo[this.Status.Exec]},
    {key: this.Status.Finish, value:this.StatusInfo[this.Status.Finish]},
    {key: this.Status.Fail, value:this.StatusInfo[this.Status.Fail]},
  ]

  public Members = {
    One: 1,
    Small: 5,
    Middle: 25,
    Large: 100,
  }

  public MembersInfo = {
    1: N18.settings.TIME.RESOURCE.QUEST.MEMBERS_NAME.ONE,
    5: N18.settings.TIME.RESOURCE.QUEST.MEMBERS_NAME.SMALL,
    25: N18.settings.TIME.RESOURCE.QUEST.MEMBERS_NAME.MIDDLE,
    100: N18.settings.TIME.RESOURCE.QUEST.MEMBERS_NAME.LARGE,
  }

  public MembersFilterOptions = [
    {key: 0, value: N18.settings.SYSTEM.RESOURCE.GENERAL.ALL },
    {key: this.Members.One, value:this.MembersInfo[this.Members.One]},
    {key: this.Members.Small, value:this.MembersInfo[this.Members.Small]},
    {key: this.Members.Middle, value:this.MembersInfo[this.Members.Middle]},
    {key: this.Members.Large, value:this.MembersInfo[this.Members.Large]},
  ]

  public Constraint = {
    ImportantAndBusy: 1,
    Important: 2,
    Busy: 3,
    Normal: 4,
  }

  public ConstraintInfo = {
    1:N18.settings.TIME.RESOURCE.QUEST.CONSTRAINT_NAME.IMPORTANT_BUSY,
    2:N18.settings.TIME.RESOURCE.QUEST.CONSTRAINT_NAME.IMPORTANT,
    3:N18.settings.TIME.RESOURCE.QUEST.CONSTRAINT_NAME.BUSY,
    4:N18.settings.TIME.RESOURCE.QUEST.CONSTRAINT_NAME.NORMAL,
  }

  public ConstraintFilterOptions = [
    {key: 0, value: N18.settings.SYSTEM.RESOURCE.GENERAL.ALL },
    {key: this.Constraint.ImportantAndBusy, value:this.ConstraintInfo[this.Constraint.ImportantAndBusy]},
    {key: this.Constraint.Important, value:this.ConstraintInfo[this.Constraint.Important]},
    {key: this.Constraint.Busy, value:this.ConstraintInfo[this.Constraint.Busy]},
    {key: this.Constraint.Normal, value:this.ConstraintInfo[this.Constraint.Normal]},
  ]

  public TimeTableType = {
    Workday: 1,
    holiday: 2,
  }

  public TimeTableTypeInfo = {
    1: N18.settings.TIME.RESOURCE.QUEST_TIMETABLE.TYPE_NAME.WORKDAY,
    2: N18.settings.TIME.RESOURCE.QUEST_TIMETABLE.TYPE_NAME.HOLIDAY,
  } 

  public TargetStatus = {
    Wait: 1,
    Finish: 2,
  }

  public TargetStatusInfo = {
    1: N18.settings.TIME.RESOURCE.QUEST_TARGET.STATUS_NAME.WAIT,
    2: N18.settings.TIME.RESOURCE.QUEST_TARGET.STATUS_NAME.FINISH,
  }

  public ResourceStatus = {
    Unmatch: 1,
    Matched: 2,
  }

  public ResourceStatusInfo = {
    1: N18.settings.TIME.RESOURCE.QUEST_RESOURCE.STATUS_NAME.UNMATCH,
    2: N18.settings.TIME.RESOURCE.QUEST_RESOURCE.STATUS_NAME.MATCHED,
  }
}

