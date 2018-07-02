import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, TreeModelSettings, NodeEvent } from 'ng2-tree';

import { Area }             from '../../../model/time/area';
import { AreaService }      from '../../../service/time/area.service';
import { FieldService }      from '../../../service/time/field.service';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  areaMap : Map<number, string> = new Map(); 
  currentField: number = 0;
  tree: TreeModel;

  ngOnInit() {
    this.fieldService.Map().subscribe(res => {
      this.areaMap = res;
      this.areaMap.delete(0);
      this.refresh(this.currentField);
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
    this.areaService.ListOneTree(this.currentField).subscribe(res => {
      this.tree = res;
      res.settings = new TreeModelSettings();
			res.settings.menuItems = [
        { action: 0, name: '新节点', cssClass: 'fa fa-arrow-right' },
        { action: 2, name: '重命名', cssClass: 'fa fa-arrow-right' },
        { action: 3, name: '删除', cssClass: 'fa fa-arrow-right' },
        { action: 4, name: '新增实体', cssClass:'fa fa-arrow-right'}
      ];
    })
  }

  handleRemoved(e: NodeEvent): void {
    if (e["lastIndex"] != 0) {
      let area = new Area();
      area.Id = +e.node.id;
      this.areaService.Delete(area.Id).subscribe(res => {
        this.refresh(this.currentField);
      })
    }
  }

  handleRenamed(e: NodeEvent): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Id = +e.node.id;
    this.areaService.Update(area).subscribe(res => {
      this.refresh(this.currentField);
    })
  }

  handleMoved(e: NodeEvent): void {
    let area = new Area();
    area.Id = +e.node.id;
    area.Del = 0;
    area.Parent = +e.node.parent.id;
    this.areaService.Update(area).subscribe(res => {
      this.refresh(this.currentField);
    })
  }

  handleCreated(e: NodeEvent, field: number): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Parent = +e.node.parent.id;
    area.FieldId = field;
    area.Del = 0;
    if (e.node.children == null) {
      e.node.removeItselfFromParent();
    } else {
      this.areaService.Add(area).subscribe(res => {
        this.refresh(this.currentField);
      })
    }
  }

  handleCustom(e: NodeEvent, field: number): void {
    console.log(e);
    console.log(field);
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
