import { Injectable } from '@angular/core'
import { Subject } from 'rxjs/Subject';
import { Router } from "@angular/router";

import { MessageService } from './message.service';
import { AlertType, httpStatusCode } from '../../shared/shared.const';
import { TranslateService } from '@ngx-translate/core';

@Injectable()
export class MessageHandlerService {

    constructor(
      private msgService: MessageService,
      private translate: TranslateService,
      private router: Router
    ) { }

    public handleError(error: any | string): void {
        if (!error) {
            return;
        }
        let msg = this.errorHandler(error);

        if (!(error.statusCode || error.status)) {
            this.msgService.announceMessage(500, msg, AlertType.DANGER);
        } else {
            let code = error.statusCode | error.status;
            if (code === httpStatusCode.Unauthorized) {
              this.router.navigate(['/login']);
            } else {
              this.msgService.announceMessage(code, msg, AlertType.DANGER);
            }
        }
    }

    public handleReadOnly(): void {
        this.msgService.announceAppLevelMessage(503, 'REPO_READ_ONLY', AlertType.WARNING);
    }

    public showError(message: string, params: any): void {
        if (!params) {
            params = {};
        }
        this.translate.get(message, params).subscribe((res: string) => {
            this.msgService.announceMessage(500, res, AlertType.DANGER);
        });
    }

    public showSuccess(message: string): void {
        if (message && message.trim() != "") {
            this.msgService.announceMessage(200, message, AlertType.SUCCESS);
        }
    }

    public showInfo(message: string): void {
        if (message && message.trim() != "") {
            this.msgService.announceMessage(200, message, AlertType.INFO);
        }
    }

    public showWarning(message: string): void {
        if (message && message.trim() != "") {
            this.msgService.announceMessage(400, message, AlertType.WARNING);
        }
    }

    public clear(): void {
        this.msgService.clear();
    }

    public isAppLevel(error: any): boolean {
        return error && error.statusCode === httpStatusCode.Unauthorized;
    }

    public error(error: any): void {
        this.handleError(error);
    }

    public warning(warning: any): void {
        this.showWarning(warning);
    }

    public info(info: any): void {
        this.showSuccess(info);
    }

    public log(log: any): void {
        this.showInfo(log);
    }

    public errorHandler(error: any): string {
        if (!error) {
            return "UNKNOWN_ERROR";
        }
        if (!(error.statusCode || error.status)) {
            return '' + error;
        } else {
            switch (error.statusCode || error.status) {
                case 400:
                    return "BAD_REQUEST_ERROR";
                case 401:
                    return "UNAUTHORIZED_ERROR";
                case 403:
                    return "FORBIDDEN_ERROR";
                case 404:
                    return "NOT_FOUND_ERROR";
                case 412:
                    return "PRECONDITION_FAILED";
                case 409:
                    return "CONFLICT_ERROR";
                case 500:
                    return "SERVER_ERROR";
                default:
                    return "UNKNOWN_ERROR";
            }
        }
    }

}
