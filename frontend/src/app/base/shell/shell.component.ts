import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { SignService  } from '../../service/base/sign.service';

@Component({
  selector: 'app-shell',
  templateUrl: './shell.component.html',
  styleUrls: ['./shell.component.css']
})
export class ShellComponent implements OnInit {

  constructor(
    private signService: SignService,
    private router: Router
  ) { }

  ngOnInit() {
  }

  logout(): void {
    this.signService.logout()
    .subscribe(res => {
      this.router.navigate(['/login']);
    })
  }
}
