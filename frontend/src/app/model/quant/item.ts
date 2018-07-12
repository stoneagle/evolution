import { Serializable } from '../base/serializable';
import { Classify } from './classify'

export class Item extends Serializable {
  Id: number;
  Code: string;
  Name: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
  Classify: Classify[];
}
