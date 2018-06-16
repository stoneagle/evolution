import { Serializable } from '../base/serializable.model';

export class Classify extends Serializable {
  Id: number;
  Asset: number;
  AssetString: string;
  Type: string;
  Main: string;
  Sub: string;
  Name: string;
  Tag: string;
  CreatedAt: string;
  UpdatedAt: string;
}
