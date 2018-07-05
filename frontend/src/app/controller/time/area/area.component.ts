import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, TreeModelSettings, NodeEvent } from 'ng2-tree';

import { Area }             from '../../../model/time/area';
import { AreaService }      from '../../../service/time/area.service';
import { FieldService }      from '../../../service/time/field.service';
import { EntityComponent } from '../entity/entity.component';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  @ViewChild(EntityComponent)
  entityComponent: EntityComponent;

  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  // 必须设置一个active，才能避免错误
  areaFirstField: number;
  areaFirstName: string;
  areaMap: Map<number, string> = new Map(); 
  currentField: number;
  tree: TreeModel;

  ngOnInit() {
    this.fieldService.Map().subscribe(res => {
      this.areaMap = res;
      this.areaFirstField = this.areaMap.keys().next().value;
      this.areaFirstName = this.areaMap.get(this.areaFirstField);
      this.currentField = this.areaFirstField;
      this.areaMap.delete(this.areaFirstField);
      this.refresh(this.currentField);
      this.entityComponent.initCurrentField(this.currentField);
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

  handleCreated(e: NodeEvent): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Parent = +e.node.parent.id;
    area.FieldId = this.currentField;
    area.Del = 0;
    if (e.node.children == null) {
      e.node.removeItselfFromParent();
    } else {
      this.areaService.Add(area).subscribe(res => {
        this.refresh(this.currentField);
      })
    }
  }

  handleCustom(e: NodeEvent): void {
    if (e.node.children.length == 0) {
      this.entityComponent.openSaveModelWithArea(+e.node.id);
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
