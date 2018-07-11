import { Basic }  from '../base/basic';
import { Entity } from './entity';
import { Area }   from './area';

export class Quest extends Basic {
  Name: string;
  StartDate: Date;
  EndDate: Date;
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
  Entity: Entity;
  Area: Area;
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
  StartDate: Date;
  EndDate: Date;
  UserId: number;
}

