import { Component, OnInit, ViewChild } from '@angular/core';

import { Quest, QuestTeam } from '../../../../model/time/quest';
import { User }             from '../../../../model/system/user';
import { QuestService }     from '../../../../service/time/quest.service';
import { QuestTeamService } from '../../../../service/time/quest-team.service';
import { UserService }      from '../../../../service/system/user.service';

@Component({
  selector: 'time-quest-team-list',
  templateUrl: './team-list.component.html',
  styleUrls: ['./team-list.component.css']
})
export class ListQuestTeamComponent implements OnInit {
  questTeams: QuestTeam[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  userNameMap: Map<number, string> = new Map(); 
  questId: number;
  modelOpened: boolean = false;

  constructor(
    private questService: QuestService,
    private questTeamService: QuestTeamService,
    private userService: UserService,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  New(questId: number): void {
    this.questId = questId;
    this.modelOpened = true;
    this.refresh();
  }            

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  delete(questTeam: QuestTeam): void {
    this.questTeamService.Delete(questTeam.Id).subscribe(res => {
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
    let questTeam = new QuestTeam();
    questTeam.QuestId = this.questId;
    this.questTeamService.ListWithCondition(questTeam).subscribe(res => {
      let user = new User();
      user.Ids = [];
      res.forEach((one, k) => {
        user.Ids.push(one.UserId);
      });
      this.userService.ListWithCondition(user).subscribe(res => {
        this.userNameMap = new Map();
        res.forEach((u, k) => {
          this.userNameMap.set(u.Id, u.Name);
        }) 
      });
      this.totalCount = res.length;
      this.questTeams = res.slice(from, to);
    })
  }
}
