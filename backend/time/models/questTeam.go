package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type QuestTeam struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int       `xorm:"unique not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	StartDate      time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"start_date,omitempty"`
	EndDate        time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"end_date,omitempty"`
	UserId         int       `xorm:"unique(quest_id) not null default 0 comment('成员id') INT(11)" structs:"user_id,omitempty"`
	Quest          *Quest    `xorm:"-" structs:"-" json:"Quest,omitempty"`
}

type QuestTeamJoin struct {
	QuestTeam `xorm:"extends" json:"-"`
	Quest     `xorm:"extends" json:"-"`
}

func NewQuestTeam() *QuestTeam {
	ret := QuestTeam{
		Quest: NewQuest(),
	}
	return &ret
}

func (m *QuestTeam) TableName() string {
	return "quest_team"
}

func (m *QuestTeam) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.Quest.BuildCondition(session)
	return
}

func (m *QuestTeam) SlicePtr() interface{} {
	ret := make([]QuestTeam, 0)
	return &ret
}

func (m *QuestTeam) Join() es.JoinGeneral {
	ret := QuestTeamJoin{}
	return &ret
}

func (m *QuestTeam) Transfer(slicePtr interface{}) *[]QuestTeam {
	ret := slicePtr.(*[]QuestTeam)
	return ret
}

func (j *QuestTeamJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	questLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Quest.TableName(),
		LeftTable:  j.Quest.TableName(),
		LeftField:  "id",
		RightTable: j.QuestTeam.TableName(),
		RightField: "quest_id",
	}
	links = append(links, questLink)
	return links
}

func (j *QuestTeamJoin) SlicePtr() interface{} {
	ret := make([]QuestTeamJoin, 0)
	return &ret
}

func (j *QuestTeamJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.QuestTeam
	ret.Quest = &join.Quest
	return &ret
}

func (j *QuestTeamJoin) TransferCopy(modelPtr es.ModelGeneral) {
	questTeamPtr := modelPtr.(*QuestTeam)
	(*questTeamPtr) = (*j).QuestTeam
	(*questTeamPtr).Quest = &(*j).Quest
	return
}

func (j *QuestTeamJoin) TransferCopySlice(slicePtr interface{}, targetPtr interface{}) {
	joinSlicePtr := slicePtr.(*[]QuestTeamJoin)
	joinSlice := *joinSlicePtr
	questTeams := make([]QuestTeam, 0)
	for _, one := range joinSlice {
		questTeamPtr := (&one).Transfer().(*QuestTeam)
		questTeams = append(questTeams, *questTeamPtr)
	}
	questTeamsPtr := targetPtr.(*[]QuestTeam)
	(*questTeamsPtr) = questTeams
	return
}
