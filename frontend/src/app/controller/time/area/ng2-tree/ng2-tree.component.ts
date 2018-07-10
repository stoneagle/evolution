import { Component, OnInit, Input, Output, ViewChild, EventEmitter  } from '@angular/core';
import { TreeModel, TreeModelSettings, NodeEvent }                    from 'ng2-tree';

import { Area }            from '../../../../model/time/area';
import { AreaService }     from '../../../../service/time/area.service';
import { Entity }          from '../../../../model/time/entity';
import { FieldService }    from '../../../../service/time/field.service';
import { AreaType }        from '../../../../shared/shared.const';

import { SaveEntityComponent }            from '../../entity/save/save.component';

@Component({
  selector: 'time-area-ng2-tree',
  templateUrl: './ng2-tree.component.html',
  styleUrls: ['./ng2-tree.component.css']
})
export class AreaNg2TreeComponent implements OnInit {
  @ViewChild(SaveEntityComponent)
  saveEntity: SaveEntityComponent;

  // TODO
  // 1. 递归删除/实体判断是否允许删除
  // 2. 转移后，节点状态变更
  // 3. 实体转移
  constructor(
    private areaService: AreaService,
    private fieldService: FieldService,
  ) { }

  // 必须设置一个active，才能避免错误
  areaFirstFieldId: number;
  
  fieldMap: Map<number, string> = new Map(); 
  currentFieldId: number;
  parents: TreeModel[];
  tree: TreeModel;

  @Output() selectAreaId = new EventEmitter<number>();

  ngOnInit() {
    this.fieldService.Map().subscribe(res => {
      this.fieldMap = res;
      this.areaFirstFieldId = this.fieldMap.keys().next().value;
      this.refreshParent(this.areaFirstFieldId);
    })
  }

  refreshParent(fieldId: number) {
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
    this.selectAreaId.emit(+e["node"]["node"]["id"]);
    return
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
    let entity = new Entity();
    entity.AreaId = +e.node.id;
    this.saveEntity.New(entity);
  }

  saved(saved: boolean): void {
    if (saved) {
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
