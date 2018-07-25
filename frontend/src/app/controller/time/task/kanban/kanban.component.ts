import { Component, OnInit, Input, Output } from '@angular/core';
import { ViewChild, EventEmitter }          from '@angular/core';
import { EJ_KANBAN_COMPONENTS }             from 'ej-angular2/src/ej/kanban.component';

import { Task, TaskSettings }         from '../../../../model/time/task';
import { Action }                     from '../../../../model/time/action';
import { Field }                      from '../../../../model/time/field';
import { Kanban }                     from '../../../../model/time/syncfusion';
import { TaskService }                from '../../../../service/time/task.service';
import { SyncfusionService }          from '../../../../service/time/syncfusion.service';
import { ActionService }              from '../../../../service/time/action.service';
import { FieldService }               from '../../../../service/time/field.service';
import { MessageHandlerService  }     from '../../../../service/base/message-handler.service';
import { ShareSettings }              from '../../../../shared/settings';
import { ErrorInfo }                  from '../../../../shared/error';
import { InternationalConfig as N18 } from '../../../../service/base/international.service';

import { TaskSaveComponent }   from '../save/save.component';
import { ActionSaveComponent } from '../../action/save/save.component';

@Component({
  selector: 'time-task-kanban',
  templateUrl: './kanban.component.html',
  styleUrls: ['./kanban.component.css']
})
export class TaskKanbanComponent implements OnInit {
  @ViewChild(TaskSaveComponent)
  taskSaveComponent: TaskSaveComponent;
  @ViewChild(ActionSaveComponent)
  actionSaveComponent: ActionSaveComponent;
  @Input() kanbanConstraint: number = 5;
  @Input() kanbanSwimCollapseFlag: boolean = true;
  @Output() taskActionSaved = new EventEmitter<Action>();
  @Output() taskDeleted = new EventEmitter<number>();
  @Output() taskUpdated = new EventEmitter<Task>();

  kanbanId: string = "TaskKanban";
  // kanbanData: any;
  kanbanData: Kanban[] = [];
  kanbanColumns: any;
  kanbanFields: any;
  kanbanWorkflow: any;
  kanbanSwimlaneSettings: any;
  kanbanCardSettings: any;
  kanbanFilterSettings: any;
  kanbanMenuItem: any;
  kanbanMenuItems: any = [];
  kanbanCustomMenuItems: any;
  updateBeforeTag: string;
  kanbanInitFlagCount: number = 0;
  fieldsMap: Map<number, Field> = new Map();

  constructor(
    private taskService: TaskService,
    private fieldService: FieldService,
    private syncfusionService: SyncfusionService,
    private actionService: ActionService,
    private taskSettings: TaskSettings,
    private messageHandlerService: MessageHandlerService,
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
  ) { }

  refreshData(filterQuestId: number, filterProjectId: number): void {
    this.syncfusionService.ListKanban(filterQuestId, filterProjectId).subscribe(res => {
      this.kanbanData = res;
    });
  }

  ngOnInit() {
    this.kanbanFields = {
      primaryKey: "Id",
      content: "Desc",
      tag: "Tags",
      title: "Name",
      color: "FieldId",
      swimlaneKey: "ProjectName",
      collapsibleCards: { 
        // field: "StatusName", 
        // key: 
        //   // array will conflict
        //   this.taskSettings.StatusName[this.taskSettings.Status.Done],
        field: "Closed", 
        key: "true" 
      },
    }
    this.fieldService.List(null).subscribe(res => {
      let colorMaps = new Object;
      // this.kanbanFilterSettings = [];
      res.forEach((one, k) => {
        this.fieldsMap.set(one.Id, one);
        colorMaps[one.Color] = one.Id.toString();
        // let filterItem = {
        //   text: one.Name,
        //   query: new ej.Query().where("FieldId", "equal", one.Id),
        // }; 
        // this.kanbanFilterSettings.push(filterItem)
      })
      this.kanbanCardSettings = {
        colorMapping: colorMaps,
      }
    });
    this.kanbanSwimlaneSettings = {
      allowDragAndDrop: false,
      unassignedGroup: {
        enable: true
      }
    }
    this.kanbanWorkflow =  [
      { 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Backlog], 
        allowedTransitions: this.taskSettings.StatusName[this.taskSettings.Status.Todo] + "," + this.taskSettings.StatusName[this.taskSettings.Status.Progress]
      },
      { 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Todo], 
        allowedTransitions: this.taskSettings.StatusName[this.taskSettings.Status.Backlog] + "," + this.taskSettings.StatusName[this.taskSettings.Status.Progress]
      },
      { 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Progress], 
        allowedTransitions: this.taskSettings.StatusName[this.taskSettings.Status.Todo] + "," + this.taskSettings.StatusName[this.taskSettings.Status.Done]
      },
    ]

    let taskName = N18.settings.TIME.RESOURCE.TASK.CONCEPT;
    let actionName = N18.settings.TIME.RESOURCE.ACTION.CONCEPT;
    this.kanbanColumns = [
      { 
        headerText: this.taskSettings.StatusInfo[this.taskSettings.Status.Backlog], 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Backlog],
        totalCount: { text: taskName },
        // isCollapsed: true,
      },
      { 
        headerText: this.taskSettings.StatusInfo[this.taskSettings.Status.Todo], 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Todo],
        constraints: { 
          max: this.kanbanConstraint,
        }, 
        totalCount: { text: taskName },
      },
      { 
        headerText: this.taskSettings.StatusInfo[this.taskSettings.Status.Progress], 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Progress],
        constraints: { 
          max: this.kanbanConstraint,
        }, 
        totalCount: { text: taskName },
      },
      { 
        headerText: this.taskSettings.StatusInfo[this.taskSettings.Status.Done], 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Done],
        totalCount: { text: taskName },
      }
    ]
    let processCreate = N18.settings.SYSTEM.PROCESS.CREATE;
    let processUpdate = N18.settings.SYSTEM.PROCESS.UPDATE;
    let processDelete = N18.settings.SYSTEM.PROCESS.DELETE;
    this.kanbanCustomMenuItems = [
      {
        text: processCreate + actionName,
      },
      {
        text: processUpdate + taskName,
      },
      {
        text: processDelete + taskName,
      },
    ]
  }

  onContextOpen($event): void {
    if ($event.cardData == undefined) {
      $event.cancel = true;
    }
  }

  onContextClick($event): void {
    let taskName      = N18.settings.TIME.RESOURCE.TASK.CONCEPT;
    let actionName      = N18.settings.TIME.RESOURCE.ACTION.CONCEPT;
    let processCreate = N18.settings.SYSTEM.PROCESS.CREATE;
    let processUpdate = N18.settings.SYSTEM.PROCESS.UPDATE;
    let processDelete = N18.settings.SYSTEM.PROCESS.DELETE;
    switch ($event.text) {
      case processCreate + actionName:
        this.actionSaveComponent.NewWithTask($event.cardData.Id);
        break;
      case processUpdate + taskName:
        if ($event.cardData.Status == this.taskSettings.Status.Done) {
          this.messageHandlerService.showWarning(
            this.shareSettings.Time.Resource.Task,
            this.shareSettings.System.Process.Update,
            this.errorInfo.Time.Finished
          );
        } else {
          this.taskSaveComponent.New($event.cardData.ProjectId, $event.cardData.Id);
          this.updateBeforeTag = $event.cardData.Tags;
        }
        break;
      case processDelete + taskName:
        let action = new Action();
        action.TaskId = $event.cardData.Id;
        this.actionService.List(action).subscribe(actions => {
          if (actions.length > 0) {
            this.messageHandlerService.showWarning(
              this.shareSettings.Time.Resource.Task,
              this.shareSettings.System.Process.Delete,
              this.errorInfo.Time.Execing
            );
            return;
          } else {
            this.taskService.Delete($event.cardData.Id).subscribe(res => {
              if (res) {
                this.taskDeleted.emit($event.cardData.Id);
                let kanbanObj = $("#" + this.kanbanId).data("ejKanban");
                kanbanObj.KanbanEdit.deleteCard($event.cardData.Id);
              }
            })
          }
        })
        break;
    }
  }

  disable($event): void {
    $event.cancel = true;
  }

  taskSaved($event: Task): void {
    let kanbanObj = $("#" + this.kanbanId).data("ejKanban"); 
    if (!$event.NewFlag) {
      let tagsArray = this.updateBeforeTag.split(",");
      tagsArray[2] = $event.Resource.Name;
      kanbanObj.updateCard(this.taskSettings.StatusName[$event.Status], [{
        Id: $event.Id, 
        Name: $event.Name, 
        Desc: $event.Desc, 
        Status: $event.Status,
        StatusName: this.taskSettings.StatusName[$event.Status],
        ProjectName: $event.Project.Name,
        Tags: tagsArray.join()
      }])
    } else {
      let field = this.fieldsMap.get($event.Area.FieldId);
      let tagsArray = [];
      tagsArray.push(field.Name);
      tagsArray.push($event.Area.Name);
      tagsArray.push($event.Resource.Name);
      let kanban = {
        Id           : $event.Id,
        Name         : $event.Name,
        Desc         : $event.Desc,
        Tags         : tagsArray.join(),
        Status       : $event.Status,
        StatusName   : this.taskSettings.StatusName[$event.Status],
        ProjectId    : $event.Project.Id,
        ProjectName  : $event.Project.Name,
        ResourceId   : $event.Resource.Id,
        ResourceName : $event.Resource.Name,
        Closed       : true,
      } 
      // TODO addCard not work
      kanbanObj.KanbanEdit.addCard(this.taskSettings.StatusName[$event.Status], [kanban])
    }
  }

  onTaskDrop($event): void {
    let statusNameMap = this.taskSettings.StatusName
    if ($event.data.length == 0) {
      $event.cancel = true;
      return
    }

    let cancelFlag = true;
    for (let key in statusNameMap) {
      if ((statusNameMap[key] == $event.data[0].StatusName) && (key != $event.data[0].Status)) {
        let task = new Task();
        task.Id = $event.data[0].Id;
        task.Status = +key;

        this.taskService.Update(task).subscribe(res => {
          let kanbanObj = $("#" + this.kanbanId).data("ejKanban"); 
          kanbanObj.updateCard($event.data[0].StatusName, [{
            Id: $event.data[0].Id, 
            Status: +key,
            StatusName: $event.data[0].StatusName,
            ProjectName: res.Project.Name,
          }])
          this.taskUpdated.emit(res);
        })
        cancelFlag = false;
        break;
      }
    }
    if (cancelFlag) {
      $event.cancel = true;
    }
  }

  onActionComplete($event): void {
    let kanbanObj = $("#" + this.kanbanId).data("ejKanban");
    switch ($event.requestType) {
      case "refresh" :
        // the third refresh can collpase swim
        if ((this.kanbanInitFlagCount == 2) && (this.kanbanSwimCollapseFlag)) {
          kanbanObj.KanbanSwimlane.collapseAll();
        } else {
          this.kanbanInitFlagCount++
        }
        break;
    }
  }

  actionSaved($event: Action): void {
    if ($event != undefined) {
      this.taskActionSaved.emit($event);
    }
  }
}
