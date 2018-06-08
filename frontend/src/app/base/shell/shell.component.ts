import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { BaseService  } from '../../service/base/base.service';

@Component({
  selector: 'app-shell',
  templateUrl: './shell.component.html',
  styleUrls: ['./shell.component.css']
})
export class ShellComponent implements OnInit {

  constructor(
    private baseService: BaseService,
    private router: Router
  ) { }

  ngOnInit() {
  }

  logout(): void {
    this.baseService.logout()
    .subscribe(res => {
      this.router.navigate(['/login']);
    })
  }
}
