/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package granthandlers

import (
	appmodel "github.com/asgardeo/thunder/internal/application/model"
	"github.com/asgardeo/thunder/internal/oauth/jwt"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/constants"
	"github.com/asgardeo/thunder/internal/oauth/oauth2/model"
)

// ClientCredentialsGrantHandler handles the client credentials grant type.
type ClientCredentialsGrantHandler struct{}

// ValidateGrant validates the client credentials grant type.
func (h *ClientCredentialsGrantHandler) ValidateGrant(tokenRequest *model.TokenRequest) *model.ErrorResponse {
	// Validate the grant type.
	if tokenRequest.GrantType != constants.GrantTypeClientCredentials {
		return &model.ErrorResponse{
			Error:            constants.ErrorUnsupportedGrantType,
			ErrorDescription: "Unsupported grant type",
		}
	}

	// Validate the client ID and secret.
	if tokenRequest.ClientID == "" || tokenRequest.ClientSecret == "" {
		return &model.ErrorResponse{
			Error:            constants.ErrorInvalidRequest,
			ErrorDescription: "Client Id and secret are required",
		}
	}

	return nil
}

// HandleGrant handles the client credentials grant type.
func (h *ClientCredentialsGrantHandler) HandleGrant(tokenRequest *model.TokenRequest,
	oauthApp *appmodel.OAuthApplication) (*model.TokenResponse, *model.ErrorResponse) {
	// Validate the client credentials.
	if tokenRequest.ClientID != oauthApp.ClientID || tokenRequest.ClientSecret != oauthApp.ClientSecret {
		return nil, &model.ErrorResponse{
			Error:            constants.ErrorInvalidClient,
			ErrorDescription: "Invalid client credentials",
		}
	}

	// Generate a JWT token for the client.
	token, err := jwt.GenerateJWT(tokenRequest.ClientID, tokenRequest.ClientID, nil)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:            constants.ErrorServerError,
			ErrorDescription: "Failed to generate token",
		}
	}

	// Return the token response.
	return &model.TokenResponse{
		AccessToken: token,
		TokenType:   constants.TokenTypeBearer,
		Scope:       tokenRequest.Scope,
		ExpiresIn:   3600,
	}, nil
}
