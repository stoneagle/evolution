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

import { QuestComponent }            from './quest/quest.component';
import { QuestSaveComponent }        from './quest/save/save.component';
import { QuestTeamListComponent }    from './quest/team-list/team-list.component';
import { ProjectComponent }          from './project/project.component';
import { ProjectGanttComponent }     from './project/gantt/gantt.component';
import { ProjectSaveComponent }      from './project/save/save.component';
import { TaskComponent }             from './task/task.component';
import { TaskSaveComponent }         from './task/save/save.component';
import { CountryComponent }          from './country/country.component';
import { CountrySaveComponent }      from './country/save/save.component';
import { PhaseComponent }            from './phase/phase.component';
import { PhaseSaveComponent }        from './phase/save/save.component';
import { ResourceListComponent }     from './resource/list/list.component';
import { ResourceSaveComponent }     from './resource/save/save.component';
import { ResourceTreeGridComponent } from './resource/tree-grid/tree-grid.component';
import { FieldComponent }            from './field/field.component';
import { FieldSaveComponent }        from './field/save/save.component';
import { AreaComponent }             from './area/area.component';
import { AreaTreeGridComponent }     from './area/tree-grid/tree-grid.component';
import { AreaNg2TreeComponent }      from './area/ng2-tree/ng2-tree.component';
import { UserResourceComponent }     from './user-resource/user-resource.component';
import { UserResourceListComponent } from './user-resource/list/list.component';

import { CountryService }         from '../../service/time/country.service';
import { ProjectService }         from '../../service/time/project.service';
import { TaskService }            from '../../service/time/task.service';
import { QuestService }           from '../../service/time/quest.service';
import { QuestTargetService }     from '../../service/time/quest-target.service';
import { QuestResourceService }   from '../../service/time/quest-resource.service';
import { QuestTimeTableService }  from '../../service/time/quest-time-table.service';
import { QuestTeamService }       from '../../service/time/quest-team.service';
import { PhaseService }           from '../../service/time/phase.service';
import { FieldService }           from '../../service/time/field.service';
import { ResourceService }        from '../../service/time/resource.service';
import { AreaService }            from '../../service/time/area.service';
import { UserResourceService }    from '../../service/time/user-resource.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';

@NgModule({
  declarations: [
    QuestComponent,
    QuestSaveComponent,
    QuestTeamListComponent,
    ProjectComponent,
    ProjectGanttComponent,
    ProjectSaveComponent,
    TaskComponent,
    TaskSaveComponent,
    CountryComponent,
    CountrySaveComponent,
    ResourceListComponent,
    ResourceSaveComponent,
    ResourceTreeGridComponent,
    PhaseComponent,
    PhaseSaveComponent,
    FieldComponent,
    FieldSaveComponent,
    AreaComponent,
    AreaNg2TreeComponent,
    AreaTreeGridComponent,
    UserResourceComponent,
    UserResourceListComponent,
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
    ProjectService,
    TaskService,
    QuestService,
    QuestTeamService,
    QuestTargetService,
    QuestResourceService,
    QuestTimeTableService,
    CountryService,
    FieldService,
    PhaseService,
    ResourceService,
    UserResourceService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
