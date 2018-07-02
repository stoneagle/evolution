import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, NodeEvent } from 'ng2-tree';

import { Area }             from '../../../model/time/area';
import { AreaService }      from '../../../service/time/area.service';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  // TODO tab单个动态记载+active的初始化异常
  
  constructor(
    private areaService: AreaService,
  ) { }

  treeMap : Map<number, TreeModel> = new Map(); 
  tree: TreeModel;

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.areaService.List().subscribe(res => {
      this.treeMap = res;
    })
  }

  handleRemoved(e: NodeEvent): void {
    if (e["lastIndex"] != 0) {
      let area = new Area();
      area.Id = +e.node.id;
      this.areaService.Delete(area.Id).subscribe(res => {
        this.refresh();
      })
    }
  }

  handleRenamed(e: NodeEvent): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Id = +e.node.id;
    this.areaService.Update(area).subscribe(res => {
      this.refresh();
    })
  }

  handleMoved(e: NodeEvent): void {
    let area = new Area();
    area.Id = +e.node.id;
    area.Del = 0;
    area.Parent = +e.node.parent.id;
    this.areaService.Update(area).subscribe(res => {
      this.refresh();
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
        this.refresh();
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
