package fixture

import "tech-tsarka/internal/storage/user/entity"

type UserBuilder struct {
	instance *entity.User
}

func User() *UserBuilder {
	return &UserBuilder{
		instance: &entity.User{},
	}
}

func (b *UserBuilder) ID(v string) *UserBuilder {
	b.instance.ID = v
	return b
}

func (b *UserBuilder) FirstName(v string) *UserBuilder {
	b.instance.FirstName = v
	return b
}

func (b *UserBuilder) LastName(v string) *UserBuilder {
	b.instance.LastName = v
	return b
}

func (b *UserBuilder) P() *entity.User {
	return b.instance
}

func (b *UserBuilder) V() entity.User {
	return *b.instance
}
