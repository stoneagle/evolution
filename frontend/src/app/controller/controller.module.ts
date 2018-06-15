import { NgModule } from '@angular/core';

import { PoolComponent } from './pool/pool.component';
import { SavePoolComponent } from './pool/save/save.component';
import { ClassifyComponent } from './classify/classify.component';
import { SyncClassifyComponent } from './classify/sync/sync.component';
import { ItemComponent } from './item/item.component';
import { SyncItemComponent } from './item/sync/sync.component';
import { PoolService } from '../service/business/pool.service';
import { ConfigService } from '../service/config/config.service';
import { ClassifyService } from '../service/config/classify.service';
import { ItemService } from '../service/business/item.service';
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
    ClassifyComponent,
    SyncClassifyComponent,
    ItemComponent,
    SyncItemComponent,
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
    ConfigService,
    ClassifyService,
    ItemService,
    MessageHandlerService,
  ],
  exports: [
    PoolComponent,
  ]
})

export class ControllerModule {
}
