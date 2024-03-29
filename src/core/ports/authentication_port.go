package ports

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-core/entities"
)

type AuthenticationPort interface {
	Authenticate(ctx context.Context, agent, token string) (entities.TokenData, error)
}
