import { Component, OnInit, Input, ViewChild } from '@angular/core';

import { Project }               from '../../../model/time/project';
import { ProjectService }        from '../../../service/time/project.service';
import { ProjectGanttComponent } from './gantt/gantt.component';

@Component({
  selector: 'time-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {
  @ViewChild(ProjectGanttComponent)
  gantt: ProjectGanttComponent;

  constructor(
    private projectService: ProjectService,
  ) { }

  ngOnInit() {
  }
}
