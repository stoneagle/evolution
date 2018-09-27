import { NgModule }                  from '@angular/core';
import { BrowserModule }             from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { HttpClientModule  }         from '@angular/common/http';
import { HttpModule   }              from '@angular/http';
import { FormsModule }               from '@angular/forms';
import { ClarityModule  }            from "@clr/angular";
import { TranslateModule  }          from "@ngx-translate/core";

import { UserComponent }     from './user/user.component';
import { UserSaveComponent } from './user/save/save.component';

import { UserService }         from '../../service/system/user.service';


@NgModule({
  declarations: [
    UserComponent,
    UserSaveComponent,
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
    UserService,
  ],
  exports: [
  ]
})

export class SystemModule {
}
