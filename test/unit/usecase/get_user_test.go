package usecase_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	users, err := userRepo.SearchByUsernameOrEmail("azis")
	assert.Nil(t, err)
	fmt.Println(users)
}
