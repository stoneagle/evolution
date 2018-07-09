import { Component, OnInit, ViewChild } from '@angular/core';

import { Project }              from '../../../../model/time/project';
import { ProjectService }       from '../../../../service/time/project.service';
import { SaveProjectComponent } from '../save/save.component';

@Component({
  selector: 'time-project-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ListProjectComponent implements OnInit {
  @ViewChild(SaveProjectComponent)
  saveProject: SaveProjectComponent;

  projects: Project[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private projectService: ProjectService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveProject.New(id);
  }

  delete(project: Project): void {
    this.projectService.Delete(project.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    this.projectService.List().subscribe(res => {
      this.totalCount = res.length;
      this.projects = res.slice(from, to);
    })
  }
}
