import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';

import { GanttProjectComponent } from './gantt/gantt.component';
import { Project }               from '../../../model/time/project';
import { ProjectService }        from '../../../service/time/project.service';

@Component({
  selector: 'time-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {
  @ViewChild(GanttProjectComponent)
  ganttProject: GanttProjectComponent;

  constructor(
    private projectService: ProjectService,
  ) { 
  }

  ngOnInit() {
  }
}
