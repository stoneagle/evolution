import { Injectable } from '@angular/core';
import { HttpHeaders  } from '@angular/common/http';
import { HttpEvent, HttpHandler, HttpRequest, HttpInterceptor  } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class CustomInterceptor implements HttpInterceptor {
    constructor() {
    }

    intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        request = request.clone({
            headers: new HttpHeaders({ 
              'Content-Type': 'application/json',
              'Access-Control-Allow-Origin': 'http://localhost:8080'
            }),
            withCredentials: true,
        });
        return next.handle(request);
    }
}
