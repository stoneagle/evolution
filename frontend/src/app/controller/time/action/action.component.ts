import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Action }                  from '../../../model/time/action';
import { ActionService }           from '../../../service/time/action.service';
import { ActionScheduleComponent } from './schedule/schedule.component';

@Component({
  selector: 'time-action',
  templateUrl: './action.component.html',
  styleUrls: ['./action.component.css']
})
export class ActionComponent implements OnInit {
  @ViewChild(ActionScheduleComponent)
  actionScheduleComponent: ActionScheduleComponent;

  actions: Action[];

  constructor(
    private actionService: ActionService,
  ) { }

  ngOnInit() {
  }
}
