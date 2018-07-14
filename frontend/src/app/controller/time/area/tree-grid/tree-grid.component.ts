import { Component, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { EJ_TREEGRID_COMPONENTS }       from 'ej-angular2/src/ej/treegrid.component';

import { Area }              from '../../../../model/time/area';
import { TreeGrid }          from '../../../../model/time/syncfusion';
import { AreaService }       from '../../../../service/time/area.service';
import { FieldService }      from '../../../../service/time/field.service';
import { SyncfusionService } from '../../../../service/time/syncfusion.service';
import { SignService }       from '../../../../service/system/sign.service';
import { AreaType }          from '../../../../shared/const';

@Component({
  selector: 'time-area-tree-grid',
  templateUrl: './tree-grid.component.html',
  styleUrls: ['./tree-grid.component.css']
})
export class AreaTreeGridComponent implements OnInit {
  constructor(
    private areaService: AreaService,
    private syncfusionService: SyncfusionService,
    private signService: SignService,
    private fieldService: FieldService,
  ) { }
  @Output() selectArea = new EventEmitter<Area>();

  firstFieldId: number;
  fieldMap: Map<number, string> = new Map(); 
  currentFieldId: number;

  treeGridData: any;
  treePageSettings: any;
  treeColumns: any[] = [];
  contextMenuSettings: any

  treeColumnIndex: number;

  ngOnInit() {
    this.treeColumns = [
      { field: "Id", headerText: "ID", width: "45", visible: false },
      { field: "Name", headerText: "Name" },
    ];

		this.treeColumnIndex = 1;

    this.fieldService.Map().subscribe(res => {
      this.fieldMap = res;
      this.firstFieldId = this.fieldMap.keys().next().value;
      this.currentFieldId = this.firstFieldId;
    })
  }

  selectAreaInfo($event): void {
    if ((!$event.data.IsParent) && ($event.data.Id != 0)) {
      let area     = new Area();
      area.Id      = $event.data.Id
      area.Name    = $event.data.Name
      area.FieldId = this.currentFieldId
      area.Type    = AreaType.Leaf
      this.selectArea.emit(area);
    }
  }

  changeField(fieldId: number): void {
    this.currentFieldId = fieldId;
    this.refresh();
  }

  refresh() {
    let dataManager = this.syncfusionService.GetTreeGridManager(this.currentFieldId);
    this.treeGridData = dataManager;
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
