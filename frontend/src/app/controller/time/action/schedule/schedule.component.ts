import { Component, OnInit, Input, Output, ViewChild } from '@angular/core';
import { EJ_SCHEDULE_COMPONENTS }                      from 'ej-angular2/src/ej/schedule.component';
import { InternationalConfig as N18 }   from '../../../../service/base/international.service';
import * as jsrender from 'jsrender'; 

import { Action }              from '../../../../model/time/action';
import { Field }              from '../../../../model/time/field';
import { Schedule }            from '../../../../model/time/syncfusion';
import { ActionService }       from '../../../../service/time/action.service';
import { SyncfusionService }   from '../../../../service/time/syncfusion.service';
import { FieldService }        from '../../../../service/time/field.service';
import { SignService }         from '../../../../service/system/sign.service';

import { ActionSaveComponent } from '../save/save.component';
import { ActionListComponent } from '../list/list.component';

@Component({
  selector: 'time-action-schedule',
  templateUrl: './schedule.component.html',
  styleUrls: ['./schedule.component.css']
})
export class ActionScheduleComponent implements OnInit {
  @ViewChild(ActionSaveComponent)
  actionSaveComponent: ActionSaveComponent;
  @ViewChild(ActionListComponent)
  actionListComponent: ActionListComponent;
  @Input() scheduleHeight: string = "600px";
  @Input() currentView: string = "week";
  @Input() scheduleViews: string[] = ["Day", "Agenda", "Week", "Month"];

  currentDate: Date = new Date();
  scheduleId: string = "ActionSchedule";
  scheduleSettings: any;
  scheduleCategorySettings: any;
  scheduleTooltipSettings: any;
  scheduleMenuItems: any;
  scheduleAppointmentTemplate: string = "#appTemplate";
  fieldsMap: Map<number, Field> = new Map();
  modelListOpened: boolean = false;

  constructor(
    private actionService: ActionService,
    private fieldService: FieldService,
    private syncfusionService: SyncfusionService,
    private signService: SignService,
  ) { 
    jsrender.views.helpers({ 
      minuteFormat: this.minuteFormat,
      conceptFormat: this.conceptFormat,
    });
  }

  minuteFormat(date: Date) {
    var dFormat = ej.format(new Date(date), "hh:mm");
    return dFormat;
  }

  conceptFormat(concept: string): string {
    let taskName     = N18.settings.TIME.RESOURCE.TASK.CONCEPT;
    let areaName     = N18.settings.TIME.RESOURCE.AREA.CONCEPT;
    let resourceName = N18.settings.TIME.RESOURCE.RESOURCE.CONCEPT;
    let actionName   = N18.settings.TIME.RESOURCE.ACTION.CONCEPT;
    let startDate    = N18.settings.TIME.RESOURCE.ACTION.STARTDATE;
    let endDate      = N18.settings.TIME.RESOURCE.ACTION.ENDDATE;
    switch (concept) {
      case "task":
        return taskName;
      case "area":
        return areaName;
      case "resource":
        return resourceName;
      case "action":
        return actionName;
      case "startDate":
        return startDate;
      case "endDate":
        return endDate;
      default:
        return "";
    }
  }

  ngOnInit() {
    this.scheduleSettings = {
      id: "Id",
      startTime: "StartDate",
      endTime: "EndDate",
      subject: "Name",
      allDay:"AllDay",
      recurrence:"Recurrence",
      recurrenceRule:"RecurrenceRule",
      categorize: "FieldId",
      dataSource: this.syncfusionService.GetScheduleManager()
    }

    this.scheduleTooltipSettings = {
      enable: true,
      templateId: "#time-schedule-tooltip-template",
    }

    this.fieldService.List(null).subscribe(fields => {
      let categories = [];
      fields.forEach((one, k) => {
        this.fieldsMap.set(one.Id, one);
        let tmp = {
          Name: one.Name,
          Id: one.Id,
          Color: one.Color,
          FontColor: "#ffffff",
        }
        categories.push(tmp);
      })
      this.scheduleCategorySettings = {
        enable: true,
        allowMultiple: false,
        dataSource: categories,
        text: "Name", id: "Id", color: "Color", fontColor: "FontColor",
      }
    })

    let actionName    = N18.settings.TIME.RESOURCE.ACTION.CONCEPT;
    let processCreate = N18.settings.SYSTEM.PROCESS.CREATE;
    let processUpdate = N18.settings.SYSTEM.PROCESS.UPDATE;
    let processDelete = N18.settings.SYSTEM.PROCESS.DELETE;
    let processList   = N18.settings.SYSTEM.PROCESS.LIST;

    this.scheduleMenuItems = {
      appointment: [
        {
          id: "list-action",
          text: actionName + processList
        },
        {
          id: "update-action",
          text: actionName + processUpdate
        },
        {
          id: "delete-action",
          text: actionName + processDelete
        },
      ],
      cells: [
        {
          id: "add-action",
          text: actionName + processCreate,
        }
      ]
    };
  }
  
  onMenuItemClick($event) {
    switch($event.events.ID) {
      case "list-action":
        let action = new Action();
        action.TaskId = $event.targetInfo.Task.Id;
        this.actionListComponent.NewWithFilter(action)
        this.modelListOpened = true;
        break;
      case "add-action":
        let startDate = $event.targetInfo.startTime;
        let endDate = $event.targetInfo.endTime;
        this.actionSaveComponent.NewWithDate(startDate, endDate);
        break;
      case "update-action":
        this.actionSaveComponent.New($event.targetInfo.Id);
        break;
      case "delete-action":
        this.actionService.Delete($event.targetInfo.Id).subscribe(res => {
          if (res) {
            let schObj = $("#" + this.scheduleId).data("ejSchedule");
            schObj.deleteAppointment($event.targetInfo.Guid); 
            $event.cancel = false;
          } else {
            $event.cancel = true;
          }
        });
        break;
    }
  }

  actionComplete($event): void {
    var schObj = $("#" + this.scheduleId).data("ejSchedule");
    let currentView = schObj["ob.values"].currentView;
    switch (currentView) {
      case "agenda":
        switch ($event.requestType) {
          case "appointmentDelete":
            if ($event.data != undefined) {
              this.actionService.Delete($event.data.Id).subscribe(res => {
                if (res) {
                  $event.cancel = false;
                } else {
                  $event.cancel = true;
                }
              });
            }
        }
        break;
    }
  }

  actionDelete($event) {
    $event.cancel = false;
  }

  onContextMenuOpen($event) {
  }

  disable($event): void {
    $event.cancel = true;
  }

  disableAndRefresh($event): void {
    $event.cancel = true;
    $("#" + this.scheduleId).data("ejSchedule").refresh();
  }

  updateTimeRange($event) {
    let action = new Action();
    action.Id = $event.appointment.Id
    action.StartDate = $event.appointment.StartDate; 
    action.EndDate = $event.appointment.EndDate; 
    let time = action.EndDate.getTime() - action.StartDate.getTime();
    action.Time = time / 1000 / 60;
    this.actionService.Update(action).subscribe(res => {
      $event.cancel = false;
    })
  }

  actionShow($event) {
    var schObj = $("#" + this.scheduleId).data("ejSchedule");
    let currentView = schObj["ob.values"].currentView;
    switch (currentView) {
      case "agenda":
        this.actionSaveComponent.New($event.appointment.Id);
        break;
    }
    $event.cancel = true;
  }

  actionSaved($event): void {
    this.actionService.Get($event.Id).subscribe(res => {
      let schedule = new Schedule();
      schedule.Id = res.Id;
      schedule.Name = res.Name;
      schedule.StartDate = res.StartDate;
      schedule.EndDate = res.EndDate;
      schedule.Field = this.fieldsMap.get(res.Area.FieldId);
      schedule.FieldId = res.Area.FieldId.toString();
      schedule.Area = res.Area;
      schedule.Task = res.Task;
      schedule.Resource = res.Resource;
      schedule.AllDay = false;
      schedule.Recurrence = false;
      var schObj = $("#" + this.scheduleId).data("ejSchedule");
      schObj.saveAppointment(schedule); 
    })
  }
}
