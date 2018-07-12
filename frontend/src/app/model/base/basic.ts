import { Serializable } from '../base/serializable';

export class Basic extends Serializable {
  Id: number;
  Ids: number[];
  UpdatedAt: Date;
  CreatedAt: Date;
  DeletedAt: Date;
}
