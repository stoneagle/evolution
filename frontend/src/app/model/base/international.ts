export interface International {
  APP: {
    TITLE: string
  }
  TIME: {
    TITLE: string
    ASSETS: {
      TYPE: {
        SCHEDULE: string
        USER: string
        CONFIG: string
      }
    }
    ERROR: {
      TARGET_NOT_FINISH: string 
    }
    RESOURCE: {
      COUNTRY: {
        CONCEPT: string
        NAME: string
        ENNAME: string
      }
      AREA: {
        CONCEPT: string
        NAME: string
        PARENT: string
        FIELDID: string
      }
      RESOURCE: {
        CONCEPT: string
        NAME: string
        DESC: string
        YEAR: string
        AREA: string
      }
      FIELD: {
        CONCEPT: string
        NAME: string
        DESC: string
        COLOR: string
      }
      PROJECT: {
        CONCEPT: string
        NAME: string
        QUEST: string
        AREA: string
        STARTDATE: string
        DURATION: string
      }
      TASK: {
        CONCEPT: string
        NAME: string
        DESC: string
        PROJECT: string
        RESOURCE: string
        USER: string
        STARTDATE: string
        ENDDATE: string
        DURATION: string
        STATUS: string
        STATUS_NAME: {
          BACKLOG: string
          TODO: string
          PROGRESS: string
          DONE: string        
        }
      }
      ACTION: {
        CONCEPT: string
        NAME: string
        STARTDATE: string 
        ENDDATE: string
        TASK: string
        USER: string
        TIME: string
      }
      PHASE: {
        CONCEPT: string
        NAME: string
        DESC: string
        LEVEL: string
        THRESHOLD: string
        FIELD: string
      }
      QUEST: {
        CONCEPT: string
        NAME: string
        STARTDATE: string
        ENDDATE: string
        FOUNDER: string
        MEMBERS: string
        CONSTRAINT: string
        STATUS: string
        STATUS_NAME: {
          RECRUIT: string
          EXEC: string
          FINISH: string
          FAIL: string
        }
        MEMBERS_NAME: {
          ONE: string
          SMALL: string
          MIDDLE: string
          LARGE: string
        }
        CONSTRAINT_NAME: {
          IMPORTANT_BUSY: string
          IMPORTANT: string
          BUSY: string
          NORMAL: string
        }
      }
      QUEST_TEAM: {
        CONCEPT: string
        QUEST: string
        STARTDATE: string
        ENDDATE: string
        USER: string
      }
      QUEST_TIMETABLE: {
        CONCEPT: string
        QUEST: string
        STARTTIME: string
        ENDDATE: string
        TYPE: string
        TYPE_NAME: {
          WORKDAY: string
          HOLIDAY: string
        } 
      }
      QUEST_TARGET: {
        CONCEPT: string
        QUEST: string
        RESOURCE: string
        DESC: string
        STATUS: string
        STATUS_NAME: {
          WAIT: string
          FINISH: string        
        } 
      }
      QUEST_RESOURCE: {
        CONCEPT: string
        QUEST: string
        RESOURCE: string
        DESC: string
        NUMBER: string
        STATUS: string
        STATUS_NAME: {
          UNMATCH: string
          MATCHED: string
        } 
      }
      USER_RESOURCE: {
        CONCEPT: string
        USER: string
        TIME: string
      }
      GENERAL: {
        CONCEPT: string
        ID: string
        OPERATION: string
        CREATED: string
        UPDATED: string
        DELETED: string
      }
    }
  }
  QUANT: {
    TITLE: string
    ASSETS: {
      TYPE: {
        STOCK: string
        EXCHANGE: string
        CONFIG: string
      }
    }
    RESOURCE: {
      POOL: {
        CONCEPT: string
        NAME: string
        STATUS: string
        STRATEGY: string
        TYPE: string
      }
      ITEM: {
        CONCEPT: string
        CODE: string
        NAME: string
        STATUS: string
      }
      CLASSIFY: {
        CONCEPT: string
        NAME: string
        TAG: string
      }
      CONFIG: {
        CONCEPT: string
        ASSET: string
        TYPE: string
        MAIN: string
        SUB: string
      }
      GENERAL: {
        CONCEPT: string
        ID: string
        OPERATION: string
        CREATED: string
        UPDATED: string
        DELETED: string
      }
    }
  }
  SYSTEM: {
    TITLE: string
    ASSETS: {
      TYPE: {
        USER: string
      }
    }
    TOOLTIP: {
      EMPTY: string
    }
    RESOURCE: {
      USER: {
        CONCEPT: string
        NAME: string
        EMAIL: string
        PASSWORD: string
      }
      GENERAL: {
        CONCEPT: string
        ID: string
        OPERATION: string
        CREATED: string
        UPDATED: string
        DELETED: string
      }
    }
    PROCESS: { 
      SIGNIN: string
      SIGNOUT: string
      GET: string 
      LIST: string 
      SYNC: string 
      CREATE: string 
      UPDATE: string
      DELETE: string
      BACK: string
      NEXT: string
      SUBMIT: string
      CANCEL: string
      ROLLBACK: string
    }
    RESULT: {
      SUCCESS: string
      EXEC: string
      FAIL: string
      EXCEPTION: string
      UNKNOWN: string
    }
  }
}
