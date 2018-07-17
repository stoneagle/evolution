import { Component, OnInit, ViewChild, Inject, forwardRef } from '@angular/core';

import { Quest, QuestTarget, QuestSettings } from '../../../model/time/quest';
import { QuestService }                      from '../../../service/time/quest.service';
import { QuestSaveComponent }                from './save/save.component';
import { QuestTeamListComponent }            from './team-list/team-list.component';
import { ShellComponent }                    from '../../../base/shell/shell.component';
import { QuestTargetService }                from '../../../service/time/quest-target.service';
import { MessageHandlerService  }            from '../../../service/base/message-handler.service';
import { ShareSettings }                     from '../../../shared/settings';
import { ErrorInfo }                         from '../../../shared/error';

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

  constructor(
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
    private questSettings: QuestSettings,
    private questTargetService: QuestTargetService,
    private questService: QuestService,
    private messageHandlerService: MessageHandlerService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { 
  }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(quest?: Quest): void {
    if (quest != undefined) {
      this.saveComponent.New(quest.FounderId, quest.Id);
    } else {
      this.saveComponent.New(this.shell.currentUser.Id);
    }
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
    quest.Status = this.questSettings.Status.Exec;
    this.questService.Update(quest).subscribe(res => {
      this.refresh();
    })
  }

  finish(quest: Quest): void {
    let questTarget = new QuestTarget();
    questTarget.QuestId = quest.Id;
    this.questTargetService.ListWithCondition(questTarget).subscribe(res => {
      let updateFlag = true;
      for (let k in res) {
        if (res[k].Status == this.questSettings.TargetStatus.Wait) {
          updateFlag = false;
          break;
        }
      }
      if (updateFlag) {
        quest.Status = this.questSettings.Status.Finish;
        this.questService.Update(quest).subscribe(res => {
          this.refresh();
        })
      } else {
        this.messageHandlerService.showWarning(
          this.shareSettings.Time.Resource.Quest,
          this.shareSettings.System.Process.Update,
          this.errorInfo.Time.TargetNotFinish
        );
      }
    })
  }

  fail(quest: Quest): void {
    quest.Status = this.questSettings.Status.Fail;
    this.questService.Update(quest).subscribe(res => {
      this.refresh();
    })
  }
}
