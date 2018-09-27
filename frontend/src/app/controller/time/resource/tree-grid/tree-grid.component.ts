import { Component, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { EJ_TREEGRID_COMPONENTS }       from 'ej-angular2/src/ej/treegrid.component';

import { Area }            from '../../../../model/time/area';
import { TreeGrid }        from '../../../../model/time/syncfusion';
import { AreaService }     from '../../../../service/time/area.service';
import { FieldService }    from '../../../../service/time/field.service';
import { Resource }        from '../../../../model/time/resource';
import { ResourceService } from '../../../../service/time/resource.service';

@Component({
  selector: 'time-resource-tree-grid',
  templateUrl: './tree-grid.component.html',
  styleUrls: ['./tree-grid.component.css']
})
export class ResourceTreeGridComponent implements OnInit {
  constructor(
    private areaService: AreaService,
    private resourceService: ResourceService,
    private fieldService: FieldService,
  ) { }
  @Output() selectResource = new EventEmitter<Resource>();

  firstFieldId: number;
  fieldMap: Map<number, string> = new Map(); 
  currentFieldId: number;

  treeGridData: TreeGrid[];
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

  selectResourceInfo($event): void {
    if (!$event.data.IsParent) {
      let resource = new Resource();
      resource.Id = $event.data.Id
      resource.Name = $event.data.Name
      resource.Area.Id = $event.data.ParentId
      resource.Area.Name = $event.data.Parent;
      this.selectResource.emit(resource);
    }
  }

  changeField(fieldId: number): void {
    this.currentFieldId = fieldId;
    this.refresh();
  }

  refresh() {
    let resource = new Resource();
    resource.Area.FieldId = this.currentFieldId;
    this.resourceService.ListGroupByLeaf(resource).subscribe(res => {
      let tmpTreeGridData = [];
      res.forEach((area, k) => {
        let Leaf: TreeGrid = new TreeGrid();
        Leaf.Id = area.Id;
        Leaf.Name = area.Name;
        Leaf.ParentId = 0;
        Leaf.IsParent = true;
        Leaf.Parent = this.fieldMap.get(area.FieldId);
        if (area.Resources.length > 0) {
          Leaf.Children = [];
          area.Resources.forEach((one, ek) => {
            let resource: TreeGrid = new TreeGrid();
            resource.Id = one.Id;
            resource.Name = one.Name;
            resource.Parent = area.Name;
            resource.ParentId = one.Area.Id;
            resource.IsParent = false;
            Leaf.Children.push(resource);
          })
        }
        tmpTreeGridData.push(Leaf);
      });
      this.treeGridData = tmpTreeGridData;
      // this.treePageSettings["pageCount"] = this.treeGridData.length;
      // this.treePageSettings["currentPage"] = 1;
    });
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
