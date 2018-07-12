import { Serializable } from '../base/serializable';
import { Item } from './item'

export class Pool extends Serializable {
  Id: number;
  Name: string;
  Status: string;
  Asset: string;
  AssetString: string;
  Type: string;
  Strategy: string;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
  Item: Item[];
}
