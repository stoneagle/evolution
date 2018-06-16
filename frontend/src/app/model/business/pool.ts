import { Serializable } from '../base/serializable.model';

export class Pool extends Serializable {
  Id: number;
  Name: string;
  Status: string;
  Strategy: number;
  Asset: number;
  Type: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
}
