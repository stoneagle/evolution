import { Basic } from '../base/basic';
import { Quest } from './quest';

export class Project extends Basic {
  Name: string;
  QuestId: number;
  AreaId: number;
  StartDate: Date;
  Duration: number;
  Quest: Quest;
}
