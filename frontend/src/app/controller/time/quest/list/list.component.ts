import { Component, OnInit, ViewChild } from '@angular/core';

import { Quest }              from '../../../../model/time/quest';
import { QuestService }       from '../../../../service/time/quest.service';
import { SaveQuestComponent } from '../save/save.component';

@Component({
  selector: 'time-quest-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class ListQuestComponent implements OnInit {
  @ViewChild(SaveQuestComponent)
  saveQuest: SaveQuestComponent;

  quests: Quest[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private questService: QuestService,
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
    this.saveQuest.New(id);
  }

  delete(quest: Quest): void {
    this.questService.Delete(quest.Id).subscribe(res => {
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
    this.questService.List().subscribe(res => {
      this.totalCount = res.length;
      this.quests = res.slice(from, to);
    })
  }
}
