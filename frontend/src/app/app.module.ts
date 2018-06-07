import { APP_INITIALIZER } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule  }    from '@angular/forms';
import { HttpClientModule  }    from '@angular/common/http';

import { ClarityModule  } from "@clr/angular";
import { AppComponent } from './app.component';
import { PoolComponent } from './controller/pool/pool.component';
import { StockComponent } from './controller/stock/stock.component';
import { AppRouteModule } from './route/app-route.module';

import { CustomInterceptor  } from './service/base/custom.interceptor';
import { HTTP_INTERCEPTORS   } from '@angular/common/http';

@NgModule({
  declarations: [
    AppComponent,
    PoolComponent,
    StockComponent
  ],
  imports: [
    ClarityModule,
    FormsModule,
    AppRouteModule,
    HttpClientModule,
    BrowserModule
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: CustomInterceptor ,
      multi: true
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
