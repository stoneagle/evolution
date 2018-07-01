import { Serializable } from '../base/serializable.model';

export class Project extends Serializable {
  id: number;
  name: string;
  status: string;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
}
