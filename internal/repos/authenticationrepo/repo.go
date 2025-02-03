package authenticationrepo

import (
	"context"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"time"

	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AuthenticationRepo struct {
	awsCognito ports.AWSCognito

	clientID string
}

func NewAuthenticationRepo(awsCognito ports.AWSCognito, clientID string) *AuthenticationRepo {
	return &AuthenticationRepo{
		awsCognito: awsCognito,
		clientID:   clientID,
	}
}

func (r *AuthenticationRepo) Authenticate(ctx context.Context, username, password string) (*domain.Token, error) {
	authResponse, err := r.awsCognito.InitiateAuth(ctx, &cognito.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       &r.clientID,
		AuthParameters: map[string]string{"USERNAME": username, "PASSWORD": password},
	})
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(time.Duration(authResponse.AuthenticationResult.ExpiresIn) * time.Second)

	return &domain.Token{
		TokenStr: *authResponse.AuthenticationResult.AccessToken,
		Expires:  expiration,
	}, nil

}
