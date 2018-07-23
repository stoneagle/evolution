import { InternationalConfig as N18 } from '../service/base/international.service';

export class ShareSettings {
  public Time = {
    Resource: {
      Country: N18.settings.TIME.RESOURCE.COUNTRY.CONCEPT,
      Area: N18.settings.TIME.RESOURCE.AREA.CONCEPT,
      Resource: N18.settings.TIME.RESOURCE.RESOURCE.CONCEPT,
      Field: N18.settings.TIME.RESOURCE.FIELD.CONCEPT,
      Project: N18.settings.TIME.RESOURCE.PROJECT.CONCEPT,
      Task: N18.settings.TIME.RESOURCE.TASK.CONCEPT,
      Action: N18.settings.TIME.RESOURCE.ACTION.CONCEPT,
      Phase: N18.settings.TIME.RESOURCE.PHASE.CONCEPT,
      Quest: N18.settings.TIME.RESOURCE.QUEST.CONCEPT,
      QuestTeam: N18.settings.TIME.RESOURCE.QUEST_TEAM.CONCEPT,
      QuestTimeTable: N18.settings.TIME.RESOURCE.QUEST_TIMETABLE.CONCEPT,
      QuestTarget: N18.settings.TIME.RESOURCE.QUEST_TARGET.CONCEPT,
      QuestResource: N18.settings.TIME.RESOURCE.QUEST_RESOURCE.CONCEPT,
      UserResource: N18.settings.TIME.RESOURCE.USER_RESOURCE.CONCEPT,
      General: N18.settings.TIME.RESOURCE.GENERAL.CONCEPT,
    }
  }
  public Quant = {
    Resource: {
      Pool: N18.settings.QUANT.RESOURCE.POOL.CONCEPT,
      Item: N18.settings.QUANT.RESOURCE.ITEM.CONCEPT,
      Classify: N18.settings.QUANT.RESOURCE.CLASSIFY.CONCEPT,
      Config: N18.settings.QUANT.RESOURCE.CONFIG.CONCEPT,
      General: N18.settings.SYSTEM.RESOURCE.GENERAL.CONCEPT,
    }
  }
  public System = {
    Resource: {
      User: N18.settings.SYSTEM.RESOURCE.USER.CONCEPT,
      General: N18.settings.SYSTEM.RESOURCE.GENERAL.CONCEPT,
    },
    Process: {
      Signin: N18.settings.SYSTEM.PROCESS.SIGNIN,
      Signout: N18.settings.SYSTEM.PROCESS.SIGNOUT,
      Get: N18.settings.SYSTEM.PROCESS.GET,
      List: N18.settings.SYSTEM.PROCESS.LIST,
      Sync: N18.settings.SYSTEM.PROCESS.SYNC,
      Create: N18.settings.SYSTEM.PROCESS.CREATE,
      Update: N18.settings.SYSTEM.PROCESS.UPDATE,
      Delete: N18.settings.SYSTEM.PROCESS.DELETE,
      Finish: N18.settings.SYSTEM.PROCESS.FINISH,
      Back: N18.settings.SYSTEM.PROCESS.BACK,
      Next: N18.settings.SYSTEM.PROCESS.NEXT,
      Submit: N18.settings.SYSTEM.PROCESS.SUBMIT,
      Cancel: N18.settings.SYSTEM.PROCESS.CANCEL,
      Rollback: N18.settings.SYSTEM.PROCESS.ROLLBACK,
    }
  }
}

