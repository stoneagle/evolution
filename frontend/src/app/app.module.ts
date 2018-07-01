import { APP_INITIALIZER }                   from '@angular/core';
import { NgModule }                          from '@angular/core';
import { HttpClient, HTTP_INTERCEPTORS   }   from '@angular/common/http';
import { TranslateModule, TranslateLoader  } from '@ngx-translate/core';
import { TranslateHttpLoader  }              from '@ngx-translate/http-loader';
import { ClarityModule  }                    from "@clr/angular";

import { AppComponent }       from './app.component';
import { AppRouteModule }     from './route/app-route.module';
import { AppConfig  }         from './service/base/config.service';
import { CustomInterceptor  } from './service/base/custom.interceptor';
import { BaseModule }         from './base/base.module';
import { QuantModule }        from './controller/quant/quant.module';
import { TimeModule }        from './controller/time/time.module';

export function initializeApp(appConfig: AppConfig) {
  return () => appConfig.load();
}

export function createTranslateLoader(http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    ClarityModule,
    AppRouteModule,
    BaseModule,
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: (createTranslateLoader),
        deps: [HttpClient]
      }
    }),
    QuantModule,
    TimeModule,
  ],
  providers: [
    AppConfig,
    { 
      provide: APP_INITIALIZER,
      useFactory: initializeApp,
      deps: [AppConfig],
      multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: CustomInterceptor ,
      multi: true
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
