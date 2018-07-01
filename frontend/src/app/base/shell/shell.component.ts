import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute }       from "@angular/router";
import { Location  }                    from '@angular/common';
import { SignService  }                 from '../../service/base/sign.service';
import { ShellNavComponent }            from './nav/shell-nav.component';

@Component({
  selector: 'app-shell',
  templateUrl: './shell.component.html',
  styleUrls: ['./shell.component.css']
})
export class ShellComponent implements OnInit {
  nav: string;

  @ViewChild(ShellNavComponent)
  shellNav: ShellNavComponent;

  constructor(
    private signService: SignService,
    private router: Router,
    private location: Location,
    private route: ActivatedRoute 
  ) { 
    // this.nav = this.route.snapshot.routeConfig.path;
  }

  ngOnInit() {
    if (this.nav === undefined) {
      this.nav = location.pathname.split("/")[1];
    }
  }

  logout(): void {
    this.signService.logout().subscribe(res => {
      this.router.navigate(['/flow/stock']);
    })
  }

	changeNav(nav: string):void {
    this.nav = nav;
    this.router.navigate([`/${nav}`]);
	}
}
