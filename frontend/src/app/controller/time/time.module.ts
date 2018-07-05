import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";
import { TreeModule }                from 'ng2-tree';

import { ProjectComponent }         from './project/project.component';
import { ProjectService }           from '../../service/time/project.service';
import { CountryComponent }         from './country/country.component';
import { SaveCountryComponent }     from './country/save/save.component';
import { EntityComponent }         from './entity/entity.component';
import { SaveEntityComponent }     from './entity/save/save.component';
import { FieldComponent }           from './field/field.component';
import { SaveFieldComponent }       from './field/save/save.component';
import { AreaComponent }            from './area/area.component';

import { CountryService }         from '../../service/time/country.service';
import { FieldService }           from '../../service/time/field.service';
import { EntityService }          from '../../service/time/entity.service';
import { AreaService }            from '../../service/time/area.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';


@NgModule({
  declarations: [
    ProjectComponent,
    CountryComponent,
    SaveCountryComponent,
    EntityComponent,
    SaveEntityComponent,
    FieldComponent,
    SaveFieldComponent,
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
    FieldService,
    EntityService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
