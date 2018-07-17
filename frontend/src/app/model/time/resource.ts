import { Basic } from '../base/basic';
import { Area } from './area';

export class Resource extends Basic {
  Name: string;
  Desc: string;
  Year: number;
  Area: Area;
  WithSub: boolean;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Area = new Area();
    }
  }
}
