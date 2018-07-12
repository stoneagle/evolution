import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Project }         from '../../../../model/time/project';
import { Quest }           from '../../../../model/time/quest';
import { Area }            from '../../../../model/time/area';
import { Entity }          from '../../../../model/time/entity';
import { ProjectService  } from '../../../../service/time/project.service';
import { EntityService  }  from '../../../../service/time/entity.service';

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
