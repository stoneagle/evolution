import { Serializable } from '../base/serializable.model';

export class SessionUser extends Serializable {
  Id: number;
  Name: string;
  Email: string;
}
