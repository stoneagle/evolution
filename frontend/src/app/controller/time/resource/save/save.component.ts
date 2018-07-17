import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Area, AreaSettings } from '../../../../model/time/area';
import { AreaService }        from '../../../../service/time/area.service';
import { Resource }           from '../../../../model/time/resource';
import { ResourceService  }   from '../../../../service/time/resource.service';

@Component({
  selector: 'time-resource-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class ResourceSaveComponent implements OnInit {
  resource: Resource;
  modelOpened: boolean = false;

  areaIds: number[];
  areas: Area[] = [];
  @Output() save = new EventEmitter<boolean>();

  constructor(
    private resourceService: ResourceService,
    private areaService: AreaService,
  ) { }

  ngOnInit() {
    this.resource = new Resource();
    let area = new Area();
    area.Type = AreaSettings.Type.Leaf; 
    this.areaService.ListWithCondition(area).subscribe(res => {
      this.areas = res;
    })
  }

  New(areaId:number, id?: number): void {
    if (id) {
      this.resourceService.Get(id).subscribe(res => {
        this.resourceService.ListAreas(id).subscribe(areas => {
          let tmpAreaIds = []; 
          areas.forEach((one, k) => {
            tmpAreaIds.push(one.Id)
          })
          this.areaIds = tmpAreaIds;
          this.resource = res;
          this.modelOpened = true;
        })
      })
    } else {
      this.resource = new Resource();
      this.areaIds = [areaId];
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    this.resource.Area.Ids = this.areaIds;
    if (this.resource.Id == null) {
      this.resourceService.Add(this.resource).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.resourceService.Update(this.resource).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
