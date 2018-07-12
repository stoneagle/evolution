import { Serializable } from '../base/serializable';

export class SessionUser extends Serializable {
  Id: number;
  Name: string;
  Email: string;
}
