import { Component, OnInit, ViewChild } from '@angular/core';

import { Area }                      from '../../../model/time/area';
import { AreaService }               from '../../../service/time/area.service';
import { FieldService }              from '../../../service/time/field.service';
import { UserResourceListComponent } from './list/list.component';

@Component({
  selector: 'time-user-resource',
  templateUrl: './user-resource.component.html',
  styleUrls: ['./user-resource.component.css']
})
export class UserResourceComponent implements OnInit {
  @ViewChild(UserResourceListComponent)
  userResourceListComponent: UserResourceListComponent;

  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  ngOnInit() {
  }

  selectAreaNode($event) {
    this.userResourceListComponent.changeFilterArea($event);
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
