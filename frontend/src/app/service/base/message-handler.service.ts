import { Injectable }                from '@angular/core'
import { Subject }                   from 'rxjs/Subject';
import { Router }                    from "@angular/router";
import { TranslateService }          from '@ngx-translate/core';
import { MessageService }            from './message.service';
import { AlertType, httpStatusCode } from '../../shared/const';

@Injectable()
export class MessageHandlerService {
  constructor(
    private msgService: MessageService,
    private translateService: TranslateService,
    private router: Router
  ) { }

  public handleError(resource: string, operation:string, error: any | string): void {
    if (!error) {
      return;
    }
    let msg = error.message;
    let code = error.statusCode | error.status;
    if (code === httpStatusCode.Unauthorized) {
      this.router.navigate(['/sign/login']);
    } else {
      this.SendMessage(resource, operation, AlertType.DANGER, 500, msg);
    }
  }

  public showError(resource: string, operation:string, message?: string): void {
    if (message && message.trim() != "") {
      this.SendMessage(resource, operation, AlertType.DANGER, 500, message);
    } else {
      this.SendMessage(resource, operation, AlertType.DANGER, 500);
    }
  }

  public showSuccess(resource: string, operation:string, message?: string): void {
    if (message && message.trim() != "") {
      this.SendMessage(resource, operation, AlertType.SUCCESS, 200, message);
    } else {
      this.SendMessage(resource, operation, AlertType.SUCCESS, 200);
    }
  }

  public showInfo(resource: string, operation:string, message?: string): void {
    if (message && message.trim() != "") {
      this.SendMessage(resource, operation, AlertType.INFO, 200, message);
    } else {
      this.SendMessage(resource, operation, AlertType.INFO, 200);
    }
  }

  public showWarning(resource: string, operation:string, message?: string): void {
    if (message && message.trim() != "") {
      this.SendMessage(resource, operation, AlertType.WARNING, 400, message);
    } else {
      this.SendMessage(resource, operation, AlertType.WARNING, 400);
    }
  }

  public clear(): void {
    this.msgService.clear();
  }

  public isAppLevel(error: any): boolean {
    return error && error.statusCode === httpStatusCode.Unauthorized;
  }

  public errorHandler(error: any): string {
    if (!error) {
      return "UNKNOWN_ERROR";
    }
  }

  private SendMessage(resource: string, operation: string, alertType: AlertType, code: number, message?: string): void {
    let result: string;
    switch (alertType) {
      case AlertType.INFO:
        result = "SYSTEM.RESULT.EXEC";
        break;
      case AlertType.SUCCESS:
        result = "SYSTEM.RESULT.SUCCESS";
        break;
      case AlertType.WARNING:
        result = "SYSTEM.RESULT.FAIL";
        break;
      case AlertType.DANGER:
        result = "SYSTEM.RESULT.EXCEPTION";
        break;
      default:
        result = "SYSTEM.RESULT.UNKNOWN";
        break;
    }

    this.translateService.get([resource, operation, result], { 'param': '' }).subscribe((res: any) => {
      let msg: string;
      if (message && message.trim() != "") {
        msg = res[resource] + res[operation] + res[result] + `【${message}】`;
      } else {
        msg = res[resource] + res[operation] + res[result];
      }
      this.msgService.announceMessage(code, msg, alertType);
    });
  }
}
