import { NgModule } from '@angular/core';
import { AppRouteModule } from '../route/app-route.module';
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule  } from '@angular/common';   
import { ClarityModule  } from "@clr/angular"; 
import { HttpClientModule  }    from '@angular/common/http';
import { FormsModule }   from '@angular/forms';
import { TranslateModule  } from "@ngx-translate/core";

import { SignInComponent } from './sign-in/sign-in.component';
import { ShellComponent } from './shell/shell.component';
import { DefaultComponent } from './default/default.component';
import { ShellNavComponent } from './shell/nav/shell-nav.component';
import { MessageComponent } from './message/message.component';
import { SignService  } from '../service/base/sign.service';
import { MessageService } from '../service/base/message.service';

@NgModule({
  declarations: [
    SignInComponent,
    MessageComponent,
    ShellComponent,
    DefaultComponent,
    ShellNavComponent,
  ],
  imports: [   
    AppRouteModule,
    BrowserModule,
    HttpClientModule,
    FormsModule,
    ClarityModule,
    TranslateModule,
  ],
  providers: [ 
    SignService,
    MessageService,
  ],
  exports: [
    SignInComponent,
  ]
})

export class BaseModule {
}
