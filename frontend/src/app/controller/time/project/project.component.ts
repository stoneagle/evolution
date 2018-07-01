import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute }       from "@angular/router";
import { Project }                         from '../../../model/time/project';
import { ProjectService }                  from '../../../service/time/project.service';

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {
  projects: Project[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private projectService: ProjectService,
    private router: Router, 
    private route: ActivatedRoute 
  ) { 
  }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  delete(project: Project): void {
    this.projectService.Delete(project.id)
    .subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshProjects(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshProjects(0, 10);
  }

  refreshProjects(from: number, to: number): void {
    this.projectService.List().subscribe(res => {
      this.totalCount = res.length;
      this.projects = res.slice(from, to);
    })
  }
}
