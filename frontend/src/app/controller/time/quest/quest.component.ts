import { Component, OnInit, ViewChild } from '@angular/core';
import { Quest }                        from '../../../model/time/quest';
import { QuestService }                 from '../../../service/time/quest.service';
import { Quest as QuestConst }          from '../../../shared/const';
import { QuestSaveComponent }           from './save/save.component';
import { QuestTeamListComponent }       from './team-list/team-list.component';

@Component({
  selector: 'time-quest',
  templateUrl: './quest.component.html',
  styleUrls: ['./quest.component.css']
})
export class QuestComponent implements OnInit {
  @ViewChild(QuestSaveComponent)
  saveComponent: QuestSaveComponent;
  @ViewChild(QuestTeamListComponent)
  questTeamListComponent: QuestTeamListComponent;

  quests: Quest[];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  membersInfoMap    = QuestConst.MembersInfo;
  constraintInfoMap = QuestConst.ConstraintInfo;
  statusInfoMap     = QuestConst.StatusInfo;
  recruitStatus     = QuestConst.Status.Recruit
  execStatus        = QuestConst.Status.Exec

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
    this.saveComponent.New(id);
  }

  openTeamList(id: number): void {
    this.questTeamListComponent.New(id);
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

  exec(quest: Quest): void {
    quest.StartDate = new Date();
    quest.Status = QuestConst.Status.Exec;
    this.questService.Update(quest).subscribe(res => {
      this.refresh();
    })
  }

  finish(quest: Quest): void {
    quest.Status = QuestConst.Status.Finish;
    this.questService.Update(quest).subscribe(res => {
      this.refresh();
    })
  }

  fail(quest: Quest): void {
    quest.Status = QuestConst.Status.Fail;
    this.questService.Update(quest).subscribe(res => {
      this.refresh();
    })
  }
}
