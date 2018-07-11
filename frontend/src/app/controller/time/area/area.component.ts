import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, TreeModelSettings, NodeEvent } from 'ng2-tree';

import { Area }                from '../../../model/time/area';
import { AreaService }         from '../../../service/time/area.service';
import { ListEntityComponent } from '../entity//list/list.component';
import { AreaType }            from '../../../shared/shared.const';

import { AreaNg2TreeComponent } from './ng2-tree/ng2-tree.component';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  @ViewChild(AreaNg2TreeComponent)
  ng2TreeComponent: AreaNg2TreeComponent;

  @ViewChild(ListEntityComponent)
  entityListComponent: ListEntityComponent;

  constructor(
    private areaService: AreaService,
  ) { }

  ngOnInit() {
  }

  selectAreaNode($event): void {
    this.entityListComponent.setFilterAreaId($event);
    this.entityListComponent.refreshClassify(0, this.entityListComponent.pageSize);
  }
}
