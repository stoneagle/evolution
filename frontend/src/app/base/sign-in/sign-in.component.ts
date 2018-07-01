import { Component, OnInit } from '@angular/core';
import { Router }            from "@angular/router";
import { AppConfig }         from '../../service/base/config.service';
import { SignService  }      from '../../service/base/sign.service';

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
    private signService: SignService,
    private router: Router
  ) { }

  ngOnInit() {
    this.CompanyName = 'quant';
  }

  submit() {
    this.signService.login(this.username, this.password)
    .subscribe(res => {
    })
  }
}
