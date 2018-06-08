import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router";
import { AppConfig } from '../../service/base/config.service';
import { BaseService  } from '../../service/base/base.service';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit {
  CompanyName: string;
  username: string;
  password: string;

  constructor(
    private baseService: BaseService,
    private router: Router
  ) { }

  ngOnInit() {
    // this.CompanyName = AppConfig.settings.app.company;
    this.CompanyName = 'quant';
  }

  submit() {
    this.baseService.login(this.username, this.password)
    .subscribe(res => {
      if (!this.baseService.checkLoginError()) {
        this.router.navigate(['/quant']);
      }
    })
  }
}
