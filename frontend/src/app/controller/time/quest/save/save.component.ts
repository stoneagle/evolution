import { Component, OnInit, ViewChild, Output, Input, EventEmitter }  from '@angular/core';
import { ClrWizard }                                                  from "@clr/angular";
import * as _                                                         from 'lodash';
import { NgForm  }                                                    from '@angular/forms';
import { EJ_SCHEDULE_COMPONENTS }       from 'ej-angular2/src/ej/schedule.component';

import { SessionUser }                    from '../../../../model/base/sign';
import { Quest, QuestTarget, QuestTeam, } from '../../../../model/time/quest';
import { Area }                           from '../../../../model/time/area';

import { QuestService  }         from '../../../../service/time/quest.service';
import { QuestTargetService  }   from '../../../../service/time/quest-target.service';
import { QuestTeamService  }     from '../../../../service/time/quest-team.service';
import { SignService  }          from '../../../../service/system/sign.service';
import { Quest as QuestConst }   from '../../../../shared/const';
import { AreaTreeGridComponent } from '../../area/tree-grid/tree-grid.component';


@Component({
  selector: 'time-quest-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class QuestSaveComponent implements OnInit {
  constructor(
    private questService: QuestService,
    private questTargetService: QuestTargetService,
    private questTeamService: QuestTeamService,
    private signService: SignService,
  ) { 
  }

  @ViewChild("wizard") 
  wizard: ClrWizard;
  @ViewChild("questForm") 
  questForm: NgForm;
  @ViewChild(AreaTreeGridComponent) 
  areaTreeGrid: AreaTreeGridComponent;

  currentUser: SessionUser = new SessionUser();
  quest: Quest             = new Quest;
  targets: QuestTarget[]   = [];
  targetAreas: Area[]      = [];
  teams: QuestTeam[]       = [];

  _: any = _;

  membersMap         = QuestConst.Members;
  membersInfoMap     = QuestConst.MembersInfo;
  constraintMap      = QuestConst.Constraint;
  constraintInfoMap  = QuestConst.ConstraintInfo;
  questRecruitStatus = QuestConst.Status.Recruit;

  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  ngOnInit() {
  }

  New(id?: number): void {
    this.wizard.reset();
    this.signService.current().subscribe( res=> {
      this.currentUser = res;
    });
    if (id) {
      this.questService.Get(id).subscribe(res => {
        this.quest = res;

        let questTarget = new QuestTarget();
        questTarget.QuestId = this.quest.Id;
        this.targetAreas = []; 
        this.questTargetService.ListWithCondition(questTarget).subscribe(res => {
          res.forEach((target, k) => {
            this.targetAreas.push(target.Area);
          });
        });

        this.modelOpened = true;
      })
    } else {
      this.quest = new Quest();
      this.targetAreas = [];
      this.modelOpened = true;
    }
  }            

  addTargetArea($event: Area) {
    if ($event.Id != undefined) {
      let addFlag = true;
      this.targetAreas.forEach((one, k) => {
        if (one.Id == $event.Id) {
          addFlag = false;
          return;
        }
      });
      if (addFlag) {
        this.targetAreas.push($event);
      }
    }
  }

  deleteTargetArea(area: Area) {
    this.targetAreas.forEach( (one, k) => {
      if (one.Id === area.Id) {
        this.targetAreas.splice(k, 1); 
        return;
      }
    });
  }

  onTargetCommit(): void {
    if (this.targetAreas.length <= 0) {
      return;
    }
    if (this.targetAreas.length >= 5) {
      return;
    }
    this.wizard.forceNext();
  }

  finish(): void {
    this.quest.EndDate = new Date(this.quest.EndDate);
    this.quest.FounderId = this.currentUser.Id;
    if (this.quest.Id == null) {
      this.quest.Status = QuestConst.Status.Recruit; 
      this.questService.Add(this.quest).subscribe(res => {
        if (res.Id != undefined) {
          this.saveQuestTarget(res.Id);
          this.saveQuestTeamFounder(res.Id, res.EndDate);
        }
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      if (this.quest.Status != QuestConst.Status.Recruit) {
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
    let questTargets: QuestTarget[] = [];
    this.targetAreas.forEach((one, k) => {
      let questTarget = new QuestTarget();
      questTarget.QuestId = questId;
      questTarget.AreaId = one.Id;
      questTarget.Status = QuestConst.TargetStatus.Wait;
      questTargets.push(questTarget);
    })
    this.questTargetService.BatchAdd(questTargets).subscribe(res => {
    })
  }

  saveQuestTeamFounder(questId: number, endDate: Date): void {
    let questTeam: QuestTeam = new QuestTeam();
    questTeam.QuestId = questId;
    questTeam.StartDate = new Date();
    questTeam.EndDate = endDate; 
    questTeam.UserId = this.currentUser.Id; 
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
