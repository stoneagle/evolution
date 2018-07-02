import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, NodeEvent } from 'ng2-tree';

import { Area }             from '../../../model/time/area';
import { AreaService }      from '../../../service/time/area.service';
import { ConfigService }      from '../../../service/time/config.service';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  constructor(
    private areaService: AreaService,
    private configService: ConfigService,
  ) { }

  areaMap : Map<number, string> = new Map(); 
  currentField: number = 0;
  tree: TreeModel;

  ngOnInit() {
    this.configService.FieldMap().subscribe(res => {
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

  getKeys(map) {
    return Array.from(map.keys());
  }
}
