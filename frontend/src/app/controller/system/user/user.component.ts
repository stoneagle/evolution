import { Component, OnInit, ViewChild } from '@angular/core';

import { User }             from '../../../model/system/user';
import { UserService }      from '../../../service/system/user.service';
import { SaveUserComponent }            from './save/save.component';

@Component({
  selector: 'system-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {
  @ViewChild(SaveUserComponent)
  saveUser: SaveUserComponent;

  users: User[];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private userService: UserService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.saveUser.New(id);
  }

  delete(user: User): void {
    this.userService.Delete(user.Id).subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    this.userService.List().subscribe(res => {
      this.totalCount = res.length;
      this.users = res.slice(from, to);
    })
  }
}
