import { Component, Input, OnInit, OnDestroy, ElementRef } from '@angular/core';
import { Router }                                          from '@angular/router';
import { Subscription }                                    from 'rxjs/Subscription';
import { TranslateService }                                from '@ngx-translate/core';

import { Message }                                                  from '../../model/base/message';
import { MessageService }                                           from '../../service/base/message.service';
import { AlertType, dismissInterval, httpStatusCode, CommonRoutes } from '../../shared/shared.const';

@Component({
  selector: 'global-message',
  templateUrl: 'message.component.html',
  styleUrls: ['message.component.css']
})
export class MessageComponent implements OnInit, OnDestroy {

  @Input() isAppLevel: boolean;
  globalMessage: Message = new Message();
  globalMessageOpened: boolean;
  messageText: string = "";
  timer: any = null;
  
  appLevelMsgSub: Subscription;
  msgSub: Subscription;
  clearSub: Subscription;

  constructor(
    private elementRef: ElementRef,
    private messageService: MessageService,
    private router: Router,
    private translate: TranslateService) { }

  ngOnInit(): void {
    if (this.isAppLevel) {
      this.appLevelMsgSub = this.messageService.appLevelAnnounced$.subscribe(
        message => {
          this.globalMessageOpened = true;
          this.globalMessage = message;
          this.messageText = message.message;
          this.translateMessage(message);
        }
      )
    } else {
      this.msgSub = this.messageService.messageAnnounced$.subscribe(
        message => {
          this.globalMessageOpened = true;
          this.globalMessage = message;
          this.messageText = message.message;
          this.translateMessage(message);
          this.timer = setTimeout(() => this.onClose(), dismissInterval);

          setTimeout(() => {
            let nativeDom: any = this.elementRef.nativeElement;
            let queryDoms: any[] = nativeDom.getElementsByClassName("alert");
            if (queryDoms && queryDoms.length > 0) {
              let hackDom: any = queryDoms[0];
              hackDom.className += ' alert-global alert-global-align';
            }
          }, 0);
        }
      );
    }

    this.clearSub = this.messageService.clearChan$.subscribe(clear => {
      this.onClose();
    });
  }

  ngOnDestroy() {
    if (this.appLevelMsgSub) {
      this.appLevelMsgSub.unsubscribe();
    }

    if (this.msgSub) {
      this.msgSub.unsubscribe();
    }

    if (this.clearSub) {
      this.clearSub.unsubscribe();
    }
  }

  translateMessage(msg: Message): void {
    let key = "UNKNOWN_ERROR", param = "";
    if (msg && msg.message) {
      key = (typeof msg.message === "string" ? msg.message.trim() : msg.message);
      if (key === "") {
        key = "UNKNOWN_ERROR";
      }
    }
    this.translate.get(key, { 'param': param }).subscribe((res: string) => this.messageText = res);
  }

  public get needAuth(): boolean {
    return this.globalMessage ?
      this.globalMessage.statusCode === httpStatusCode.Unauthorized : false;
  }

  public get message(): string {
    return this.messageText;
  }

  signIn(): void {
    this.router.navigateByUrl(CommonRoutes.SIGN_IN);
  }

  onClose() {
    if (this.timer) {
      clearTimeout(this.timer);
    }
    this.globalMessageOpened = false;
  }
}
