package ports

import (
	"context"
	"github.com/Braly-Ltd/t2v-api-core/entities"
)

type AuthenticationPort interface {
	Authenticate(ctx context.Context, token string) (entities.TokenData, error)
}
