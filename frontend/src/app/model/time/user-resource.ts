import { Basic }    from '../base/basic';
import { Resource } from './resource';
import { Area }     from './area';

export class UserResource extends Basic {
  UserId: number;
  ResourceId: number;
  Time: number;
  Resource: Resource;
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Resource = new Resource();
      this.Resource.Area = new Area();
    }
  }
}
