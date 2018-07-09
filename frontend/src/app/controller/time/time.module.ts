import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";
import { TreeModule }                from 'ng2-tree';
import { EJ_GANTT_COMPONENTS }                     from 'ej-angular2/src/ej/gantt.component';

import { ProjectComponent }        from './project/project.component';
import { ListProjectComponent }    from './project/list/list.component';
import { SaveProjectComponent }    from './project/save/save.component';
import { GanttProjectComponent }   from './project/gantt/gantt.component';
import { CountryComponent }        from './country/country.component';
import { SaveCountryComponent }    from './country/save/save.component';
import { PhaseComponent }          from './phase/phase.component';
import { SavePhaseComponent }      from './phase/save/save.component';
import { EntityComponent }         from './entity/entity.component';
import { SaveEntityComponent }     from './entity/save/save.component';
import { FieldComponent }          from './field/field.component';
import { SaveFieldComponent }      from './field/save/save.component';
import { AreaComponent }           from './area/area.component';
import { TreasureComponent }       from './treasure/treasure.component';
import { TreasureEntityComponent } from './treasure/entity-list/entity-list.component';

import { CountryService }         from '../../service/time/country.service';
import { ProjectService }          from '../../service/time/project.service';
import { PhaseService }           from '../../service/time/phase.service';
import { FieldService }           from '../../service/time/field.service';
import { EntityService }          from '../../service/time/entity.service';
import { AreaService }            from '../../service/time/area.service';
import { TreasureService }        from '../../service/time/treasure.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';

@NgModule({
  declarations: [
    ProjectComponent,
    ListProjectComponent,
    SaveProjectComponent,
    GanttProjectComponent,
    CountryComponent,
    SaveCountryComponent,
    EntityComponent,
    SaveEntityComponent,
    PhaseComponent,
    SavePhaseComponent,
    FieldComponent,
    SaveFieldComponent,
    AreaComponent,
    TreasureComponent,
    TreasureEntityComponent,
    EJ_GANTT_COMPONENTS,
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
    TreasureService,
    ProjectService,
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
