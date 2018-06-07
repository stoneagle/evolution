import { Serializable } from './base/serializable.model';

export class Pool extends Serializable {
  id: number;
  name: string;
  strategy: string;
  type: number;
}
