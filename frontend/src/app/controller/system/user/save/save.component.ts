import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { User }           from '../../../../model/system/user';
import { UserService  }   from '../../../../service/system/user.service';

@Component({
  selector: 'system-save-user',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveUserComponent implements OnInit {
  user: User = new User;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private userService: UserService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.userService.Get(id).subscribe(res => {
        this.user = res;
        this.modelOpened = true;
      })
    } else {
      this.user = new User();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.user.Id == null) {
      this.userService.Add(this.user).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.userService.Update(this.user).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
