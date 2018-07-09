import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';

import { GanttQuestComponent } from './gantt/gantt.component';
import { Quest }               from '../../../model/time/quest';
import { QuestService }        from '../../../service/time/quest.service';

@Component({
  selector: 'time-quest',
  templateUrl: './quest.component.html',
  styleUrls: ['./quest.component.css']
})
export class QuestComponent implements OnInit {
  @ViewChild(GanttQuestComponent)
  ganttQuest: GanttQuestComponent;

  constructor(
    private questService: QuestService,
  ) { 
  }

  ngOnInit() {
  }
}
