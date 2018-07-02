import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";
import { TreeModule }                from 'ng2-tree';

import { ProjectComponent }       from './project/project.component';
import { ProjectService }         from '../../service/time/project.service';
import { CountryComponent }       from './country/country.component';
import { SaveCountryComponent }   from './country/save/save.component';
import { AreaComponent }          from './area/area.component';
import { EntityLifeComponent }       from './entity/life/life.component';
import { SaveEntityLifeComponent }   from './entity/life/save/save.component';

import { CountryService }         from '../../service/time/country.service';
import { EntityLifeService }         from '../../service/time/entity-life.service';
import { AreaService }         from '../../service/time/area.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';


@NgModule({
  declarations: [
    ProjectComponent,
    CountryComponent,
    SaveCountryComponent,
    EntityLifeComponent,
    SaveEntityLifeComponent,
    AreaComponent,
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
    CountryService,
    EntityLifeService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
