import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Area }             from '../../../../model/time/area';
import { AreaService }      from '../../../../service/time/area.service';
import { Resource }         from '../../../../model/time/resource';
import { ResourceService  } from '../../../../service/time/resource.service';
import { AreaType  }        from '../../../../shared/const';

@Component({
  selector: 'time-resource-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class ResourceSaveComponent implements OnInit {
  resource: Resource = new Resource;
  modelOpened: boolean = false;

  areaMaps: Map<number, string> = new Map();
  @Output() save = new EventEmitter<boolean>();

  constructor(
    private resourceService: ResourceService,
    private areaService: AreaService,
  ) { }

  ngOnInit() {
    let area = new Area();
    area.Type = AreaType.Leaf; 
    this.areaService.ListAreaMap(area).subscribe(res => {
      this.areaMaps = res;
    })
  }

  New(resource: Resource): void {
    if (resource.Id) {
      this.resourceService.Get(resource.Id).subscribe(res => {
        this.resource = res;
        this.modelOpened = true;
      })
    } else {
      this.resource = resource;
      this.modelOpened = true;
    }
  }            

  Submit(): void {
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
}
