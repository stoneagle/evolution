import { NgModule } from '@angular/core';

import { PoolComponent } from './pool/pool.component';
import { SavePoolComponent } from './pool/save/save.component';
import { PoolService } from '../service/business/pool.service';
import { MessageHandlerService  } from '../service/base/message-handler.service';

import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule   } from '@angular/platform-browser/animations';
import { ClarityModule  } from "@clr/angular"; 
import { HttpClientModule  }    from '@angular/common/http';
import { HttpModule   } from '@angular/http';
import { FormsModule }   from '@angular/forms';
import { TranslateModule  } from "@ngx-translate/core";

@NgModule({
  declarations: [
    PoolComponent,
    SavePoolComponent,
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
    PoolService,
    MessageHandlerService,
  ],
  exports: [
    PoolComponent,
  ]
})

export class ControllerModule {
}
