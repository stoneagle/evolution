import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'shell-nav',
  templateUrl: './shell-nav.component.html',
  styleUrls: ['./shell-nav.component.css']
})
export class ShellNavComponent implements OnInit {
  @Input() nav: string;

  constructor(
  ) { 
  }

  ngOnInit() {
  }
}
