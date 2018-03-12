// Package yggdrasil provides methods to utilize Mojang's Yggdrasil API.
package yggdrasil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client holds an access token and a client token.
// After a successful authentication, it will also hold the currently selected profile and the current user.
type Client struct {
	AccessToken     string
	ClientToken     string
	SelectedProfile Profile
	User            User
}

// Error holds data about a Yggdrasil or internal error.
type Error struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Cause        string `json:"cause"`
	StatusCode   int

	FuncError error
}

// AuthenticationRequest holds data used to make an authentication request.
type AuthenticationRequest struct {
	Agent       Agent  `json:"agent"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToken"`
	RequestUser bool   `json:"requestUser"`
}

// AuthenticationResponse holds data returned from a successful authentication request.
type AuthenticationResponse struct {
	AccessToken       string    `json:"accessToken"`
	ClientToken       string    `json:"clientToken"`
	AvailableProfiles []Profile `json:"availableProfiles"`
	SelectedProfile   Profile   `json:"selectedProfile"`
	User              User      `json:"user"`
}

// RefreshRequest holds data used to make a refresh request.
type RefreshRequest struct {
	AccessToken     string  `json:"accessToken"`
	ClientToken     string  `json:"clientToken"`
	SelectedProfile Profile `json:"selectedProfile"`
	RequestUser     bool    `json:"requestUser"`
}

// RefreshResponse holds data returned from a successful refresh request.
type RefreshResponse struct {
	AccessToken     string  `json:"accessToken"`
	ClientToken     string  `json:"clientToken"`
	SelectedProfile Profile `json:"selectedProfile"`
	User            User    `json:"user"`
}

// ValidateRequest holds data used to make a validate request.
type ValidateRequest struct {
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken"`
}

// SignoutRequest holds data used to make a signout request.
type SignoutRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// InvalidateRequest holds data used to make an invalidate request.
type InvalidateRequest struct {
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken"`
}

// Agent holds data about the game that was authenticated for.
type Agent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

// Profile holds data about an authenticated user's profile.
type Profile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Legacy bool   `json:"legacy"`
}

// User holds data about an authenticated user.
type User struct {
	ID         string     `json:"id"`
	Properties []Property `json:"properties"`
}

// Property holds data about an authenticated user's property.
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Authenticate attempts to authenticate with Yggdrasil.
func (client *Client) Authenticate(username, password, gameName string, gameVersion int) (*AuthenticationResponse, *Error) {
	authRequest := &AuthenticationRequest{
		Agent: Agent{
			Name:    gameName,
			Version: gameVersion},
		Username:    username,
		Password:    password,
		ClientToken: client.ClientToken,
		RequestUser: true}

	requestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	requestJSONBuffer := bytes.NewBuffer([]byte(requestJSON))

	request, err := http.NewRequest("POST", "https://authserver.mojang.com/authenticate", requestJSONBuffer)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &Error{FuncError: err}
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 200 {
		var errorResponse *Error
		err = json.Unmarshal(body, &errorResponse)

		if err != nil {
			return nil, &Error{FuncError: err}
		}

		errorResponse.StatusCode = response.StatusCode
		return nil, errorResponse
	}

	var authResponse *AuthenticationResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	client.AccessToken = authResponse.AccessToken
	client.SelectedProfile = authResponse.SelectedProfile
	client.User = authResponse.User

	return authResponse, nil
}

// Refresh attempts to refresh an existing access/client token pair to get a new valid access token.
func (client *Client) Refresh() (*RefreshResponse, *Error) {
	refreshRequest := &RefreshRequest{
		AccessToken: client.AccessToken,
		ClientToken: client.ClientToken,
		RequestUser: true}

	requestJSON, err := json.Marshal(refreshRequest)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	requestJSONBuffer := bytes.NewBuffer([]byte(requestJSON))

	request, err := http.NewRequest("POST", "https://authserver.mojang.com/refresh", requestJSONBuffer)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &Error{FuncError: err}
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 200 {
		var errorResponse *Error
		err = json.Unmarshal(body, &errorResponse)

		if err != nil {
			return nil, &Error{FuncError: err}
		}

		errorResponse.StatusCode = response.StatusCode
		return nil, errorResponse
	}

	var refreshResponse *RefreshResponse
	err = json.Unmarshal(body, &refreshResponse)
	if err != nil {
		return nil, &Error{FuncError: err}
	}

	client.AccessToken = refreshResponse.AccessToken
	client.SelectedProfile = refreshResponse.SelectedProfile
	client.User = refreshResponse.User

	return refreshResponse, nil
}

// Validate attempts to check whether or not an existing access/client token pair is valid.
func (client *Client) Validate() (bool, *Error) {
	validateRequest := &ValidateRequest{
		AccessToken: client.AccessToken,
		ClientToken: client.ClientToken}

	requestJSON, err := json.Marshal(validateRequest)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	requestJSONBuffer := bytes.NewBuffer([]byte(requestJSON))

	request, err := http.NewRequest("POST", "https://authserver.mojang.com/validate", requestJSONBuffer)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, &Error{FuncError: err}
	}
	defer response.Body.Close()

	if response.StatusCode == 204 {
		return true, nil
	} else if response.StatusCode == 403 {
		var errorResponse *Error
		err = json.Unmarshal(body, &errorResponse)

		if err != nil {
			return false, &Error{FuncError: err}
		}

		errorResponse.StatusCode = response.StatusCode
		return false, errorResponse
	}

	return false, nil
}

// Signout attempts to signout of a legacy Minecraft account.
func (client *Client) Signout(username, password string) (bool, *Error) {
	signoutRequest := &SignoutRequest{
		Username: username,
		Password: password}

	requestJSON, err := json.Marshal(signoutRequest)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	requestJSONBuffer := bytes.NewBuffer([]byte(requestJSON))

	request, err := http.NewRequest("POST", "https://authserver.mojang.com/signout", requestJSONBuffer)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, &Error{FuncError: err}
	}
	defer response.Body.Close()

	if len(body) == 0 {
		return true, nil
	}
	
	var errorResponse *Error
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return false, &Error{FuncError: err}
	}

	errorResponse.StatusCode = response.StatusCode
	return false, errorResponse
}

// Invalidate attempts to invalidate an existing access/client token pair.
func (client *Client) Invalidate() *Error {
	invalidateRequest := &InvalidateRequest{
		AccessToken: client.AccessToken,
		ClientToken: client.ClientToken}

	requestJSON, err := json.Marshal(invalidateRequest)
	if err != nil {
		return &Error{FuncError: err}
	}

	requestJSONBuffer := bytes.NewBuffer([]byte(requestJSON))

	request, err := http.NewRequest("POST", "https://authserver.mojang.com/invalidate", requestJSONBuffer)
	if err != nil {
		return &Error{FuncError: err}
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return &Error{FuncError: err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Error{FuncError: err}
	}
	defer response.Body.Close()

	if len(body) == 0 {
		return nil
	}

	var errorResponse *Error
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return &Error{FuncError: err}
	}

	errorResponse.StatusCode = response.StatusCode
	return errorResponse
}