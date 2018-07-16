import { Injectable }     from '@angular/core';
import { Http, Response } from '@angular/http';
import { environment }    from '../../../environments/environment';
import { International }  from '../../model/base/international';
// import { AppConfig }      from './config.service';

@Injectable()  
export class InternationalConfig {
  static settings: International;
  constructor(private http: Http) {
  }
  load() {   
    // const jsonFile = `assets/i18n/${AppConfig.settings.app.language}.json`;
    const jsonFile = `assets/i18n/${environment.language}.json`;
    return new Promise<void>((resolve, reject) => {
      this.http.get(jsonFile).toPromise().then((response : Response) => {
        InternationalConfig.settings = <International>response.json();
        resolve();
      }).catch((response: any) => {
        reject(`Could not load file '${jsonFile}': ${JSON.stringify(response)}`);
      });
    });
  }
}
