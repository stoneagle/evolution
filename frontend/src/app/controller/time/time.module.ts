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
import { MessageHandlerService  } from '../../service/base/message-handler.service';


@NgModule({
  declarations: [
    ProjectComponent,
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
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class TimeModule {
}
