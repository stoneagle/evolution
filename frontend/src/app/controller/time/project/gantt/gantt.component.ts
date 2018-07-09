import { Component, OnInit, ViewChild } from '@angular/core';
import { EJ_GANTT_COMPONENTS }                      from 'ej-angular2/src/ej/gantt.component';

import { Project }              from '../../../../model/time/project';
import { ProjectService }       from '../../../../service/time/project.service';

@Component({
  selector: 'time-gantt-project',
  templateUrl: './gantt.component.html',
  styleUrls: ['./gantt.component.css']
})
export class GanttProjectComponent implements OnInit {
  constructor(
    private projectService: ProjectService,
  ) { }

  ganttData =  [{
       taskID: 1,
       taskName: "Planning",
       startDate: "02/03/2014",
       endDate: "02/07/2014",
       progress: 100,
       duration: 5,
       subtasks: [
           { taskID: 2, taskName: "Plan timeline", startDate: "02/03/2014", endDate: "02/07/2014", duration: 5, progress: 100 },
           { taskID: 3, taskName: "Plan budget", startDate: "02/03/2014", endDate: "02/07/2014", duration: 5 },
           { taskID: 4, taskName: "Allocate resources", startDate: "02/03/2014", endDate: "02/07/2014", duration: 5, progress: 100 },
           { taskID: 5, taskName: "Planning complete", startDate: "02/07/2014", endDate: "02/07/2014", duration: 0, progress: 0 }
       ]
   }];

  toolbarsettings = {
    showToolbar: true,
      toolbarItems: [
        ej.Gantt.ToolbarItems.Add,
        ej.Gantt.ToolbarItems.Edit,
        ej.Gantt.ToolbarItems.Delete,
        ej.Gantt.ToolbarItems.Update,
        ej.Gantt.ToolbarItems.Cancel,
        ej.Gantt.ToolbarItems.Indent,
        ej.Gantt.ToolbarItems.Outdent,
        ej.Gantt.ToolbarItems.ExpandAll,
        ej.Gantt.ToolbarItems.CollapseAll
      ]
  };

  treeColumnIndex = 1;

  editsettings = {
      allowEditing: true,
      allowAdding: true,
      allowDeleting: true,
      allowIndent: true,
      editMode: 'cellEditing'
  };

  ngOnInit() {
  }

  ngAfterViewInit(): void {
  }
}
