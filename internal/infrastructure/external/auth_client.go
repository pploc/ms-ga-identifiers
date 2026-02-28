package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

type RolePermission struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type AuthClient struct {
	baseURL string
	client  *http.Client
}

func NewAuthClient(cfg *config.AuthConfig) *AuthClient {
	return &AuthClient{
		baseURL: cfg.ServiceURL,
		client:  &http.Client{},
	}
}

func (c *AuthClient) GetUserRolesAndPermissions(userID uuid.UUID) ([]RolePermission, error) {
	url := fmt.Sprintf("%s/auth/users/%s/roles-with-permissions", c.baseURL, userID.String())
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth service returned status code: %d", resp.StatusCode)
	}

	var result struct {
		Success bool             `json:"success"`
		Data    []RolePermission `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *AuthClient) ExtractRolesAndPermissions(userID uuid.UUID) ([]string, []string, error) {
	rolePerms, err := c.GetUserRolesAndPermissions(userID)
	if err != nil {
		utils.Errorf("Failed to get roles and permissions", utils.ErrorField(err.Error()))
		// Return empty roles/permissions if auth service is unavailable
		return []string{}, []string{}, nil
	}

	roles := make([]string, 0)
	permissions := make([]string, 0)

	for _, rp := range rolePerms {
		roles = append(roles, rp.Role)
		permissions = append(permissions, rp.Permissions...)
	}

	return roles, permissions, nil
}
