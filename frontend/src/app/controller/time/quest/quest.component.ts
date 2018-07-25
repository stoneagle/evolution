import { Component, OnInit, ViewChild, Inject, forwardRef } from '@angular/core';
import { Comparator, State, SortOrder}                      from "clarity-angular";

import { PageSet }                           from '../../../model/base/basic';
import { Quest, QuestTarget, QuestSettings } from '../../../model/time/quest';
import { Project }                           from '../../../model/time/project';
import { QuestService }                      from '../../../service/time/quest.service';
import { ProjectService }                    from '../../../service/time/project.service';
import { QuestTargetService }                from '../../../service/time/quest-target.service';
import { MessageHandlerService  }            from '../../../service/base/message-handler.service';

import { QuestSaveComponent }                from './save/save.component';
import { QuestTeamListComponent }            from './team-list/team-list.component';
import { ShellComponent }                    from '../../../base/shell/shell.component';

import { ShareSettings }                                from '../../../shared/settings';
import { ErrorInfo }                                    from '../../../shared/error';
import { CustomComparator }                             from '../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../shared/utils';
import { PageSize }                                     from '../../../shared/const';

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
  preSorted = SortOrder.Desc;
  currentState: State;
  pageSet: PageSet = new PageSet();

  statusOptions: Array<{key: number, value: number}>;
  constraintOptions: Array<{key: number, value: number}>;
  membersOptions: Array<{key: number, value: number}>;

  modelListOpened: boolean = false;

  constructor(
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
    private questSettings: QuestSettings,
    private questTargetService: QuestTargetService,
    private questService: QuestService,
    private projectService: ProjectService,
    private messageHandlerService: MessageHandlerService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { 
  }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
    this.pageSet.Current = 1;
    this.statusOptions = this.questSettings.StatusFilterOptions;
    this.membersOptions = this.questSettings.MembersFilterOptions;
    this.constraintOptions = this.questSettings.ConstraintFilterOptions;
  }

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  load(state: State): void {
    let quest = new Quest();
    quest = loadPageFilterSort<Quest>(quest, state);
    this.pageSet.Current = quest.Page.Current;
    this.currentState = state;
    this.questService.Count(quest).subscribe(count => {
      this.pageSet.Count = count;
      this.questService.List(quest).subscribe(res => {
        this.quests = res;
      })
    })
  }

  saved(savedQuest: Quest): void {
    if (savedQuest.Id) {
      this.load(this.currentState);
    } else {
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
    this.modelListOpened = true;
  }

  exec(quest: Quest): void {
    quest.Status = this.questSettings.Status.Exec;
    quest.StartDate = new Date();
    this.questService.Update(quest).subscribe(res => {
      this.load(this.currentState);
    })
  }

  finish(quest: Quest): void {
    let questTarget = new QuestTarget();
    questTarget.QuestId = quest.Id;
    this.questTargetService.List(questTarget).subscribe(res => {
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
          this.load(this.currentState);
        })
      } else {
        this.messageHandlerService.showWarning(
          this.shareSettings.Time.Resource.Quest,
          this.shareSettings.System.Process.Update,
          this.errorInfo.Time.NotFinish
        );
      }
    })
  }

  fail(quest: Quest): void {
    let project = new Project();
    project.QuestTarget.QuestId = quest.Id;
    this.projectService.List(project).subscribe(projects => {
      if (projects.length > 0) {
        this.messageHandlerService.showWarning(
          this.shareSettings.Time.Resource.Quest,
          this.shareSettings.System.Process.Update,
          this.shareSettings.Time.Resource.Project + this.errorInfo.Time.NotFinish
        );
      } else {
        quest.Status = this.questSettings.Status.Fail;
        this.questService.Update(quest).subscribe(res => {
          this.load(this.currentState);
        })
      }
    })
  }

  delete(quest: Quest): void {
    let project = new Project();
    project.QuestTarget.QuestId = quest.Id;
    this.projectService.List(project).subscribe(projects => {
      if (projects.length > 0) {
        this.messageHandlerService.showWarning(
          this.shareSettings.Time.Resource.Quest,
          this.shareSettings.System.Process.Update,
         this.shareSettings.Time.Resource.Project + this.errorInfo.Time.Execing
        );
      } else {
        this.questService.Delete(quest.Id).subscribe(res => {
          let state = deleteState(this.pageSet, this.currentState, 1);
          this.load(state);
        })
      }
    })
  }
}
