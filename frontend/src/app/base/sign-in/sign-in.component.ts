import { Component, OnInit } from '@angular/core';
import { Router }            from "@angular/router";
import { AppConfig }         from '../../service/base/config.service';
import { SignService  }      from '../../service/system/sign.service';

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
    this.CompanyName = AppConfig.settings.app.name;
    this.signService.current().subscribe(res => {
      if (res.Name != undefined) {
        this.router.navigate(['/time']);
      }
    });
  }

  submit() {
    this.signService.login(this.username, this.password).subscribe(res => {
      if (res) {
        this.router.navigate(['/time']);
      }
    })
  }
}
