import { Component, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { EJ_TREEGRID_COMPONENTS }       from 'ej-angular2/src/ej/treegrid.component';

import { Area, LeafWithEntities }          from '../../../../model/time/area';
import { AreaService }   from '../../../../service/time/area.service';
import { FieldService }  from '../../../../service/time/field.service';
import { Entity }        from '../../../../model/time/entity';
import { EntityService } from '../../../../service/time/entity.service';

@Component({
  selector: 'time-entity-tree-grid',
  templateUrl: './tree-grid.component.html',
  styleUrls: ['./tree-grid.component.css']
})
export class TreeGridEntityComponent implements OnInit {
  constructor(
    private areaService: AreaService,
    private entityService: EntityService,
    private fieldService: FieldService,
  ) { }
  @Output() selectEntity = new EventEmitter<Entity>();

  firstFieldId: number;
  fieldMap: Map<number, string> = new Map(); 
  currentFieldId: number;

  treeGridData: LeafWithEntities[] = [];
  treePageSettings: any;
  treeColumns: any[] = [];
  contextMenuSettings: any

	contextMenuOpen(sender) {
		sender.contextMenuItems.push({
			headerText: "Add Target",
			menuId: "target",
			eventHandler: this.customMenuClick,
		});
	}

	customMenuClick(args) {
	}

  treeColumnIndex: number = 1;

  ngOnInit() {
    // TODO init field load failed
    // this.treeGridData.push(new LeafWithEntities({"Id":0,"Name":"init","Parent":"init", "Children":[
    //   {"Id":1,"Name":"test","Parent":"test"}
    // ]}));
    
    this.treeColumns = [
      { field: "Id", headerText: "ID", width: "45", visible: false },
      { field: "Name", headerText: "Name" },
      { field: "Parent", headerText: "Parent" },
    ];

    // TODO contextMenu not fire
    this.contextMenuSettings = {
      showContextMenu: true,
      contextMenuItems: ["add", "edit", "delete"]
    }

    // TODO page conflict with collapse
    // this.treePageSettings = {
    //   pageSizeMode: "all",
    //   pageSize: "10",
    //   currentPage: 1,
    // }
    this.fieldService.Map().subscribe(res => {
      this.fieldMap = res;
      this.firstFieldId = this.fieldMap.keys().next().value;
      this.currentFieldId = this.firstFieldId;
    })
  }

  selectEntityInfo($event): void {
    if (!$event.data.IsParent) {
      let entity = new Entity();
      entity.Id = $event.data.Id
      entity.Name = $event.data.Name
      entity.AreaId = $event.data.ParentId
      entity.Area = new Area();
      entity.Area.Name = $event.data.Parent;
      this.selectEntity.emit(entity);
    }
  }

  changeField(fieldId: number): void {
    this.currentFieldId = fieldId;
    this.refresh();
  }

  refresh() {
    let entity = new Entity();
    entity.Area = new Area();
    entity.Area.FieldId = this.currentFieldId;
    this.entityService.ListGroupByLeaf(entity).subscribe(res => {
      this.treeGridData = [];
      res.forEach((area, k) => {
        let Leaf: LeafWithEntities = new LeafWithEntities();
        Leaf.Id = area.Id;
        Leaf.Name = area.Name;
        Leaf.ParentId = 0;
        Leaf.IsParent = true;
        Leaf.Parent = this.fieldMap.get(area.FieldId);
        if (area.Entities.length > 0) {
          Leaf.Children = [];
          area.Entities.forEach((one, ek) => {
            let entity: LeafWithEntities = new LeafWithEntities();
            entity.Id = one.Id;
            entity.Name = one.Name;
            entity.Parent = area.Name;
            entity.ParentId = one.AreaId;
            entity.IsParent = false;
            Leaf.Children.push(entity);
          })
        }
        this.treeGridData.push(Leaf);
      });
      // this.treePageSettings["pageCount"] = this.treeGridData.length;
      // this.treePageSettings["currentPage"] = 1;
    });
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
