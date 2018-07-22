import { Component, OnInit, ViewChild, Output, Input, EventEmitter }  from '@angular/core';
import { ClrWizard }                                                  from "@clr/angular";
import * as _                                                         from 'lodash';
import { NgForm  }                                                    from '@angular/forms';
import { EJ_SCHEDULE_COMPONENTS }       from 'ej-angular2/src/ej/schedule.component';

import { User }                                         from '../../../../model/system/user';
import { Quest, QuestTarget, QuestTeam, QuestSettings } from '../../../../model/time/quest';
import { Area }                                         from '../../../../model/time/area';
import { Project }                                      from '../../../../model/time/project';

import { QuestService  }              from '../../../../service/time/quest.service';
import { ProjectService  }            from '../../../../service/time/project.service';
import { QuestTargetService  }        from '../../../../service/time/quest-target.service';
import { QuestTeamService  }          from '../../../../service/time/quest-team.service';
import { UserService  }               from '../../../../service/system/user.service';
import { MessageHandlerService  }     from '../../../../service/base/message-handler.service';
import { AreaTreeGridComponent }      from '../../area/tree-grid/tree-grid.component';
import { ShareSettings }              from '../../../../shared/settings';
import { ErrorInfo }                  from '../../../../shared/error';


@Component({
  selector: 'time-quest-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class QuestSaveComponent implements OnInit {
  constructor(
    private errorInfo: ErrorInfo,
    private shareSettings: ShareSettings,
    private questService: QuestService,
    private projectService: ProjectService,
    private questTargetService: QuestTargetService,
    private questTeamService: QuestTeamService,
    private userService: UserService,
    private questSettings: QuestSettings,
    private messageHandlerService: MessageHandlerService,
  ) { 
  }

  @ViewChild("wizard") 
  wizard: ClrWizard;
  @ViewChild("questForm") 
  questForm: NgForm;
  @ViewChild(AreaTreeGridComponent) 
  areaTreeGrid: AreaTreeGridComponent;

  founderUser: User           = new User();
  quest: Quest                = new Quest;
  questTargets: QuestTarget[] = [];
  teams: QuestTeam[]          = [];

  _: any = _;

  membersMap         = this.questSettings.Members;
  constraintMap      = this.questSettings.Constraint;

  modelOpened: boolean = false;
  @Output() save = new EventEmitter<boolean>();

  ngOnInit() {
  }

  New(userId: number, id?: number): void {
    this.wizard.reset();
    this.userService.Get(userId).subscribe(user => {
      this.founderUser = user;
    });
    if (id) {
      this.questService.Get(id).subscribe(res => {
        this.quest = res;
        let questTarget = new QuestTarget();
        questTarget.QuestId = this.quest.Id;
        this.questTargets = []; 
        this.questTargetService.List(questTarget).subscribe(res => {
          res.forEach((target, k) => {
            this.questTargets.push(target);
          });
        });
        this.modelOpened = true;
      })
    } else {
      this.quest        = new Quest();
      this.questTargets = [];
      this.modelOpened  = true;
    }
  }            

  addTargetArea($event: Area) {
    if ($event.Id != undefined) {
      let addFlag = true;
      this.questTargets.forEach((one, k) => {
        if (one.Area.Id == $event.Id) {
          addFlag = false;
          return;
        }
      });
      if (addFlag) {
        let questTarget = new QuestTarget();
        questTarget.QuestId = this.quest.Id;
        questTarget.AreaId = $event.Id;
        questTarget.Area = $event;
        this.questTargets.push(questTarget);
      }
    }
  }

  deleteQuestTarget(questTarget: QuestTarget) {
    this.questTargets.forEach((one, k) => {
      if (one.AreaId === questTarget.AreaId) {
        if (one.Id != undefined) {
          let project = new Project();
          project.QuestTarget.Id = one.Id;
          this.projectService.List(project).subscribe(projects => {
            if (projects.length > 0) {
              this.messageHandlerService.showWarning(
                this.shareSettings.Time.Resource.Quest,
                this.shareSettings.System.Process.Update,
                this.errorInfo.Time.TargetNotFinish
              );
            } else {
              this.questTargets.splice(k, 1); 
            } 
            return;
          })
        } else {
          this.questTargets.splice(k, 1); 
          return;
        }
      }
    });
  }

  onTargetCommit(): void {
    if (this.questTargets.length <= 0) {
      return;
    }
    if (this.questTargets.length >= 10) {
      return;
    }
    // TODO 新增quest其他配置后再支持
    // this.wizard.forceNext();
  }

  finish(): void {
    this.quest.EndDate = new Date(this.quest.EndDate);
    this.quest.FounderId = this.founderUser.Id;
    if (this.quest.Id == null) {
      this.quest.Status = this.questSettings.Status.Recruit; 
      this.questService.Add(this.quest).subscribe(res => {
        this.saveQuestTarget(res.Id);
        this.saveQuestTeamFounder(res.Id, res.EndDate);
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      if (this.quest.Status != this.questSettings.Status.Recruit) {
        this.quest.StartDate = new Date(this.quest.StartDate);
      }
      this.questService.Update(this.quest).subscribe(res => {
        this.saveQuestTarget(this.quest.Id);
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }

  saveQuestTarget(questId: number): void {
    let batchQuestTargets: QuestTarget[] = [];
    for (let k in this.questTargets) {
      let tmpQuestTarget: QuestTarget = this.questTargets[k];
      tmpQuestTarget.QuestId = questId;
      if ((tmpQuestTarget.Status == 0) || (tmpQuestTarget.Status == undefined)) {
        tmpQuestTarget.Status = this.questSettings.TargetStatus.Wait; 
      }
      batchQuestTargets.push(tmpQuestTarget);
    }
    this.questTargetService.BatchSave(batchQuestTargets).subscribe(res => {
    })
  }

  saveQuestTeamFounder(questId: number, endDate: Date): void {
    let questTeam: QuestTeam = new QuestTeam();
    questTeam.QuestId = questId;
    questTeam.StartDate = new Date();
    questTeam.EndDate = endDate;
    questTeam.UserId = this.founderUser.Id; 
    this.questTeamService.Add(questTeam).subscribe(res => {
    });
  }

  onCancel(): void {
    this.wizard.reset();
    this.wizard.close();
  }

  onPrevious(): void {
    this.wizard.previous();
  }
}
