import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";

import { ProjectComponent }         from './project/project.component';
import { ProjectService }            from '../../service/time/project.service';
import { CountryComponent }         from './country/country.component';
import { SaveCountryComponent }         from './country/save/save.component';
import { CountryService }            from '../../service/time/country.service';
import { MessageHandlerService  } from '../../service/base/message-handler.service';


@NgModule({
  declarations: [
    ProjectComponent,
    CountryComponent,
    SaveCountryComponent,
  ],
  imports: [   
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    HttpModule,
    FormsModule,
    ClarityModule,
    TranslateModule,
  ],
  providers: [ 
    ProjectService,
    CountryService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
