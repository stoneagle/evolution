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
  listComponent: UserResourceListComponent;

  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  // 必须设置一个active，才能避免错误
  areaFirstField: number;
  areaFirstName: string;
  areaMaps: Map<number, string> = new Map(); 
  currentField: number;

  ngOnInit() {
    this.fieldService.Map().subscribe(res => {
      this.areaMaps = res;
      this.areaFirstField = this.areaMaps.keys().next().value;
      this.areaFirstName = this.areaMaps.get(this.areaFirstField);
      this.currentField = this.areaFirstField;
      this.areaMaps.delete(this.areaFirstField);
      this.refresh(this.currentField);
      this.listComponent.initCurrentField(this.currentField);
    })
  }

  checkActive(fieldId: number): boolean {
    let ret = false;
    if (fieldId === this.currentField) {
      ret = true;
    }
    return ret
  }

  refresh(fieldId: number) {
    this.currentField = fieldId;
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
