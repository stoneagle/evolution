package models

type ModelPackage struct {
	ActionModel         *Action
	AreaModel           *Area
	CountryModel        *Country
	FieldModel          *Field
	PhaseModel          *Phase
	ProjectModel        *Project
	QuestModel          *Quest
	QuestResourceModel  *QuestResource
	QuestTargetModel    *QuestTarget
	QuestTeamModel      *QuestTeam
	QuestTimeTableModel *QuestTimeTable
	ResourceModel       *Resource
	TaskModel           *Task
	UserResourceModel   *UserResource
}

func (m *ModelPackage) PrepareModel() {
	m.ActionModel = &Action{}
	m.AreaModel = &Area{}
	m.CountryModel = &Country{}
	m.FieldModel = &Field{}
	m.PhaseModel = &Phase{}
	m.ProjectModel = &Project{}
	m.QuestModel = &Quest{}
	m.QuestResourceModel = &QuestResource{}
	m.QuestTargetModel = &QuestTarget{}
	m.QuestTeamModel = &QuestTeam{}
	m.QuestTimeTableModel = &QuestTimeTable{}
	m.ResourceModel = &Resource{}
	m.TaskModel = &Task{}
	m.UserResourceModel = &UserResource{}
}
