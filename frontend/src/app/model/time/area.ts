import { Basic }  from '../base/basic';
import { Entity } from './entity';

export class Area extends Basic {
  Name: string;
  Parent: number;
  FieldId: number;
  Type: number;
  Entities: Entity[];
}

export class LeafWithEntities extends Basic {
  Id: number;
  Name: string;
  Parent: string;
  ParentId: number;
  IsParent: boolean;
  Children: LeafWithEntities[];
}
