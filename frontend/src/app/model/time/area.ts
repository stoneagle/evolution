import { Basic }    from '../base/basic';
import { Resource } from './resource';

export class Area extends Basic {
  Name: string;
  Parent: number;
  FieldId: number;
  Type: number;
  Resources: Resource[];
}
