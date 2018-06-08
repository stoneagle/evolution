import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { CommonModule  } from '@angular/common';   
import { ClarityModule  } from "@clr/angular"; 
import { HttpClientModule  }    from '@angular/common/http';
import { FormsModule }   from '@angular/forms';
import { PoolComponent } from './pool/pool.component';
import { StockComponent } from './stock/stock.component';
import { MessageService } from '../service/base/message.service';
import { MessageHandlerService  } from '../service/base/message-handler.service';

@NgModule({
  declarations: [
    PoolComponent,
    StockComponent,
  ],
  imports: [   
    BrowserModule,
    HttpClientModule,
    FormsModule,
    ClarityModule,
  ],
  providers: [ 
    MessageService,
    MessageHandlerService,
  ],
  exports: [
  ]
})

export class ControllerModule {
}
