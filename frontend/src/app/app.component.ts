import { Component }         from '@angular/core';
import { TranslateService  } from '@ngx-translate/core';
import { AppConfig }         from './service/base/config.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';

  constructor(
    translate: TranslateService
  ) {
    let lang = AppConfig.settings.app.language;
    translate.setDefaultLang(lang);
    translate.use(lang);
  }
}
