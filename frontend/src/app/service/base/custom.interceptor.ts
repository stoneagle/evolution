import { Injectable }                                            from '@angular/core';
import { HttpHeaders  }                                          from '@angular/common/http';
import { HttpEvent, HttpHandler, HttpRequest, HttpInterceptor  } from '@angular/common/http';
import { Observable }                                            from 'rxjs';
import { SignService }      from '../system/sign.service';

@Injectable()
export class CustomInterceptor implements HttpInterceptor {
  constructor(
    protected signService: SignService ,
  ) {
  }

    intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
      let headersMap = { 
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': 'http://localhost:8080'
      };
      if (request.url.indexOf("login") != -1) {
        headersMap = this.signService.setAuthHeader(headersMap);
      }
      request = request.clone({
          headers: new HttpHeaders(headersMap),
          withCredentials: true,
      });
      return next.handle(request);
    }
}
