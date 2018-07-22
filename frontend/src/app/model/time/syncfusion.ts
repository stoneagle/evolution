import { Basic }  from '../base/basic';

export class TreeGrid extends Basic {
  Id: number;
  Name: string;
  Parent: string;
  ParentId: number;
  IsParent: boolean;
  Children: TreeGrid[];
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Children = [];
    }
  }
}

export class Gantt extends Basic {
  Id: number;
  Name: string;
  Relate: string;
  StartDate: Date;
  EndDate: Date;
  Progress: number;
  Duration: number;
  Expended: boolean;
  Parent: number;
  Status: number;
  Color: string;
  Children: Gantt[];
  constructor (json?: any) {
    if (json != undefined) {
      super(json)
    } else {
      this.Children = [];
    }
  }
}

export class Schedule extends Basic {
  Id: number;
  Name: string;
  StartDate: Date;
  EndDate: Date;
  AllDay: boolean;
  Recurrence: boolean;
}

export class Kanban extends Basic {
	Id          :number
	Name        :string
	Desc        :string
	Status      :number
	StatusName  :string
	ProjectName :number
	Tags        :string
	FieldName   :string
	ResourceId  :number
	ProjectId   :number
	FieldId     :number
}
