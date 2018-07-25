import { Serializable } from '../base/serializable';

export class Basic extends Serializable{
  Id: number;
  Ids: number[];
  UpdatedAt: Date;
  CreatedAt: Date;
  DeletedAt: Date;
  Page: PageSet;
  Sort: Sort;
}

export class PageSet {
  Size: number;
  Current: number;
  Count: number;
} 

export class Sort {
  By: string;
  Reverse: boolean;
}
