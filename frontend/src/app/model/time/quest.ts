import { Basic } from '../base/basic';

export class Quest extends Basic {
  Name: string;
  StartDate: string;
  EndDate: string;
  FounderId: number;
  Members: number;
  Constraint: number;
  Status: number;
}

export class QuestTarget extends Basic {
  QuestId: number;
  EntityId: number;
  Desc: string;
  Status: number;
}

export class QuestTimeTable extends Basic {
  QuestId: number;
  StartTime: string;
  EndTime: string;
  Type: number;
}

export class QuestEntity extends Basic {
  QuestId: number;
  EntityId: number;
  PhaseId: number;
  Desc: string;
  Number: number;
  Status: number;
}

export class QuestTeam extends Basic {
  QuestId: number;
  StartDate: string;
  EndDate: string;
  UserId: number;
  TimeTablePercent: number;
}

