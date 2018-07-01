import { NgModule }          from '@angular/core';
import { BrowserModule }     from '@angular/platform-browser';
import { CommonModule  }     from '@angular/common';
import { HttpClientModule  } from '@angular/common/http';
import { FormsModule }       from '@angular/forms';
import { ClarityModule  }    from "@clr/angular";
import { TranslateModule  }  from "@ngx-translate/core";

import { DefaultComponent }  from './default/default.component';
import { MessageComponent }  from './message/message.component';
import { SignInComponent }   from './sign-in/sign-in.component';
import { ShellComponent }    from './shell/shell.component';
import { ShellNavComponent } from './shell/nav/shell-nav.component';
import { AppRouteModule }    from '../route/app-route.module';
import { SignService  }      from '../service/base/sign.service';
import { MessageService }    from '../service/base/message.service';

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
