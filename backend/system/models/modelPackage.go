package models

type ModelPackage struct {
	UserModel *User
}

func (m *ModelPackage) PrepareModel() {
	m.UserModel = &User{}
}
