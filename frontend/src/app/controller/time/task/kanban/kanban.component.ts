import { Component, OnInit, Input, ViewChild } from '@angular/core';
import { EJ_KANBAN_COMPONENTS }                from 'ej-angular2/src/ej/kanban.component';

import { Task, TaskSettings }         from '../../../../model/time/task';
import { Kanban }                     from '../../../../model/time/syncfusion';
import { TaskService }                from '../../../../service/time/task.service';
import { SyncfusionService }          from '../../../../service/time/syncfusion.service';
import { FieldService }               from '../../../../service/time/field.service';
import { InternationalConfig as N18 } from '../../../../service/base/international.service';

import { TaskSaveComponent }   from '../save/save.component';

@Component({
  selector: 'time-task-kanban',
  templateUrl: './kanban.component.html',
  styleUrls: ['./kanban.component.css']
})
export class TaskKanbanComponent implements OnInit {
  @ViewChild(TaskSaveComponent)
  taskSaveComponent: TaskSaveComponent;

  kanbanId: string = "TaskKanban";
  kanbanData: Kanban[] = [];
  kanbanColumns: any;
  kanbanFields: any;
  kanbanWorkflow: any;
  kanbanSwimlaneSettings: any;
  kanbanCardSettings: any;
  kanbanFilterSettings: any;
  kanbanConstraint: number;
  kanbanMenuItem: any;
  kanbanMenuItems: any = [];
  kanbanCustomMenuItems: any;
  updateBeforeTag: string;

  constructor(
    private taskService: TaskService,
    private fieldService: FieldService,
    private syncfusionService: SyncfusionService,
    private taskSettings: TaskSettings,
  ) { }

  ngOnInit() {
    this.syncfusionService.ListKanban().subscribe(res => {
      this.kanbanData = res;
    });
    this.kanbanConstraint = 2;
    this.kanbanFields = {
      primaryKey: "Id",
      content: "Desc",
      tag: "Tags",
      title: "Name",
      color: "FieldId",
      swimlaneKey: "ProjectName",
    }
    this.fieldService.List(null).subscribe(res => {
      let colorMaps = new Object;
      // this.kanbanFilterSettings = [];
      res.forEach((one, k) => {
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
    this.kanbanColumns = [
      { 
        headerText: this.taskSettings.StatusInfo[this.taskSettings.Status.Backlog], 
        key: this.taskSettings.StatusName[this.taskSettings.Status.Backlog],
        totalCount: { text: taskName },
        isCollapsed: true,
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
      // {
      //   text: processCreate + taskName,
      // },
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
    let processCreate = N18.settings.SYSTEM.PROCESS.CREATE;
    let processUpdate = N18.settings.SYSTEM.PROCESS.UPDATE;
    let processDelete = N18.settings.SYSTEM.PROCESS.DELETE;
    switch ($event.text) {
      case processUpdate + taskName:
        this.taskSaveComponent.New($event.cardData.ProjectId, $event.cardData.Id);
        this.updateBeforeTag = $event.cardData.Tags;
        break;
      case processDelete + taskName:
        this.taskService.Delete($event.cardData.Id).subscribe(res => {
          let kanbanObj = $("#" + this.kanbanId).data("ejKanban");
          kanbanObj.KanbanEdit.deleteCard($event.cardData.Id);
        });
        break;
    }
  }

  disable($event): void {
    $event.cancel = true;
  }

  taskSaved($event): void {
    console.log($event);
    console.log(this.updateBeforeTag);
    let tagsArray = this.updateBeforeTag.split(",");
    tagsArray[2] = $event.Resource.Name;
    let kanbanObj = $("#" + this.kanbanId).data("ejKanban"); 
    kanbanObj.updateCard(this.taskSettings.StatusName[$event.Status], [{
      Id: $event.Id, 
      Name: $event.Name, 
      Desc: $event.Desc, 
      Tags: tagsArray.join()
    }])
  }

  onTaskDrop($event): void {
    let statusNameMap = this.taskSettings.StatusName
    if ($event.data.length == 0) {
      return
    }
    for (let key in statusNameMap) {
      if ((statusNameMap[key] == $event.data[0].StatusName) && (key != $event.data[0].Status)) {
        let task = new Task();
        task.Id = $event.data[0].Id;
        task.Status = +key;

        // source status
        if ($event.data[0].Status == this.taskSettings.Status.Backlog) {
          task.StartDate = new Date();
        }
        if ($event.data[0].Status == this.taskSettings.Status.Done) {
          task.EndDateReset = true;
          task.Duration = 0;
        }

        // target status
        if (+key == this.taskSettings.Status.Backlog) {
          task.StartDateReset = true;
        }
        if (+key == this.taskSettings.Status.Done) {
          task.EndDate = new Date();
        }
        this.taskService.Update(task).subscribe(res => {
          let kanbanObj = $("#" + this.kanbanId).data("ejKanban"); 
          kanbanObj.updateCard($event.data[0].StatusName, [{Id: $event.data[0].Id, Status: +key}])
        })
      }
    }
  }
}
