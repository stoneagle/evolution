import { Serializable } from '../base/serializable.model';

export class Basic extends Serializable {
  Id: number;
  UpdatedAt: string;
  CreatedAt: string;
  DeletedAt: string;
}
