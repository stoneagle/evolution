import { Component, OnInit, ViewChild } from '@angular/core';
import { EJ_GANTT_COMPONENTS }                      from 'ej-angular2/src/ej/gantt.component';

import { Quest }                from '../../../../model/time/quest';
import { Project }              from '../../../../model/time/project';
import { Gantt }                from '../../../../model/time/syncfusion';
import { QuestService }         from '../../../../service/time/quest.service';
import { ProjectService }       from '../../../../service/time/project.service';
import { SignService }          from '../../../../service/system/sign.service';
import { TranslateService }     from '@ngx-translate/core';
import { ProjectSaveComponent } from '../save/save.component';
import { QuestSaveComponent }   from '../../quest/save/save.component';

@Component({
  selector: 'time-project-gantt',
  templateUrl: './gantt.component.html',
  styleUrls: ['./gantt.component.css']
})
export class ProjectGanttComponent implements OnInit {
  @ViewChild(ProjectSaveComponent)
  projectSaveComponent: ProjectSaveComponent;
  @ViewChild(QuestSaveComponent)
  questSaveComponent: QuestSaveComponent;

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
          this.translateService.get('TIME.RESOURCE.PROJECT.NAME').subscribe(result => {
            columns[k].headerText = result;
          });
          columns[k].width = "300";
          break;
      }
    });

  }

  onGanttContextMenuOpen($event): void {
    let questName = this.translateService.instant('TIME.RESOURCE.QUEST.CONCEPT');
    let projectName = this.translateService.instant('TIME.RESOURCE.PROJECT.CONCEPT');
    let taskName = this.translateService.instant('TIME.RESOURCE.TASK.CONCEPT');
    let processCreate = this.translateService.instant('SYSTEM.PROCESS.CREATE');
    let processUpdate = this.translateService.instant('SYSTEM.PROCESS.UPDATE');
    let processDelete = this.translateService.instant('SYSTEM.PROCESS.DELETE');
    let self = this;
    $event.contextMenuItems = [];
    let taskUpdateItem = {
      headerText: processUpdate + questName,
      menuId: "quest-update",
      eventHandler: function(args) {
        self.questSaveComponent.New(args.data.item.Id);
      },
    }
    let projectAddItem = {
      headerText: processCreate + projectName,
      menuId: "project-add",
      eventHandler: function(args) {
        self.projectSaveComponent.New(args.data.item.Id);
      },
    }
    let projectUpdateItem = {
      headerText: processUpdate + projectName,
      menuId: "project-update",
      eventHandler: function(args) {
        self.projectSaveComponent.New(args.data.parentItem.taskId, args.data.taskId);
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
        $event.contextMenuItems.push(taskUpdateItem);
        $event.contextMenuItems.push(projectAddItem);
        break;
      case 1:
        $event.contextMenuItems.push(projectUpdateItem);
        $event.contextMenuItems.push(projectDeleteItem);
        break;
    }
  }

  projectSaved($event): void {
    if ($event == undefined) {
      return
    }

    let gantt = $("#GanttPanel").ejGantt("instance");
    if ($event.Id == null) {
      let newRecord = new Gantt;
      newRecord.Id = $event.Id;
      newRecord.Name = $event.Name;
      newRecord.StartDate = $event.StartDate;
      newRecord.EndDate = null;
      newRecord.Parent = $event.Quest.Id
      newRecord.Progress = 0;
      newRecord.Duration = 0;
      let rowPosition: any;
      rowPosition = ej.TreeGrid.RowPosition.Child;
      gantt.addRecord(newRecord, rowPosition);
    }
    // let data = {Id: args.data.item.Id, Name: "updated value" };
    // this.updateRecordByTaskId(data);
  }

  questSaved($event): void {
  }
}
