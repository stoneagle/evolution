import { Serializable } from '../base/serializable.model';

export class Item extends Serializable {
  Id: number;
  Code: string;
  Name: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
}
