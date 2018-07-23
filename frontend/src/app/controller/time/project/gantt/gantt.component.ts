import { Component, OnInit, ViewChild, Output, Inject, forwardRef  } from '@angular/core';
import { EJ_GANTT_COMPONENTS }                                       from 'ej-angular2/src/ej/gantt.component';
import { InternationalConfig as N18 }                                from '../../../../service/base/international.service';

import { Quest, QuestTarget, QuestSettings } from '../../../../model/time/quest';
import { Project, ProjectSettings }          from '../../../../model/time/project';
import { Task, TaskSettings }                from '../../../../model/time/task';
import { Action }                            from '../../../../model/time/action';
import { Gantt, SyncfusionSettings }         from '../../../../model/time/syncfusion';
import { Field }                             from '../../../../model/time/field';
import { QuestService }                      from '../../../../service/time/quest.service';
import { QuestTargetService }                from '../../../../service/time/quest-target.service';
import { ProjectService }                    from '../../../../service/time/project.service';
import { TaskService }                       from '../../../../service/time/task.service';
import { ActionService }                     from '../../../../service/time/action.service';
import { AreaService }                       from '../../../../service/time/area.service';
import { ResourceService }                   from '../../../../service/time/resource.service';
import { SyncfusionService }                 from '../../../../service/time/syncfusion.service';
import { FieldService }                      from '../../../../service/time/field.service';
import { SignService }                       from '../../../../service/system/sign.service';
import { MessageHandlerService  }            from '../../../../service/base/message-handler.service';
import { ShareSettings }                     from '../../../../shared/settings';
import { ErrorInfo }                         from '../../../../shared/error';

import { ProjectSaveComponent }   from '../save/save.component';
import { ProjectFinishComponent } from '../finish/finish.component';
import { QuestSaveComponent }     from '../../quest/save/save.component';
import { TaskSaveComponent }      from '../../task/save/save.component';
import { ShellComponent }         from '../../../../base/shell/shell.component';

@Component({
  selector: 'time-project-gantt',
  templateUrl: './gantt.component.html',
  styleUrls: ['./gantt.component.css']
})
export class ProjectGanttComponent implements OnInit {
  @ViewChild(ProjectSaveComponent)
  projectSaveComponent: ProjectSaveComponent;
  @ViewChild(ProjectFinishComponent)
  projectFinishComponent: ProjectFinishComponent;
  @ViewChild(QuestSaveComponent)
  questSaveComponent: QuestSaveComponent;
  @ViewChild(TaskSaveComponent)
  taskSaveComponent: TaskSaveComponent;

  constructor(
    private questService: QuestService,
    private projectService: ProjectService,
    private projectSettings: ProjectSettings,
    private taskService: TaskService,
    private taskSettings: TaskSettings,
    private questSettings: QuestSettings,
    private questTargetService: QuestTargetService,
    private syncfusionSettings: SyncfusionSettings,
    private actionService: ActionService,
    private areaService: AreaService,
    private resourceService: ResourceService,
    private fieldService: FieldService,
    private syncfusionService: SyncfusionService,
    private signService: SignService,
    private messageHandlerService: MessageHandlerService,
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  data: any;
	editSettings
  toolbarSettings: any;
  treeColumnIndex: number;
  selectGanttData: any;
  fieldMaps: Map<number, Field> = new Map();

  ganttLevel: string;
  ganttStatus: string;

  ganttStartDate: Date;
  ganttEndDate: Date;
  ganttHeaderSettings: any;

  ganttInitFlag: boolean = false;
  ganttSaveFlag: boolean = false;

  ngOnInit() {
    this.fieldService.List(null).subscribe(res => {
      let fieldMaps = new Object;
      res.forEach((one, k) => {
        this.fieldMaps.set(one.Id, one);
      })
    });
    this.ganttLevel = this.syncfusionSettings.GanttLevel.Task;
    this.ganttStatus = this.syncfusionSettings.GanttStatus.Wait;
    this.data = this.syncfusionService.GetGanttManager(this.ganttLevel, this.ganttStatus);
    this.setGanttStartAndEndDate(this.ganttStatus);
    this.ganttHeaderSettings = {
        scheduleHeaderType: ej.Gantt.ScheduleHeaderType.Month,
        monthHeaderFormat: "yyyy MMM",
        weekHeaderFormat: "M/dd",
    }
	  this.editSettings = {
	    allowDeleting: true,
	    // allowEditing: true,
	    // allowAdding: true,
	    // allowIndent: true,
	    editMode: 'cellEditing'
	  };
    this.treeColumnIndex = 1;
    let projectStatusWait   = N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.WAIT;
    let projectStatusFinish = N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.FINISH;
    this.toolbarSettings    = {
      showToolbar: true,
      // toolbarItems: [
      //   ej.Gantt.ToolbarItems.ExpandAll,
      //   ej.Gantt.ToolbarItems.CollapseAll,
      // ],
      customToolbarItems: [
        {
          text: "gantt-show-finish",
          tooltipText: projectStatusFinish
        },
      ]
    }
  }

  setGanttStartAndEndDate(ganttStatus: string) {
    this.syncfusionService.ListGantt(this.syncfusionSettings.GanttLevel.Task, ganttStatus).subscribe(gantts => {
      let tmpEarlyDate: Date = new Date();
      let tmpLateDate: Date = new Date();
      gantts.forEach((one, k) => {
        let startDate = new Date(one.StartDate)
        if (startDate < tmpEarlyDate) {
          tmpEarlyDate = startDate;
        }
        let endDate = new Date(one.EndDate)
        if (endDate > tmpLateDate) {
          tmpLateDate = endDate;
        }
      })
      this.ganttStartDate = tmpEarlyDate;
      this.ganttEndDate = new Date(tmpLateDate.setMonth(tmpLateDate.getMonth() + 1));
    })
  }

  onGanttAction($event): void {
    let gantt = $("#GanttPanel").ejGantt("instance");
    switch ($event.requestType) {
      case "create" :
        if (!this.ganttInitFlag) {
          this.ganttInitFlag = true;
          gantt.collapseAllItems(); 
        }
        break;
      case "refresh" :
        if (this.ganttSaveFlag) {
          this.ganttSaveFlag = false;
        } else {
          gantt.collapseAllItems(); 
        }
        break;
    }
  }

  onGanttToolbarClick($event): void {
    let projectStatusWait   = N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.WAIT;
    let projectStatusFinish = N18.settings.TIME.RESOURCE.PROJECT.STATUS_NAME.FINISH;
    switch ($event.itemName) {
      case (projectStatusFinish):
        $($event.currentTarget).find("a").removeClass("gantt-show-finish").addClass("gantt-hide-finish");
        $($event.currentTarget).attr("data-content", projectStatusWait);
        this.ganttStatus = this.syncfusionSettings.GanttStatus.Finish;
        this.data = this.syncfusionService.GetGanttManager(this.ganttLevel, this.ganttStatus);
        break;
      case (projectStatusWait):
        $($event.currentTarget).find("a").removeClass("gantt-hide-finish").addClass("gantt-show-finish");
        $($event.currentTarget).attr("data-content", projectStatusFinish);
        this.ganttStatus = this.syncfusionSettings.GanttStatus.Wait;
        this.data = this.syncfusionService.GetGanttManager(this.ganttLevel, this.ganttStatus);
        break;
    }
    this.setGanttStartAndEndDate(this.ganttStatus);
  }

  onGanttLoad($event): void {
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
          let projectName       = N18.settings.TIME.RESOURCE.PROJECT.CONCEPT;
          columns[k].headerText = projectName;
          columns[k].width      = "300";
          break;
      }
    });
    var relateColumn = {
      field: "Relate",
      mappingName: "Relate",
      allowEditing: true,
      visible: true,
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
    var statusColumn = {
      field: "Status",
      mappingName: "Status",
      allowEditing: false,
      visible: false,
      headerText: "Status",
    };
    var relateIdColumn = {
      field: "RelateId",
      mappingName: "RelateId",
      allowEditing: false,
      visible: false,
      headerText: "RelateId",
    };
    columns.push(relateColumn);
    columns.push(colorColumn);
    columns.push(statusColumn);
    columns.push(relateIdColumn);
    // gantt data init empty will cause console alert
    try {
      gantt.sortColumn("StartDate","ascending");
    } catch (e) {
      // console.log(e);
    }
  }

  onGanttContextMenuOpen($event): void {
    let questName       = N18.settings.TIME.RESOURCE.QUEST.CONCEPT;
    let projectName     = N18.settings.TIME.RESOURCE.PROJECT.CONCEPT;
    let taskName        = N18.settings.TIME.RESOURCE.TASK.CONCEPT;
    let processCreate   = N18.settings.SYSTEM.PROCESS.CREATE;
    let processUpdate   = N18.settings.SYSTEM.PROCESS.UPDATE;
    let processDelete   = N18.settings.SYSTEM.PROCESS.DELETE;
    let processFinish   = N18.settings.SYSTEM.PROCESS.FINISH;
    let processRollback = N18.settings.SYSTEM.PROCESS.ROLLBACK;

    let self = this;
    $event.contextMenuItems = [];
    let questUpdateItem = {
      headerText: processUpdate + questName,
      menuId: "quest-update",
      eventHandler: function(args) {
        self.questSaveComponent.New(self.shell.currentUser.Id, args.data.item.RelateId);
      },
    }
    let projectAddItem = {
      headerText: processCreate + projectName,
      menuId: "project-add",
      eventHandler: function(args) {
        self.projectSaveComponent.New(args.data.item.RelateId);
      },
    }
    let projectUpdateItem = {
      headerText: processUpdate + projectName,
      menuId: "project-update",
      eventHandler: function(args) {
        self.selectGanttData = args.data;
        self.projectSaveComponent.New(args.data.parentItem.RelateId, args.data.RelateId);
      },
    }
    let projectFinishItem = {
      headerText: processFinish + projectName,
      menuId: "project-finish",
      eventHandler: function(args) {
        self.selectGanttData = args.data;
        self.projectFinishComponent.New(args.data.RelateId);
      },
    }
    let projectRollbackItem = {
      headerText: processRollback + projectName,
      menuId: "project-rollback",
      eventHandler: function(args) {
        if (args.data.Status != self.projectSettings.Status.Finish) {
          self.messageHandlerService.showWarning(
            self.shareSettings.Time.Resource.Project,
            self.shareSettings.System.Process.Rollback,
            self.errorInfo.Time.NotFinish,
          );
          return;
        }
        let project = new Project();
        project.Id = args.data.RelateId;
        project.Status = self.projectSettings.Status.Wait; 
        self.projectService.Update(project).subscribe(project => {
          let questTarget = new QuestTarget();
          questTarget.Id = project.QuestTargetId;
          questTarget.Status = self.questSettings.TargetStatus.Wait;
          self.questTargetService.Update(questTarget).subscribe(questTarget => {
            let gantt = $("#GanttPanel").ejGantt("instance");
            let newRecord = new Gantt;
            newRecord.Id = project.UuidNumber;
            newRecord.Status = project.Status;
            args.data.Status = project.Status;
            args.data.item.Status = project.Status;
            gantt.updateRecordByTaskId(newRecord);
          })
        })
      },
    }
    let projectDeleteItem = {
      headerText: processDelete + projectName,
      menuId: "project-delete",
      eventHandler: function(args) {
        let task = new Task();
        task.ProjectId = args.data.RelateId;
        self.taskService.List(task).subscribe(tasks => {
          if (tasks.length > 0) {
            self.messageHandlerService.showWarning(
              self.shareSettings.Time.Resource.Project,
              self.shareSettings.System.Process.Delete,
              self.errorInfo.Time.Execing
            );
            return;
          } else {
            self.projectService.Delete(args.data.RelateId).subscribe(res => {
              if (res) {
                this.deleteItem();
              }
            })
          }
        })
      }
    }
    let taskAddItem = {
      headerText: processCreate + taskName,
      menuId: "task-add",
      eventHandler: function(args) {
        self.taskSaveComponent.New(args.data.RelateId);
      },
    }
    let taskUpdateItem = {
      headerText: processUpdate + taskName,
      menuId: "task-update",
      eventHandler: function(args) {
        self.selectGanttData = args.data;
        self.taskSaveComponent.New(args.data.parentItem.RelateId, args.data.RelateId);
      },
    }
    let taskFinishItem = {
      headerText: processFinish + taskName,
      menuId: "task-finish",
      eventHandler: function(args) {
        if (args.data.Status != self.taskSettings.Status.Progress) {
          self.messageHandlerService.showWarning(
            self.shareSettings.Time.Resource.Task,
            self.shareSettings.System.Process.Finish,
            self.errorInfo.Time.NotExec
          );
          return;
        }
        let action = new Action();
        action.TaskId = args.data.RelateId;
        self.actionService.List(action).subscribe(actions => {
          if (actions.length == 0) {
            self.messageHandlerService.showWarning(
              self.shareSettings.Time.Resource.Task,
              self.shareSettings.System.Process.Finish,
              self.errorInfo.Time.NotExec
            );
            return;
          }
          let task = new Task();
          task.Id = args.data.RelateId;
          task.Status = self.taskSettings.Status.Done; 
          self.taskService.Update(task).subscribe(task => {
            let gantt = $("#GanttPanel").ejGantt("instance");
            let newRecord = new Gantt;
            newRecord.Id = task.UuidNumber;
            newRecord.Status = task.Status;
            newRecord.Progress = 100;
            args.data.Status = task.Status;
            args.data.item.Status = task.Status;
            args.data.Progress = 100;
            args.data.item.Progress = 100;
            gantt.updateRecordByTaskId(newRecord);
          })
        })
      },
    }
    let taskRollbackItem = {
      headerText: processRollback + taskName,
      menuId: "task-rollback",
      eventHandler: function(args) {
        if (args.data.Status != self.taskSettings.Status.Done) {
          self.messageHandlerService.showWarning(
            self.shareSettings.Time.Resource.Task,
            self.shareSettings.System.Process.Finish,
            self.errorInfo.Time.NotFinish,
          );
          return;
        }
        let task = new Task();
        task.Id = args.data.RelateId;
        task.Status = self.taskSettings.Status.Progress; 
        self.taskService.Update(task).subscribe(task => {
          let gantt = $("#GanttPanel").ejGantt("instance");
          let newRecord = new Gantt;
          newRecord.Id = task.UuidNumber;
          newRecord.Status = task.Status;
          newRecord.Progress = 0;
          args.data.Status = task.Status;
          args.data.item.Status = task.Status;
          args.data.Progress = 0;
          args.data.item.Progress = 0;
          gantt.updateRecordByTaskId(newRecord);
        })
      },
    }
    let taskDeleteItem = {
      headerText: processDelete + taskName,
      menuId: "task-delete",
      eventHandler: function(args) {
        let action = new Action();
        action.TaskId = args.data.RelateId;
        self.actionService.List(action).subscribe(actions => {
          if (actions.length > 0) {
            self.messageHandlerService.showWarning(
              self.shareSettings.Time.Resource.Task,
              self.shareSettings.System.Process.Delete,
              self.errorInfo.Time.Execing
            );
            return;
          } else {
            self.taskService.Delete(args.data.RelateId).subscribe(res => {
              if (res) {
                this.deleteItem();
              }
            })
          }
        })
      }
    }
    switch($event.item.level) {
      case 0:
        $event.contextMenuItems.push(questUpdateItem);
        $event.contextMenuItems.push(projectAddItem);
        break;
      case 1:
        if ($event.item.Status == this.projectSettings.Status.Wait) {
          $event.contextMenuItems.push(projectUpdateItem);
        }
        if (!$event.item.childRecords) {
          $event.contextMenuItems.push(projectDeleteItem);
        } else if ($event.item.status == 100 && $event.item.Status == this.projectSettings.Status.Wait) {
          $event.contextMenuItems.push(projectFinishItem);
        } else if ($event.item.Status == this.projectSettings.Status.Finish) {
          $event.contextMenuItems.push(projectRollbackItem);
        }
        $event.contextMenuItems.push(taskAddItem);
        break;
      case 2:
        if ($event.item.Status != this.taskSettings.Status.Done) {
          $event.contextMenuItems.push(taskUpdateItem);
        }
        if ($event.item.Status == this.taskSettings.Status.Backlog) {
          $event.contextMenuItems.push(taskDeleteItem);
        } else if ($event.item.Status == this.taskSettings.Status.Progress) {
          $event.contextMenuItems.push(taskFinishItem);
        } else if ($event.item.Status == this.taskSettings.Status.Done) {
          $event.contextMenuItems.push(taskRollbackItem);
        }
        break;
    }
  }

  projectSaved($event: Project): void {
    if ($event == undefined) {
      return
    }
    // need ensure index and id unique
    let gantt = $("#GanttPanel").ejGantt("instance");
    if ($event.NewFlag) {
      let newRecord = new Gantt;
      newRecord.Id = $event.UuidNumber;
      newRecord.Name = $event.Name;
      newRecord.Relate = $event.QuestTarget.Area.Name;
      newRecord.RelateId = $event.Id;
      newRecord.StartDate = $event.StartDate;
      newRecord.EndDate = null;
      newRecord.Parent = $event.Quest.UuidNumber;
      newRecord.Status = this.projectSettings.Status.Wait;
      newRecord.Progress = 0;
      newRecord.Duration = 0;
      newRecord.Color = this.fieldMaps.get($event.QuestTarget.Area.FieldId).Color;
      let rowPosition: any;
      rowPosition = ej.Gantt.RowPosition.Child;
      this.ganttSaveFlag = true;
      gantt.addRecord(newRecord, rowPosition);
    } else {
      let newRecord = new Gantt;
      newRecord.Id = $event.UuidNumber;
      newRecord.Name = $event.Name;
      newRecord.Relate = $event.QuestTarget.Area.Name;
      this.selectGanttData.Relate = $event.QuestTarget.Area.Name;
      this.selectGanttData.item.Relate = $event.QuestTarget.Area.Name;
      this.ganttSaveFlag = true;
      gantt.updateRecordByTaskId(newRecord);
    }
  }

  projectFinished($event: Project): void {
    if ($event == undefined) {
      return
    }
    let gantt = $("#GanttPanel").ejGantt("instance");
    let newRecord = new Gantt;
    newRecord.Id = $event.UuidNumber;
    newRecord.Status = $event.Status;
    this.selectGanttData.Status = $event.Status;
    this.selectGanttData.item.Status = $event.Status;
    gantt.updateRecordByTaskId(newRecord);
  }

  taskSaved($event: Task): void {
    if ($event == undefined) {
      return
    }
    let gantt = $("#GanttPanel").ejGantt("instance");
    if ($event.NewFlag) {
      let newRecord = new Gantt;
      newRecord.Id = $event.UuidNumber;
      newRecord.Name = $event.Name;
      newRecord.Relate = $event.Resource.Name;
      newRecord.RelateId = $event.Id;
      newRecord.StartDate = $event.StartDate;
      newRecord.EndDate = null;
      newRecord.Parent = $event.Project.UuidNumber;
      newRecord.Progress = 0;
      newRecord.Duration = 0;
      newRecord.Status = this.taskSettings.Status.Backlog;
      let rowPosition: any;
      rowPosition = ej.Gantt.RowPosition.Child;
      this.ganttSaveFlag = true;
      gantt.addRecord(newRecord, rowPosition);
    } else {
      let newRecord = new Gantt;
      newRecord.Id = $event.UuidNumber;
      newRecord.Name = $event.Name;
      newRecord.Relate = $event.Resource.Name;
      this.selectGanttData.Relate = $event.Resource.Name;
      this.selectGanttData.item.Relate = $event.Resource.Name;
      this.ganttSaveFlag = true;
      gantt.updateRecordByTaskId(newRecord);
    }
  }

  questSaved($event): void {
  }

  onGanttQueryTaskbarInfo($event): void {
    if (($event.data.Color != undefined) && ($event.data.Color != "")) {
      let color = $event.data.Color;
      // $event.parentTaskbarBackground = color;
      $event.parentProgressbarBackground = color;
    }
  }

  onGanttRowSelected($event): void {
  }

  disabled($event) {
  }
}
