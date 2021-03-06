import { InternationalConfig as N18 } from '../service/base/international.service';

export class ErrorInfo {
  public Server = {
    Error: {
      1: N18.settings.SYSTEM.ERROR.PARAMS,
      2: N18.settings.SYSTEM.ERROR.FILES,
      3: N18.settings.SYSTEM.ERROR.DATABASE,
      4: N18.settings.SYSTEM.ERROR.CACHE,
      5: N18.settings.SYSTEM.ERROR.AUTH,
      6: N18.settings.SYSTEM.ERROR.ENGINE,
      7: N18.settings.SYSTEM.ERROR.DATA_SERVICE,
      8: N18.settings.SYSTEM.ERROR.LOGIN,
      9: N18.settings.SYSTEM.ERROR.API,
      500: N18.settings.SYSTEM.ERROR.SERVER,
    },
    Exception: {
      NoResponse: N18.settings.SYSTEM.EXCEPTION.NORESPONSE,
    }
  }
  public Time = {
    Execing: N18.settings.TIME.ERROR.EXECING,
    Finished: N18.settings.TIME.ERROR.FINISHED,
    NotExec: N18.settings.TIME.ERROR.NOT_EXEC,
    NotFinish: N18.settings.TIME.ERROR.NOT_FINISH,
  }
}
