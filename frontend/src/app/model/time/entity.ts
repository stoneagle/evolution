import { Basic } from '../base/basic';
import { Area } from './area';

export class Entity extends Basic {
  Name: string;
  Desc: string;
  AreaId: number;
  Year: number;
  Area: Area;
}
