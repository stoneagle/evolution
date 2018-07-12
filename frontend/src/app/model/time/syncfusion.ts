import { Basic }  from '../base/basic';

export class TreeGrid extends Basic {
  Id: number;
  Name: string;
  Parent: string;
  ParentId: number;
  IsParent: boolean;
  Children: TreeGrid[];
}
