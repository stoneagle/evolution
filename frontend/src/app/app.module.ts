import { APP_INITIALIZER } from '@angular/core';
import { NgModule } from '@angular/core';
import { HttpClient  } from '@angular/common/http';
import { TranslateModule, TranslateLoader  } from '@ngx-translate/core';
import { TranslateHttpLoader  } from '@ngx-translate/http-loader';
import { HTTP_INTERCEPTORS   } from '@angular/common/http';
import { CustomInterceptor  } from './service/base/custom.interceptor';

import { AppComponent } from './app.component';
import { AppRouteModule } from './route/app-route.module';
import { AppConfig  } from './service/base/config.service';
import { ClarityModule  } from "@clr/angular";
import { BaseModule } from './base/base.module';
import { ControllerModule } from './controller/controller.module';

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
    ControllerModule,
    BaseModule,
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: (createTranslateLoader),
        deps: [HttpClient]
      }
    })
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
