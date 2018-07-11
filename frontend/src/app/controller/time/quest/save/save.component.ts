import { Component, OnInit, ViewChild, Output, Input, EventEmitter }  from '@angular/core';
import { ClrWizard }                                                  from "@clr/angular";
import * as _                                                         from 'lodash';
import { NgForm  }                                                    from '@angular/forms';
import { Quest, QuestTarget, QuestTeam, QuestEntity, QuestTimeTable } from '../../../../model/time/quest';

import { SessionUser }         from '../../../../model/base/sign';
import { Entity }         from '../../../../model/time/entity';
import { QuestService  }       from '../../../../service/time/quest.service';
import { SignService  }        from '../../../../service/system/sign.service';
import { Quest as QuestConst } from '../../../../shared/shared.const';

import { TreeGridEntityComponent } from '../../entity/tree-grid/tree-grid.component';

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
  @ViewChild(TreeGridEntityComponent) 
  treeGridEntity: TreeGridEntityComponent;

  currentUser: SessionUser = new SessionUser();
  quest: Quest                 = new Quest;
  targets: QuestTarget[]       = [];
  teams: QuestTeam[]           = [];
  entities: QuestEntity[]      = [];
  timeTables: QuestTimeTable[] = [];

  targetEntities: Entity[] = [];

  _: any = _;
  membersMap = QuestConst.MembersMap;
  constraintMap = QuestConst.ConstraintMap;

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

  addTargetEntity($event: Entity) {
    if ($event.Id != undefined) {
      let addFlag = true;
      this.targetEntities.forEach((one, k) => {
        if (one.Id == $event.Id) {
          addFlag = false;
          return;
        }
      });
      if (addFlag) {
        this.targetEntities.push($event);
      }
    }
  }

  deleteTargetEntity(entity: Entity) {
    this.targetEntities.forEach( (one, k) => {
      if (one.Id === entity.Id) {
        this.targetEntities.splice(k, 1); 
        return;
      }
    });
  }

  onTargetCommit(): void {
    if (this.targetEntities.length <= 0) {
      return;
    }
    if (this.targetEntities.length >= 5) {
      return;
    }
    this.wizard.forceNext();
  }

  finish(): void {
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
