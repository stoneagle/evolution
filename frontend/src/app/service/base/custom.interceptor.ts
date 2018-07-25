import { Injectable }                                            from '@angular/core';
import { HttpHeaders  }                                          from '@angular/common/http';
import { HttpEvent, HttpHandler, HttpRequest, HttpInterceptor  } from '@angular/common/http';
import { Observable }                                            from 'rxjs';
import { SignService }                                           from '../system/sign.service';

@Injectable()
export class CustomInterceptor implements HttpInterceptor {
  constructor(
    protected signService: SignService,
  ) {
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let headersMap = { 
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': 'http://localhost:8080'
    };

    let token = this.signService.getAuthToken();
    if (token != undefined) {
      if (token.length > 0) {
        headersMap["Authorization"] = token;
      }
    }

    request = request.clone({
      setHeaders: headersMap,
      withCredentials: true,
    });
    return next.handle(request);
  }
}
