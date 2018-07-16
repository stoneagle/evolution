import { Serializable } from '../base/serializable';

export class Resp {
  status: number;
  code: number;
  data: any;
  desc: string;
}

export class RespObject extends Serializable {
  status: number;
  code: number;
  data: any;
  desc: string;
}
