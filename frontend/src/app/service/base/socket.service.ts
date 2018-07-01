import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Observer }   from 'rxjs/Observer';
import * as socketIo  from 'socket.io-client';
import { Event }      from '../../model/base/socket';

@Injectable()
export class SocketService {
  private socket;

  public initSocket(url: string): any {
    return this.socket = socketIo(url);
  }

  public send(key: string, message: string): void {
    this.socket.emit(key, message);
  }

  public onMessage(key: string): Observable<any> {
      return new Observable<any>(observer => {
          this.socket.on('message', (data: any) => observer.next(data));
      });
  }

  public onEvent(event: Event): Observable<any> {
    return new Observable<Event>(observer => {
        this.socket.on(event, () => observer.next());
    });
  }
}
