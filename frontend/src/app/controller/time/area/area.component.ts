import { Component, OnInit, ViewChild } from '@angular/core';
import { TreeModel, TreeModelSettings, NodeEvent } from 'ng2-tree';

import { Area }            from '../../../model/time/area';
import { AreaService }     from '../../../service/time/area.service';
import { FieldService }    from '../../../service/time/field.service';
import { EntityComponent } from '../entity/entity.component';
import { AreaType }        from '../../../shared/shared.const';

@Component({
  selector: 'time-area',
  templateUrl: './area.component.html',
  styleUrls: ['./area.component.css']
})
export class AreaComponent implements OnInit {
  // TODO
  // 1. 递归删除/实体判断是否允许删除
  // 2. 转移后，节点状态变更
  // 3. 实体转移
  @ViewChild(EntityComponent)
  entityComponent: EntityComponent;

  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  // 必须设置一个active，才能避免错误
  areaFirstFieldId: number;
  areaFirstName: string;
  areaMap: Map<number, string> = new Map(); 
  currentFieldId: number;
  parents: TreeModel[];
  tree: TreeModel;

  ngOnInit() {
    this.fieldService.Map().subscribe(res => {
      this.areaMap = res;
      this.areaFirstFieldId = this.areaMap.keys().next().value;
      this.areaFirstName = this.areaMap.get(this.areaFirstFieldId);
      this.currentFieldId = this.areaFirstFieldId;
      this.areaMap.delete(this.areaFirstFieldId);
      this.refreshParent(this.currentFieldId);
      this.entityComponent.initCurrentField(this.currentFieldId);
    })
  }

  checkActive(fieldId: number): boolean {
    let ret = false;
    if (fieldId === this.currentFieldId) {
      ret = true;
    }
    return ret
  }

  refreshParent(fieldId: number) {
    this.entityComponent.setCurrentAreaId(undefined);
    this.currentFieldId = fieldId;
    this.areaService.ListParent(this.currentFieldId).subscribe(res => {
      let parents: TreeModel[] = [];
      res.forEach((one, k) => {
        let tree = this.buildTree(one);
        parents.push(tree);
      });
      this.parents = parents;
    });
  }

	buildTree(one: Area): TreeModel {
    let tree: TreeModel
    tree = {
      id: one.Id,
      value: one.Name
    };
    tree.settings = new TreeModelSettings();
    tree.settings.isCollapsedOnInit = true; 
    if (one.Type == AreaType.Leaf) {
      tree.settings.menuItems = [
        { action: 0, name: '新节点', cssClass: 'fa fa-arrow-right' },
        { action: 2, name: '重命名', cssClass: 'fa fa-arrow-right' },
        // { action: 3, name: '删除', cssClass: 'fa fa-arrow-right' },
        { action: 4, name: '新增实体', cssClass:'fa fa-arrow-right'}
      ];
    } else {
      tree.settings.menuItems = [
        { action: 0, name: '新节点', cssClass: 'fa fa-arrow-right' },
        { action: 2, name: '重命名', cssClass: 'fa fa-arrow-right' },
        // { action: 3, name: '删除', cssClass: 'fa fa-arrow-right' },
      ];
      tree.loadChildren = (callback) => {
        setTimeout(() => {
          this.areaService.ListChildren(one.Id).subscribe(res => {
            let childrens: TreeModel[] = [];
            res.forEach((child, k) => {
              let childTree = this.buildTree(child);
              childrens.push(childTree);
            });
            callback(childrens);
          })
        }, 200);
      }
    }
    return tree;
  }

  handleRemoved(e: NodeEvent): boolean {
    if (e["node"]["children"] != null )  {
      return false;
    }
    if (e["lastIndex"] != 0) {
      let area = new Area();
      area.Id = +e.node.id;
      this.areaService.Delete(area.Id).subscribe(res => {
        // this.refreshParent(this.currentFieldId);
        return true;
      })
    }
  }

  handleSelected(e: NodeEvent): void {
    this.entityComponent.setCurrentAreaId(+e["node"]["node"]["id"]);
    this.entityComponent.refreshClassify(0, this.entityComponent.pageSize);
  }

  handleRenamed(e: NodeEvent): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Id = +e.node.id;
    this.areaService.Update(area).subscribe(res => {
      // this.refreshParent(this.currentFieldId);
    })
  }

  handleMoved(e: NodeEvent): void {
    let area = new Area();
    area.Id = +e.node.id;
    area.Parent = +e.node.parent.id;
    this.areaService.Update(area).subscribe(res => {
      // this.refreshParent(this.currentFieldId);
    })
  }

  handleCreated(e: NodeEvent): void {
    let area = new Area();
    area.Name = e.node.value;
    area.Parent = +e.node.parent.id;
    area.FieldId = this.currentFieldId;
    if (e.node.children == null) {
      e.node.removeItselfFromParent();
    } else {
      this.areaService.Add(area).subscribe(res => {
        // this.refreshParent(this.currentFieldId);
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
