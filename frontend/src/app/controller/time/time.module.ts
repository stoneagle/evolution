import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";
import { TreeModule }                from 'ng2-tree';
import { EJ_GANTT_COMPONENTS }       from 'ej-angular2/src/ej/gantt.component';
import { EJ_TREEGRID_COMPONENTS }    from 'ej-angular2/src/ej/treegrid.component';
import { EJ_SCHEDULE_COMPONENTS }    from 'ej-angular2/src/ej/schedule.component';
import { EJ_DATEPICKER_COMPONENTS }  from 'ej-angular2/src/ej/datepicker.component';

import { QuestComponent }          from './quest/quest.component';
import { QuestSaveComponent }      from './quest/save/save.component';
import { QuestTeamListComponent }  from './quest/team-list/team-list.component';
import { ProjectComponent }        from './project/project.component';
import { ProjectGanttComponent }   from './project/gantt/gantt.component';
import { ProjectSaveComponent }    from './project/save/save.component';
import { CountryComponent }        from './country/country.component';
import { CountrySaveComponent }    from './country/save/save.component';
import { PhaseComponent }          from './phase/phase.component';
import { PhaseSaveComponent }      from './phase/save/save.component';
import { EntityListComponent }     from './entity/list/list.component';
import { EntitySaveComponent }     from './entity/save/save.component';
import { EntityTreeGridComponent } from './entity/tree-grid/tree-grid.component';
import { FieldComponent }          from './field/field.component';
import { FieldSaveComponent }      from './field/save/save.component';
import { AreaComponent }           from './area/area.component';
import { AreaNg2TreeComponent }    from './area/ng2-tree/ng2-tree.component';
import { ResourceComponent }       from './resource/resource.component';
import { ResourceEntityComponent } from './resource/entity-list/entity-list.component';

import { CountryService }         from '../../service/time/country.service';
import { ProjectService }         from '../../service/time/project.service';
import { QuestService }           from '../../service/time/quest.service';
import { QuestTargetService }     from '../../service/time/quest-target.service';
import { QuestEntityService }     from '../../service/time/quest-entity.service';
import { QuestTimeTableService }  from '../../service/time/quest-time-table.service';
import { QuestTeamService }       from '../../service/time/quest-team.service';
import { PhaseService }           from '../../service/time/phase.service';
import { FieldService }           from '../../service/time/field.service';
import { EntityService }          from '../../service/time/entity.service';
import { AreaService }            from '../../service/time/area.service';
import { ResourceService }        from '../../service/time/resource.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';

@NgModule({
  declarations: [
    QuestComponent,
    QuestSaveComponent,
    QuestTeamListComponent,
    ProjectComponent,
    ProjectGanttComponent,
    ProjectSaveComponent,
    CountryComponent,
    CountrySaveComponent,
    EntityListComponent,
    EntitySaveComponent,
    EntityTreeGridComponent,
    PhaseComponent,
    PhaseSaveComponent,
    FieldComponent,
    FieldSaveComponent,
    AreaComponent,
    AreaNg2TreeComponent,
    ResourceComponent,
    ResourceEntityComponent,
    EJ_GANTT_COMPONENTS,
    EJ_TREEGRID_COMPONENTS,
    EJ_SCHEDULE_COMPONENTS,
    EJ_DATEPICKER_COMPONENTS
  ],
  imports: [   
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    HttpModule,
    FormsModule,
    ClarityModule,
    TranslateModule,
    TreeModule,
  ],
  providers: [ 
    AreaService,
    ResourceService,
    ProjectService,
    QuestService,
    QuestTeamService,
    QuestTargetService,
    QuestEntityService,
    QuestTimeTableService,
    CountryService,
    FieldService,
    PhaseService,
    EntityService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
