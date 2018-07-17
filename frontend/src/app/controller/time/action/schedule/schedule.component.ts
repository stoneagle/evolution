import { Component, OnInit, Input, Output, ViewChild } from '@angular/core';
import { EJ_SCHEDULE_COMPONENTS }                      from 'ej-angular2/src/ej/schedule.component';
import { TranslateService }                            from '@ngx-translate/core';

import { Action }              from '../../../../model/time/action';
import { Schedule }            from '../../../../model/time/syncfusion';
import { ActionService }       from '../../../../service/time/action.service';
import { SyncfusionService }   from '../../../../service/time/syncfusion.service';
import { SignService }         from '../../../../service/system/sign.service';
import { ActionSaveComponent } from '../save/save.component';

@Component({
  selector: 'time-action-schedule',
  templateUrl: './schedule.component.html',
  styleUrls: ['./schedule.component.css']
})
export class ActionScheduleComponent implements OnInit {
  @ViewChild(ActionSaveComponent)
  actionSaveComponent: ActionSaveComponent;

  currentDate: Date = new Date();
  scheduleId: string = "ActionSchedule";
  scheduleSettings: any;
  scheduleData: any;
  scheduleMenuItems: any;

  constructor(
    private actionService: ActionService,
    private syncfusionService: SyncfusionService,
    private signService: SignService,
    private translateService: TranslateService,
  ) { }

  ngOnInit() {
    this.scheduleSettings = {
      id: "Id",
      startTime: "StartDate",
      endTime: "EndDate",
      subject: "Name",
      allDay:"AllDay",
      recurrence:"Recurrence",
      recurrenceRule:"RecurrenceRule",
    }
    let dataManager = this.syncfusionService.GetScheduleManager();
    this.scheduleData = dataManager;
    this.translateService.get('TIME.RESOURCE.ACTION.CONCEPT').subscribe(actionName => {
      let processCreate = this.translateService.instant('SYSTEM.PROCESS.CREATE');
      let processUpdate = this.translateService.instant('SYSTEM.PROCESS.UPDATE');
      let processDelete = this.translateService.instant('SYSTEM.PROCESS.DELETE');
      this.scheduleMenuItems = {
        appointment: [
          {
              id: "update-action",
              text: processUpdate + actionName
          },
          {
              id: "delete-action",
              text: processDelete + actionName
          },
        ],
        cells: [
          {
            id: "add-action",
            text: processCreate + actionName,
          }
        ]
      };
    });
  }
  
  onMenuItemClick($event) {
    switch($event.events.ID) {
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
          let schObj = $("#" + this.scheduleId).data("ejSchedule");
          schObj.deleteAppointment($event.targetInfo.Guid); 
        });
        break;
    }
  }

  actionDelete($event) {
    this.actionService.Delete($event.appointment.Id).subscribe(res => {
      $event.cancel = false;
    });
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
    $event.cancel = true;
  }

  actionSaved($event): void {
    let schedule = new Schedule();
    schedule.Id = $event.Id;
    schedule.Name = $event.Name;
    schedule.StartDate = $event.StartDate;
    schedule.EndDate = $event.EndDate;
    schedule.AllDay = false;
    schedule.Recurrence = false;
		var schObj = $("#" + this.scheduleId).data("ejSchedule");
		schObj.saveAppointment(schedule); 
  }
}
