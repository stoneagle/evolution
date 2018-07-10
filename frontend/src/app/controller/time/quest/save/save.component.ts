import { Component, OnInit, ViewChild, Output, Input, EventEmitter }  from '@angular/core';
import { ClrWizard }                                                  from "@clr/angular";
import * as _                                                         from 'lodash';
import { NgForm  }                                                    from '@angular/forms';
import { Quest, QuestTarget, QuestTeam, QuestEntity, QuestTimeTable } from '../../../../model/time/quest';

import { SessionUser }         from '../../../../model/base/sign';
import { QuestService  }       from '../../../../service/time/quest.service';
import { SignService  }        from '../../../../service/system/sign.service';
import { Quest as QuestConst } from '../../../../shared/shared.const';

@Component({
  selector: 'time-save-quest',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SaveQuestComponent implements OnInit {
  @ViewChild("wizard") 
  wizard: ClrWizard;
  @ViewChild("questForm") 
  questForm: NgForm;

  currentUser: SessionUser = new SessionUser();
  quest: Quest                 = new Quest;
  targets: QuestTarget[]       = [];
  teams: QuestTeam[]           = [];
  entities: QuestEntity[]      = [];
  timeTables: QuestTimeTable[] = [];

  _: any = _;
  membersMap = QuestConst.Members;
  constraintMap = QuestConst.Constraint;

  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private questService: QuestService,
    private signService: SignService,
  ) { 
    this.signService.current().subscribe( res=> {
      this.currentUser = res;
    });
  }

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

  finish(): void {
    console.log("finish");
  }

  onCancel(): void {
    this.wizard.reset();
    this.wizard.close();
  }

  onQuestCommit(): void {
    // this.wizard.forceNext();
    
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
