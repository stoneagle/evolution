import { Component, OnInit, ViewChild } from '@angular/core';
import { EJ_GANTT_COMPONENTS }          from 'ej-angular2/src/ej/gantt.component';
import { TranslateService }             from '@ngx-translate/core';

import { Quest }             from '../../../../model/time/quest';
import { Project }           from '../../../../model/time/project';
import { Gantt }             from '../../../../model/time/syncfusion';
import { QuestService }      from '../../../../service/time/quest.service';
import { ProjectService }    from '../../../../service/time/project.service';
import { AreaService }       from '../../../../service/time/area.service';
import { ResourceService }   from '../../../../service/time/resource.service';
import { SyncfusionService } from '../../../../service/time/syncfusion.service';
import { SignService }       from '../../../../service/system/sign.service';

import { ProjectSaveComponent } from '../save/save.component';
import { QuestSaveComponent }   from '../../quest/save/save.component';
import { TaskSaveComponent }    from '../../task/save/save.component';

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
  @ViewChild(TaskSaveComponent)
  taskSaveComponent: TaskSaveComponent;

  constructor(
    private questService: QuestService,
    private projectService: ProjectService,
    private areaService: AreaService,
    private resourceService: ResourceService,
    private syncfusionService: SyncfusionService,
    private signService: SignService,
    private translateService: TranslateService
  ) { }

  data: Gantt[] = [];
	editSettings
  toolbarSettings: any;
  treeColumnIndex: number;
  selectGanttData: any;

  ganttStartDate: Date;
  ganttEndDate: Date;
  ganttHeaderSettings: any;

  ngOnInit() {
    this.ganttStartDate = new Date(new Date().setMonth(new Date().getMonth() - 3));
    this.ganttEndDate = new Date(new Date().setMonth(new Date().getMonth() + 3));
    this.ganttHeaderSettings = {
        scheduleHeaderType: ej.Gantt.ScheduleHeaderType.Month,
        monthHeaderFormat: "yyyy MMM",
        weekHeaderFormat: "M/dd",
    }
    this.syncfusionService.ListGantt().subscribe(gantts => {
      this.data = gantts;
    })
		this.editSettings = {
				allowDeleting: true,
				// allowEditing: true,
				// allowAdding: true,
				// allowIndent: true,
				editMode: 'cellEditing'
		};
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
    var relateColumn = {
      field: "Relate",
      mappingName: "Relate",
      allowEditing: true,
      visible: false,
      headerText: "Relate",
      // isTemplateColumn: true,
      // template: "{{if eResourceTaskType=='resourceTask'}} <span style='padding:10px;'> {{if eOverlapped}} Yes {{else}} No {{/if}} </span> {{/if}}"
    };
    var colorColumn = {
      field: "Color",
      mappingName: "Color",
      allowEditing: false,
      visible: false,
      headerText: "Color",
    };
    columns.push(relateColumn);
    columns.push(colorColumn);
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
    let questUpdateItem = {
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
        self.selectGanttData = args.data;
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
    let taskAddItem = {
      headerText: processCreate + taskName,
      menuId: "task-add",
      eventHandler: function(args) {
        self.taskSaveComponent.New(args.data.taskId);
      },
    }
    let taskUpdateItem = {
      headerText: processUpdate + taskName,
      menuId: "task-update",
      eventHandler: function(args) {
        self.selectGanttData = args.data;
        self.taskSaveComponent.New(args.data.parentItem.taskId, args.data.taskId);
      },
    }
    let taskDeleteItem = {
      headerText: processDelete + taskName,
      menuId: "task-close",
      eventHandler: function(args) {
        this.deleteItem();
      }
    }
    switch($event.item.level) {
      case 0:
        $event.contextMenuItems.push(questUpdateItem);
        $event.contextMenuItems.push(projectAddItem);
        break;
      case 1:
        $event.contextMenuItems.push(projectUpdateItem);
        $event.contextMenuItems.push(projectDeleteItem);
        $event.contextMenuItems.push(taskAddItem);
        break;
      case 2:
        $event.contextMenuItems.push(taskUpdateItem);
        $event.contextMenuItems.push(taskDeleteItem);
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
    } else {
      this.areaService.Get($event.AreaId).subscribe(area => {
        let newRecord = new Gantt;
        newRecord.Name = $event.Name;
        this.selectGanttData.Relate = area.Name;
        this.selectGanttData.item.Relate = area.Name;
        gantt.updateRecordByIndex(this.selectGanttData.index, newRecord);
      })
    }
  }

  taskSaved($event): void {
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
      newRecord.Parent = $event.Project.Id
      newRecord.Progress = 0;
      newRecord.Duration = 0;
      let rowPosition: any;
      rowPosition = ej.TreeGrid.RowPosition.Child;
      gantt.addRecord(newRecord, rowPosition);
    } else {
      this.resourceService.Get($event.ResourceId).subscribe(resource => {
        let newRecord = new Gantt;
        newRecord.Name = $event.Name;
        this.selectGanttData.Relate = resource.Name;
        this.selectGanttData.item.Relate = resource.Name;
        gantt.updateRecordByIndex(this.selectGanttData.index, newRecord);
      })
    }
  }

  questSaved($event): void {
  }

  onGanttQueryTaskbarInfo($event): void {
    if (($event.data.Color != undefined) && ($event.data.Color != "")) {
      let color = $event.data.Color;
      $event.parentTaskbarBackground = color;
      $event.parentProgressbarBackground = color;
    }
  }

  onGanttRowSelected($event): void {
  }
}
