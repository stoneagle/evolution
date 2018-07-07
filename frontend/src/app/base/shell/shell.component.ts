import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute }       from "@angular/router";
import { Location  }                    from '@angular/common';
import { SignService  }                 from '../../service/system/sign.service';
import { ShellNavComponent }            from './nav/shell-nav.component';
import { SessionUser }                  from '../../model/base/sign';

@Component({
  selector: 'app-shell',
  templateUrl: './shell.component.html',
  styleUrls: ['./shell.component.css']
})
export class ShellComponent implements OnInit {
  nav: string;
  currentUser: SessionUser = new SessionUser();

  @ViewChild(ShellNavComponent)
  shellNav: ShellNavComponent;

  constructor(
    private signService: SignService,
    private router: Router,
    private location: Location,
    private route: ActivatedRoute 
  ) { 
  }

  ngOnInit() {
    if (this.nav === undefined) {
      this.nav = location.pathname.split("/")[1];
    }
    this.signService.current().subscribe(res => {
      this.currentUser = res;
      if (this.currentUser.Name == undefined) {
        this.router.navigate(['/sign/login']);
      }
    });
  }

  logout(): void {
    this.signService.logout().subscribe(res => {
      this.router.navigate(['/sign/login']);
    })
  }

	changeNav(nav: string):void {
    this.nav = nav;
    this.router.navigate([`/${nav}`]);
	}
}
