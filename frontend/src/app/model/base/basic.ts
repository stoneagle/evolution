import { Serializable } from '../base/serializable.model';

export class Basic extends Serializable {
  Id: number;
  Ids: number[];
  UpdatedAt: Date;
  CreatedAt: Date;
  DeletedAt: Date;
}
