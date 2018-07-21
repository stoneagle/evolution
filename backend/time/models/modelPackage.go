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
	m.ActionModel = NewAction()
	m.AreaModel = NewArea()
	m.CountryModel = NewCountry()
	m.FieldModel = NewField()
	m.PhaseModel = NewPhase()
	m.ProjectModel = NewProject()
	m.QuestModel = NewQuest()
	m.QuestResourceModel = NewQuestResource()
	m.QuestTargetModel = NewQuestTarget()
	m.QuestTeamModel = NewQuestTeam()
	m.QuestTimeTableModel = NewQuestTimeTable()
	m.ResourceModel = NewResource()
	m.TaskModel = NewTask()
	m.UserResourceModel = NewUserResource()
}
