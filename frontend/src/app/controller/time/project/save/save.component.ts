import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Project }          from '../../../../model/time/project';
import { Quest }            from '../../../../model/time/quest';
import { Area }             from '../../../../model/time/area';
import { Resource }         from '../../../../model/time/resource';
import { ProjectService  }  from '../../../../service/time/project.service';
import { ResourceService  } from '../../../../service/time/resource.service';

@Component({
  selector: 'time-project-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class ProjectSaveComponent implements OnInit {
  project: Project = new Project;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private projectService: ProjectService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.projectService.Get(id).subscribe(res => {
        this.project = res;
        this.modelOpened = true;
      })
    } else {
      this.project = new Project();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.project.Id == null) {
      this.projectService.Add(this.project).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.projectService.Update(this.project).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
