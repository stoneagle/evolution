import { Component, OnInit, ViewChild, Inject, forwardRef } from '@angular/core';
import { Comparator, State } from "clarity-angular";

import { ShellComponent }    from '../../../base/shell/shell.component';

@Component({
  selector: 'time-base',
  templateUrl: './base.component.html',
  styleUrls: ['./base.component.css']
})
export class BaseComponent {
  @Inject(forwardRef(() => ShellComponent))
  protected shell: ShellComponent;

  constructor(
  ) { 
  }
}
