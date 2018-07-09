import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Quest }           from '../../../../model/time/quest';
import { QuestService  }   from '../../../../service/time/quest.service';

@Component({
  selector: 'time-save-quest',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveQuestComponent implements OnInit {
  quest: Quest = new Quest;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private questService: QuestService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.questService.Get(id).subscribe(res => {
        this.quest = res;
        this.modelOpened = true;
      })
    } else {
      this.quest = new Quest();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.quest.Id == null) {
      this.questService.Add(this.quest).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.questService.Update(this.quest).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
