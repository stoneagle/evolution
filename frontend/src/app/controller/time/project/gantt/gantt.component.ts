import { Component, OnInit, ViewChild } from '@angular/core';
import { EJ_GANTT_COMPONENTS }                      from 'ej-angular2/src/ej/gantt.component';

import { Quest }                from '../../../../model/time/quest';
import { QuestService }         from '../../../../service/time/quest.service';
import { ProjectService }       from '../../../../service/time/project.service';
import { SignService }          from '../../../../service/system/sign.service';
import { TranslateService }     from '@ngx-translate/core';
import { ProjectSaveComponent } from '../save/save.component';

@Component({
  selector: 'time-project-gantt',
  templateUrl: './gantt.component.html',
  styleUrls: ['./gantt.component.css']
})
export class ProjectGanttComponent implements OnInit {
  @ViewChild(ProjectSaveComponent)
  saveComponent: ProjectSaveComponent;

  constructor(
    private questService: QuestService,
    private projectService: ProjectService,
    private signService: SignService,
    private translateService: TranslateService
  ) { }

  data: any
	editSettings
  toolbarSettings: any;
  treeColumnIndex: number;

  ngOnInit() {
    let dataManager = new ej.DataManager({
      url: this.projectService.getGanttUrl(),
      crossDomain: true,
      adaptor: new ej.WebApiAdaptor(),
      headers: [{
        "Authorization": this.signService.getAuthToken(),
      }],
    });
		this.editSettings = {
				allowDeleting: true,
				// allowEditing: true,
				// allowAdding: true,
				// allowIndent: true,
				editMode: 'cellEditing'
		};
    this.data = dataManager;
    this.treeColumnIndex = 1;
    this.toolbarSettings = {
        showToolbar: true,
          toolbarItems: [
            ej.Gantt.ToolbarItems.ExpandAll,
            ej.Gantt.ToolbarItems.CollapseAll
        ]
    }
  }

  onGanttLoad():void {
		var gantt = $("#GanttPanel").ejGantt("instance");
		var columns = gantt.getColumns();
    columns.forEach((one, k) => {
      switch(one["mappingName"]) {
        case "Id":
        case "EndDate":
        case "StartDate":
        case "Progress":
        case "Duration":
          columns[k].visible = false;
          break;
        case "Name":
          columns[k].headerText = this.translateService.instant('TIME.RESOURCE.PROJECT.NAME');
          break;
      }
    });

  }

  onGanttContextMenuOpen($event): void {
    let questName = this.translateService.instant('TIME.RESOURCE.QUEST.NAME');
    let projectName = this.translateService.instant('TIME.RESOURCE.PROJECT.NAME');
    let taskName = this.translateService.instant('TIME.RESOURCE.TASK.NAME');
    let processAdd = this.translateService.instant('SYSTEM.PROCESS.CREATE');
    let processDelete = this.translateService.instant('SYSTEM.PROCESS.CREATE');
    let self = this;
    $event.contextMenuItems = [];
    let projectAddItem = {
      headerText: processAdd + projectName,
      menuId: "project-add",
      eventHandler: function(args) {
        self.saveComponent.New();
      },
    }
    let projectDeleteItem = {
      headerText: processDelete + projectName,
      menuId: "project-close",
      eventHandler: function(args) {
        this.deleteItem();
      }
    }
    switch($event.item.level) {
      case 0:
        $event.contextMenuItems.push(projectAddItem);
        $event.contextMenuItems.push(projectDeleteItem);
    }
  }

  saved($event): void {
    // let data = {Id: args.data.item.Id, Name: "updated value" };
    // this.updateRecordByTaskId(data);
    // let tempData = {
    //   Id: 5,
    //   Name: "test",
    //   Parent: 2,
    // };
    // let obj = $("#GanttPanel").ejGantt("instance");
    // let rowPosition: any;
    // rowPosition = ej.TreeGrid.RowPosition.Child;
    // obj.addRecord(tempData, rowPosition);
  }
}
