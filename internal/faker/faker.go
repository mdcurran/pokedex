package faker

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mdcurran/pokedex/models"
)

type Faker struct {
	instance *gofakeit.Faker
}

func NewFaker() *Faker {
	return &Faker{instance: gofakeit.New(0)}
}

func (f *Faker) GenerateNature() *models.Nature {
	return &models.Nature{
		ID:   f.instance.Rand.Int(),
		Name: f.instance.Name(),
		HatesFlavor: &models.NamedApiResource{
			Name: f.instance.Name(),
			Url:  f.instance.URL(),
		},
	}
}
