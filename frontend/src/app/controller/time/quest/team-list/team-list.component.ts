import { Component, OnInit, ViewChild } from '@angular/core';
import { Comparator, State, SortOrder}  from "clarity-angular";

import { PageSet }          from '../../../../model/base/basic';
import { Quest, QuestTeam } from '../../../../model/time/quest';
import { User }             from '../../../../model/system/user';
import { QuestService }     from '../../../../service/time/quest.service';
import { QuestTeamService } from '../../../../service/time/quest-team.service';
import { UserService }      from '../../../../service/system/user.service';

import { CustomComparator }                             from '../../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../../shared/utils';
import { PageSize }                                     from '../../../../shared/const';

@Component({
  selector: 'time-quest-team-list',
  templateUrl: './team-list.component.html',
  styleUrls: ['./team-list.component.css']
})
export class QuestTeamListComponent implements OnInit {
  questTeams: QuestTeam[] = [];

  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  userNameMap: Map<number, string> = new Map(); 
  questId: number;

  constructor(
    private questService: QuestService,
    private questTeamService: QuestTeamService,
    private userService: UserService,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: State): void {
    if (this.questId != undefined) {
      let questTeam = new QuestTeam();
      questTeam.QuestId = this.questId;
      questTeam = loadPageFilterSort<QuestTeam>(questTeam, state);
      this.pageSet.Current = questTeam.Page.Current;
      this.currentState = state;
      this.questTeamService.Count(questTeam).subscribe(count => {
        this.pageSet.Count = count;
        this.questTeamService.List(questTeam).subscribe(res => {
          this.questTeams = res;
          let user = new User();
          user.Ids = [];
          res.forEach((one, k) => {
            user.Ids.push(one.UserId);
          });
          this.userService.List(user).subscribe(res => {
            this.userNameMap = new Map();
            res.forEach((u, k) => {
              this.userNameMap.set(u.Id, u.Name);
            }) 
          });
        })
      })
    }
  }

  New(questId: number): void {
    this.questId = questId;
    this.refresh();
  }            
}
