export const enum AlertType {
  DANGER, WARNING, INFO, SUCCESS
}

export const dismissInterval = 10 * 1000;
export const httpStatusCode = {
  "Unauthorized": 401,
}

export const AuthType = {
  BasicAuth: "basic-auth",
  BasicAuthJwt: "basic-auth-jwt",
}

export const AreaType = {
  Root: 1,
  Node: 2,
  Leaf: 3,
}

export const Task = {
  Status: {
    Backlog: 1,
    Todo: 2,
    Progress: 3,
    Done: 4
  },
  StatusName: {
    1: "backlog",
    2: "todo",
    3: "progress",
    4: "done",
  },
  StatusInfo: {
    1: "TIME.RESOURCE.TASK.STATUS-NAME.BACKLOG",
    2: "TIME.RESOURCE.TASK.STATUS-NAME.TODO",
    3: "TIME.RESOURCE.TASK.STATUS-NAME.PROGRESS",
    4: "TIME.RESOURCE.TASK.STATUS-NAME.DONE",
  }
}

export const Quest = {
  Status: {
    Recruit: 1,
    Exec: 2,
    Finish: 3,
    Fail: 4,
  },
  StatusInfo: {
    1: "TIME.RESOURCE.QUEST.STATUS-NAME.RECRUIT",
    2: "TIME.RESOURCE.QUEST.STATUS-NAME.EXEC",
    3: "TIME.RESOURCE.QUEST.STATUS-NAME.FINISH",
    4: "TIME.RESOURCE.QUEST.STATUS-NAME.FAIL",
  },
  Members: {
    One: 1,
    Small: 5,
    Middle: 25,
    Large: 100,
  },
  MembersInfo: {
    1: "TIME.RESOURCE.QUEST.MEMBERS-NAME.ONE",
    5: "TIME.RESOURCE.QUEST.MEMBERS-NAME.SMALL",
    25: "TIME.RESOURCE.QUEST.MEMBERS-NAME.MIDDLE",
    100: "TIME.RESOURCE.QUEST.MEMBERS-NAME.LARGE",
  },
  Constraint: {
    ImportantAndBusy: 1,
    Important: 2,
    Busy: 3,
    Normal: 4,
  },
  ConstraintInfo: {
    1:"TIME.RESOURCE.QUEST.CONSTRAINT-NAME.IMPORTANT-BUSY",
    2:"TIME.RESOURCE.QUEST.CONSTRAINT-NAME.IMPORTANT",
    3:"TIME.RESOURCE.QUEST.CONSTRAINT-NAME.BUSY",
    4:"TIME.RESOURCE.QUEST.CONSTRAINT-NAME.NORMAL",
  },
  TimeTableType: {
    Workday: 1,
    holiday: 2,
  },
  TimeTableTypeInfo : {
    1: "TIME.RESOURCE.QUEST-TIMETABLE.TYPE-NAME.WORKDAY",
    2: "TIME.RESOURCE.QUEST-TIMETABLE.TYPE-NAME.HOLIDAY",
  },
  TargetStatus: {
    Wait: 1,
    Finish: 2,
  },
  TargetStatusInfo : {
    1: "TIME.RESOURCE.QUEST-TARGET.STATUS-NAME.WAIT",
    2: "TIME.RESOURCE.QUEST-TARGET.STATUS-NAME.FINISH",
  },
  ResourceStatus: {
    Unmatch: 1,
    Matched: 2,
  },
  ResourceStatusInfo : {
    1: "TIME.RESOURCE.QUEST-RESOURCE.STATUS-NAME.UNMATCH",
    2: "TIME.RESOURCE.QUEST-RESOURCE.STATUS-NAME.MATCHED",
  }
}

export const WsStatus = {
  Message: 0,
  Connect: 1,
  Disconnect: 2,
  Error: 3,
}

export const CommonRoutes = {  
  SIGN_IN: "/login",
}

