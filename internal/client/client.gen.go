// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"
)

// Device defines model for device.
type Device struct {
	Icon     *string     `json:"icon,omitempty"`
	Id       *string     `json:"id,omitempty"`
	Location *string     `json:"location,omitempty"`
	Make     *string     `json:"make,omitempty"`
	Model    *string     `json:"model,omitempty"`
	Name     *string     `json:"name,omitempty"`
	Tags     *DeviceTags `json:"tags,omitempty"`
}

// DeviceTags defines model for device_tags.
type DeviceTags map[string]interface{}

// Error defines model for error.
type Error struct {
	ErrorReason *string `json:"error_reason,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// Hello defines model for hello.
type Hello struct {
	AbCohort     *string       `json:"ab_cohort,omitempty"`
	AccessToken  *string       `json:"access_token,omitempty"`
	AccountId    *int          `json:"account_id,omitempty"`
	Authorized   *bool         `json:"authorized,omitempty"`
	BridgeServer *string       `json:"bridge_server,omitempty"`
	DateCreated  *time.Time    `json:"date_created,omitempty"`
	Monitors     *[]Monitor    `json:"monitors,omitempty"`
	RefreshToken *string       `json:"refresh_token,omitempty"`
	Settings     *UserSettings `json:"settings,omitempty"`
	TotpEnabled  *bool         `json:"totp_enabled,omitempty"`
	UserId       *int          `json:"user_id,omitempty"`
}

// Monitor defines model for monitor.
type Monitor struct {
	Attributes               *MonitorAttributes `json:"attributes,omitempty"`
	AuxIgnore                *bool              `json:"aux_ignore,omitempty"`
	AuxPort                  *string            `json:"aux_port,omitempty"`
	DataSharing              *[]interface{}     `json:"data_sharing,omitempty"`
	EthernetSupported        *bool              `json:"ethernet_supported,omitempty"`
	HardwareType             *string            `json:"hardware_type,omitempty"`
	Id                       *int               `json:"id,omitempty"`
	Online                   *bool              `json:"online,omitempty"`
	SerialNumber             *string            `json:"serial_number,omitempty"`
	SignalCheckCompletedTime *time.Time         `json:"signal_check_completed_time,omitempty"`
	SolarConfigured          *bool              `json:"solar_configured,omitempty"`
	SolarConnected           *bool              `json:"solar_connected,omitempty"`
	TimeZone                 *string            `json:"time_zone,omitempty"`
	ZigbeeSupported          *bool              `json:"zigbee_supported,omitempty"`
}

// MonitorAttributes defines model for monitor_attributes.
type MonitorAttributes struct {
	BasementType        *string      `json:"basement_type,omitempty"`
	BasementTypeKey     *interface{} `json:"basement_type_key,omitempty"`
	Cost                *float32     `json:"cost,omitempty"`
	CycleStart          *interface{} `json:"cycle_start,omitempty"`
	ElectricityCost     *interface{} `json:"electricity_cost,omitempty"`
	HomeSizeType        *string      `json:"home_size_type,omitempty"`
	HomeSizeTypeKey     *interface{} `json:"home_size_type_key,omitempty"`
	HomeType            *string      `json:"home_type,omitempty"`
	HomeTypeKey         *interface{} `json:"home_type_key,omitempty"`
	Id                  *int         `json:"id,omitempty"`
	Name                *interface{} `json:"name,omitempty"`
	NumberOfOccupants   *string      `json:"number_of_occupants,omitempty"`
	OccupancyType       *string      `json:"occupancy_type,omitempty"`
	OccupancyTypeKey    *interface{} `json:"occupancy_type_key,omitempty"`
	Panel               *interface{} `json:"panel,omitempty"`
	PostalCode          *string      `json:"postal_code,omitempty"`
	PowerRegion         *string      `json:"power_region,omitempty"`
	SellBackRate        *float32     `json:"sell_back_rate,omitempty"`
	ShowCost            *bool        `json:"show_cost,omitempty"`
	SolarTouEnabled     *bool        `json:"solar_tou_enabled,omitempty"`
	State               *interface{} `json:"state,omitempty"`
	ToGridThreshold     *interface{} `json:"to_grid_threshold,omitempty"`
	TouEnabled          *bool        `json:"tou_enabled,omitempty"`
	UserSetCost         *bool        `json:"user_set_cost,omitempty"`
	UserSetSellBackRate *bool        `json:"user_set_sell_back_rate,omitempty"`
	YearBuiltType       *string      `json:"year_built_type,omitempty"`
	YearBuiltTypeKey    *interface{} `json:"year_built_type_key,omitempty"`
}

// NotificationSettings defines model for notification_settings.
type NotificationSettings struct {
	AlwaysOnChangePush       *bool `json:"always_on_change_push,omitempty"`
	ComparisonChangePush     *bool `json:"comparison_change_push,omitempty"`
	DailyChangePush          *bool `json:"daily_change_push,omitempty"`
	GeneratorOffPush         *bool `json:"generator_off_push,omitempty"`
	GeneratorOnPush          *bool `json:"generator_on_push,omitempty"`
	GridOutagePush           *bool `json:"grid_outage_push,omitempty"`
	GridRestoredPush         *bool `json:"grid_restored_push,omitempty"`
	MonitorMonthlyEmail      *bool `json:"monitor_monthly_email,omitempty"`
	MonitorOfflineEmail      *bool `json:"monitor_offline_email,omitempty"`
	MonitorOfflinePush       *bool `json:"monitor_offline_push,omitempty"`
	MonthlyChangePush        *bool `json:"monthly_change_push,omitempty"`
	NewNamedDeviceEmail      *bool `json:"new_named_device_email,omitempty"`
	NewNamedDevicePush       *bool `json:"new_named_device_push,omitempty"`
	NewPeakEmail             *bool `json:"new_peak_email,omitempty"`
	NewPeakPush              *bool `json:"new_peak_push,omitempty"`
	RelayUpdateAvailablePush *bool `json:"relay_update_available_push,omitempty"`
	RelayUpdateInstalledPush *bool `json:"relay_update_installed_push,omitempty"`
	TimeOfUse                *bool `json:"time_of_use,omitempty"`
	WeeklyChangePush         *bool `json:"weekly_change_push,omitempty"`
}

// UserSettings defines model for user_settings.
type UserSettings struct {
	Settings *struct {
		LabsEnabled   *bool                            `json:"labs_enabled,omitempty"`
		Notifications *map[string]NotificationSettings `json:"notifications,omitempty"`
	} `json:"settings,omitempty"`
	UserId  *int `json:"user_id,omitempty"`
	Version *int `json:"version,omitempty"`
}

// GetDevicesParams defines parameters for GetDevices.
type GetDevicesParams struct {
	IncludeMerged *bool `form:"include_merged,omitempty" json:"include_merged,omitempty"`
}

// AuthenticateFormdataBody defines parameters for Authenticate.
type AuthenticateFormdataBody struct {
	Email    *string `form:"email,omitempty" json:"email,omitempty"`
	MfaToken *string `form:"mfa_token,omitempty" json:"mfa_token,omitempty"`
	Password *string `form:"password,omitempty" json:"password,omitempty"`
	Totp     *string `form:"totp,omitempty" json:"totp,omitempty"`
}

// RenewAuthTokenFormdataBody defines parameters for RenewAuthToken.
type RenewAuthTokenFormdataBody struct {
	IsAccessToken *bool   `form:"is_access_token,omitempty" json:"is_access_token,omitempty"`
	RefreshToken  *string `form:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	UserId        *int    `form:"user_id,omitempty" json:"user_id,omitempty"`
}

// AuthenticateFormdataRequestBody defines body for Authenticate for application/x-www-form-urlencoded ContentType.
type AuthenticateFormdataRequestBody AuthenticateFormdataBody

// RenewAuthTokenFormdataRequestBody defines body for RenewAuthToken for application/x-www-form-urlencoded ContentType.
type RenewAuthTokenFormdataRequestBody RenewAuthTokenFormdataBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetDevices request
	GetDevices(ctx context.Context, monitorId int, params *GetDevicesParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AuthenticateWithBody request with any body
	AuthenticateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AuthenticateWithFormdataBody(ctx context.Context, body AuthenticateFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetEnvironments request
	GetEnvironments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RenewAuthTokenWithBody request with any body
	RenewAuthTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RenewAuthTokenWithFormdataBody(ctx context.Context, body RenewAuthTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetDevices(ctx context.Context, monitorId int, params *GetDevicesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDevicesRequest(c.Server, monitorId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AuthenticateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthenticateRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AuthenticateWithFormdataBody(ctx context.Context, body AuthenticateFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthenticateRequestWithFormdataBody(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetEnvironments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetEnvironmentsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RenewAuthTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRenewAuthTokenRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RenewAuthTokenWithFormdataBody(ctx context.Context, body RenewAuthTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRenewAuthTokenRequestWithFormdataBody(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetDevicesRequest generates requests for GetDevices
func NewGetDevicesRequest(server string, monitorId int, params *GetDevicesParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "monitor_id", runtime.ParamLocationPath, monitorId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/monitors/%s/devices/overview", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.IncludeMerged != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "include_merged", runtime.ParamLocationQuery, *params.IncludeMerged); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAuthenticateRequestWithFormdataBody calls the generic Authenticate builder with application/x-www-form-urlencoded body
func NewAuthenticateRequestWithFormdataBody(server string, body AuthenticateFormdataRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	bodyStr, err := runtime.MarshalForm(body, nil)
	if err != nil {
		return nil, err
	}
	bodyReader = strings.NewReader(bodyStr.Encode())
	return NewAuthenticateRequestWithBody(server, "application/x-www-form-urlencoded", bodyReader)
}

// NewAuthenticateRequestWithBody generates requests for Authenticate with any type of body
func NewAuthenticateRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/authenticate")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetEnvironmentsRequest generates requests for GetEnvironments
func NewGetEnvironmentsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/public/monitors/environments")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRenewAuthTokenRequestWithFormdataBody calls the generic RenewAuthToken builder with application/x-www-form-urlencoded body
func NewRenewAuthTokenRequestWithFormdataBody(server string, body RenewAuthTokenFormdataRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	bodyStr, err := runtime.MarshalForm(body, nil)
	if err != nil {
		return nil, err
	}
	bodyReader = strings.NewReader(bodyStr.Encode())
	return NewRenewAuthTokenRequestWithBody(server, "application/x-www-form-urlencoded", bodyReader)
}

// NewRenewAuthTokenRequestWithBody generates requests for RenewAuthToken with any type of body
func NewRenewAuthTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/renew")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetDevicesWithResponse request
	GetDevicesWithResponse(ctx context.Context, monitorId int, params *GetDevicesParams, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error)

	// AuthenticateWithBodyWithResponse request with any body
	AuthenticateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error)

	AuthenticateWithFormdataBodyWithResponse(ctx context.Context, body AuthenticateFormdataRequestBody, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error)

	// GetEnvironmentsWithResponse request
	GetEnvironmentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEnvironmentsResponse, error)

	// RenewAuthTokenWithBodyWithResponse request with any body
	RenewAuthTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RenewAuthTokenResponse, error)

	RenewAuthTokenWithFormdataBodyWithResponse(ctx context.Context, body RenewAuthTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*RenewAuthTokenResponse, error)
}

type GetDevicesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		DeviceDataChecksum *string   `json:"device_data_checksum,omitempty"`
		Devices            *[]Device `json:"devices,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r GetDevicesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDevicesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AuthenticateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Hello
	JSON401      *struct {
		ErrorReason *string `json:"error_reason,omitempty"`
		MfaToken    *string `json:"mfa_token,omitempty"`
		Status      *string `json:"status,omitempty"`
	}
	JSONDefault *Error
}

// Status returns HTTPResponse.Status
func (r AuthenticateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AuthenticateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetEnvironmentsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]struct {
		ApiUrl              *string `json:"api_url,omitempty"`
		ClientBridgelinkUrl *string `json:"client_bridgelink_url,omitempty"`
		DisplayName         *string `json:"display_name,omitempty"`
		Environment         *string `json:"environment,omitempty"`
	}
	JSONDefault *Error
}

// Status returns HTTPResponse.Status
func (r GetEnvironmentsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetEnvironmentsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RenewAuthTokenResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		AccessToken  *string    `json:"access_token,omitempty"`
		Expires      *time.Time `json:"expires,omitempty"`
		RefreshToken *string    `json:"refresh_token,omitempty"`
	}
	JSONDefault *Error
}

// Status returns HTTPResponse.Status
func (r RenewAuthTokenResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RenewAuthTokenResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetDevicesWithResponse request returning *GetDevicesResponse
func (c *ClientWithResponses) GetDevicesWithResponse(ctx context.Context, monitorId int, params *GetDevicesParams, reqEditors ...RequestEditorFn) (*GetDevicesResponse, error) {
	rsp, err := c.GetDevices(ctx, monitorId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDevicesResponse(rsp)
}

// AuthenticateWithBodyWithResponse request with arbitrary body returning *AuthenticateResponse
func (c *ClientWithResponses) AuthenticateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error) {
	rsp, err := c.AuthenticateWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthenticateResponse(rsp)
}

func (c *ClientWithResponses) AuthenticateWithFormdataBodyWithResponse(ctx context.Context, body AuthenticateFormdataRequestBody, reqEditors ...RequestEditorFn) (*AuthenticateResponse, error) {
	rsp, err := c.AuthenticateWithFormdataBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthenticateResponse(rsp)
}

// GetEnvironmentsWithResponse request returning *GetEnvironmentsResponse
func (c *ClientWithResponses) GetEnvironmentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEnvironmentsResponse, error) {
	rsp, err := c.GetEnvironments(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetEnvironmentsResponse(rsp)
}

// RenewAuthTokenWithBodyWithResponse request with arbitrary body returning *RenewAuthTokenResponse
func (c *ClientWithResponses) RenewAuthTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RenewAuthTokenResponse, error) {
	rsp, err := c.RenewAuthTokenWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRenewAuthTokenResponse(rsp)
}

func (c *ClientWithResponses) RenewAuthTokenWithFormdataBodyWithResponse(ctx context.Context, body RenewAuthTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*RenewAuthTokenResponse, error) {
	rsp, err := c.RenewAuthTokenWithFormdataBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRenewAuthTokenResponse(rsp)
}

// ParseGetDevicesResponse parses an HTTP response from a GetDevicesWithResponse call
func ParseGetDevicesResponse(rsp *http.Response) (*GetDevicesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDevicesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			DeviceDataChecksum *string   `json:"device_data_checksum,omitempty"`
			Devices            *[]Device `json:"devices,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAuthenticateResponse parses an HTTP response from a AuthenticateWithResponse call
func ParseAuthenticateResponse(rsp *http.Response) (*AuthenticateResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AuthenticateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Hello
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest struct {
			ErrorReason *string `json:"error_reason,omitempty"`
			MfaToken    *string `json:"mfa_token,omitempty"`
			Status      *string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetEnvironmentsResponse parses an HTTP response from a GetEnvironmentsWithResponse call
func ParseGetEnvironmentsResponse(rsp *http.Response) (*GetEnvironmentsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetEnvironmentsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []struct {
			ApiUrl              *string `json:"api_url,omitempty"`
			ClientBridgelinkUrl *string `json:"client_bridgelink_url,omitempty"`
			DisplayName         *string `json:"display_name,omitempty"`
			Environment         *string `json:"environment,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseRenewAuthTokenResponse parses an HTTP response from a RenewAuthTokenWithResponse call
func ParseRenewAuthTokenResponse(rsp *http.Response) (*RenewAuthTokenResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RenewAuthTokenResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			AccessToken  *string    `json:"access_token,omitempty"`
			Expires      *time.Time `json:"expires,omitempty"`
			RefreshToken *string    `json:"refresh_token,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
