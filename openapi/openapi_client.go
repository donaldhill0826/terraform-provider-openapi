package openapi

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/dikhan/terraform-provider-openapi/openapi/version"

	"github.com/dikhan/http_goclient"
)

type httpMethodSupported string

const (
	httpGet    httpMethodSupported = "GET"
	httpPost   httpMethodSupported = "POST"
	httpPut    httpMethodSupported = "PUT"
	httpDelete httpMethodSupported = "DELETE"
)

// ClientOpenAPI defines the behaviour expected to be implemented for the OpenAPI Client used in the Terraform OpenAPI Provider
type ClientOpenAPI interface {
	Post(resource SpecResource, requestPayload interface{}, responsePayload interface{}) (*http.Response, error)
	Put(resource SpecResource, id string, requestPayload interface{}, responsePayload interface{}) (*http.Response, error)
	Get(resource SpecResource, id string, responsePayload interface{}) (*http.Response, error)
	Delete(resource SpecResource, id string) (*http.Response, error)
}

// ProviderClient defines a client that is configured based on the OpenAPI server side documentation
// The CRUD operations accept an OpenAPI operation which defines among other things the security scheme applicable to
// the API when making the HTTP request
type ProviderClient struct {
	openAPIBackendConfiguration SpecBackendConfiguration
	httpClient                  http_goclient.HttpClientIface
	providerConfiguration       providerConfiguration
	apiAuthenticator            specAuthenticator
}

// Post performs a POST request to the server API based on the resource configuration and the payload passed in
// TODO: This function will have to accept an array of IDs now. For instance a URI like this /v1/cdns/1234/firewalls will be represented in the array like []string{"1234"} where 1234 will be the cdns parent ID.
func (o *ProviderClient) Post(resource SpecResource, requestPayload interface{}, responsePayload interface{}) (*http.Response, error) {
	// TODO: pass in the ids args from post method
	ids := []string{}
	resourceURL, err := o.getResourceURL(resource, ids)
	if err != nil {
		return nil, err
	}
	operation := resource.getResourceOperations().Post
	return o.performRequest(httpPost, resourceURL, operation, requestPayload, responsePayload)
}

// Put performs a PUT request to the server API based on the resource configuration and the payload passed in
// TODO: Replace param from being just a string to being an array of strings. For instance a URI like this /v1/cdns/1234/firewalls will be represented in the array like []string{"1234", "567"} where 1234 will be the cdns parent ID AND 567 will be the firewall ID
func (o *ProviderClient) Put(resource SpecResource, id string, requestPayload interface{}, responsePayload interface{}) (*http.Response, error) {
	// TODO: pass in the ids args from Put method containing both the parent ids as well as the instance id
	ids := []string{}
	resourceURL, err := o.getResourceIDURL(resource, ids)
	if err != nil {
		return nil, err
	}
	operation := resource.getResourceOperations().Put
	return o.performRequest(httpPut, resourceURL, operation, requestPayload, responsePayload)
}

// Get performs a GET request to the server API based on the resource configuration and the resource instance id passed in
// TODO: Replace param from being just a string to being an array of strings. For instance a URI like this /v1/cdns/1234/firewalls will be represented in the array like []string{"1234", "567"} where 1234 will be the cdns parent ID AND 567 will be the firewall ID
func (o *ProviderClient) Get(resource SpecResource, id string, responsePayload interface{}) (*http.Response, error) {
	// TODO: pass in the ids args from Put method containing both the parent ids as well as the instance id
	ids := []string{}
	resourceURL, err := o.getResourceIDURL(resource, ids)
	if err != nil {
		return nil, err
	}
	operation := resource.getResourceOperations().Get
	return o.performRequest(httpGet, resourceURL, operation, nil, responsePayload)
}

// Delete performs a DELETE request to the server API based on the resource configuration and the resource instance id passed in
// TODO: Replace param from being just a string to being an array of strings. For instance a URI like this /v1/cdns/1234/firewalls will be represented in the array like []string{"1234", "567"} where 1234 will be the cdns parent ID AND 567 will be the firewall ID
func (o *ProviderClient) Delete(resource SpecResource, id string) (*http.Response, error) {
	// TODO: pass in the ids args from Put method containing both the parent ids as well as the instance id
	ids := []string{}
	resourceURL, err := o.getResourceIDURL(resource, ids)
	if err != nil {
		return nil, err
	}
	operation := resource.getResourceOperations().Delete
	return o.performRequest(httpDelete, resourceURL, operation, nil, nil)
}

func (o *ProviderClient) performRequest(method httpMethodSupported, resourceURL string, operation *specResourceOperation, requestPayload interface{}, responsePayload interface{}) (*http.Response, error) {
	reqContext, err := o.apiAuthenticator.prepareAuth(resourceURL, operation.SecuritySchemes, o.providerConfiguration)
	if err != nil {
		return nil, err
	}
	o.appendOperationHeaders(operation.HeaderParameters, o.providerConfiguration, reqContext.headers)
	log.Printf("[DEBUG] Performing %s %s", method, reqContext.url)

	userAgentHeader := version.BuildUserAgent(runtime.GOOS, runtime.GOARCH)
	o.appendUserAgentHeader(reqContext.headers, userAgentHeader)

	o.logHeadersSafely(reqContext.headers)

	switch method {
	case httpPost:
		return o.httpClient.PostJson(reqContext.url, reqContext.headers, requestPayload, &responsePayload)
	case httpPut:
		return o.httpClient.PutJson(reqContext.url, reqContext.headers, requestPayload, &responsePayload)
	case httpGet:
		return o.httpClient.Get(reqContext.url, reqContext.headers, &responsePayload)
	case httpDelete:
		return o.httpClient.Delete(reqContext.url, reqContext.headers)
	}
	return nil, fmt.Errorf("method '%s' not supported", method)
}

func (o *ProviderClient) appendUserAgentHeader(headers map[string]string, value string) {
	headers[userAgent] = value
}

// logHeadersSafely logs the header names sent to the APIs but the values are redacted for security reasons in case
// values contain secrets. However, the logging will display whether the values contained data or not so it's easier
// to debug whether the headers sent had data.
func (o *ProviderClient) logHeadersSafely(headers map[string]string) {
	for headerName, headerValue := range headers {
		if headerValue == "" {
			log.Printf("[DEBUG] Request Header '%s' sent with empty value :(", headerName)
		}
		log.Printf("[DEBUG] Request Header '%s' sent", headerName)
	}
}

// appendOperationHeaders returns a maps containing the headers passed in and adds whatever headers the operation requires. The values
// are retrieved from the provider configuration.
func (o ProviderClient) appendOperationHeaders(operationHeaders []SpecHeaderParam, providerConfig providerConfiguration, headers map[string]string) {
	if operationHeaders != nil && len(operationHeaders) > 0 {
		for _, headerParam := range operationHeaders {
			// Setting the actual name of the header with the expectedValue coming from the provider configuration
			headers[headerParam.Name] = providerConfig.getHeaderValueFor(headerParam)
		}
	}
}

func (o ProviderClient) getResourceURL(resource SpecResource, ids []string) (string, error) {
	var host string
	var err error

	isMultiRegion, _, regions, err := o.openAPIBackendConfiguration.isMultiRegion()
	if err != nil {
		return "", err
	}
	if isMultiRegion {
		// get region value provided by user in the terraform configuration file
		region := o.providerConfiguration.getRegion()
		// otherwise, if not provided falling back to the default value specified in the service provider swagger file
		if region == "" {
			region, err = o.openAPIBackendConfiguration.getDefaultRegion(regions)
			if err != nil {
				return "", err
			}
		}
		host, err = o.openAPIBackendConfiguration.getHostByRegion(region)
		if err != nil {
			return "", err
		}
	} else {
		host, err = o.openAPIBackendConfiguration.getHost()
		if err != nil {
			return "", err
		}
	}

	basePath := o.openAPIBackendConfiguration.getBasePath()
	resourceRelativePath, err := resource.getResourcePath(ids)
	if err != nil {
		return "", err
	}

	// Fall back to override the host if value is not empty; otherwise global host will be used as usual
	hostOverride, err := resource.getHost()
	if err != nil {
		return "", err
	}
	if hostOverride != "" {
		log.Printf("[INFO] resource '%s' is configured with host override, API calls will be made against '%s' instead of '%s'", resourceRelativePath, hostOverride, host)
		host = hostOverride
	}

	if endPointHost := o.providerConfiguration.getEndPoint(resource.getResourceName()); endPointHost != "" {
		log.Printf("[INFO] resource '%s' is configured with endpoint override, API calls will be made against '%s' instead of '%s'", resourceRelativePath, endPointHost, host)
		host = endPointHost
	}

	if host == "" || resourceRelativePath == "" {
		return "", fmt.Errorf("host and path are mandatory attributes to get the resource URL - host['%s'], path['%s']", host, resourceRelativePath)
	}

	// TODO: use resource operation schemes if specified
	defaultScheme := "http"
	for _, scheme := range o.openAPIBackendConfiguration.getHTTPSchemes() {
		if scheme == "https" {
			defaultScheme = "https"
		}
	}
	path := resourceRelativePath
	if strings.Index(resourceRelativePath, "/") != 0 {
		path = fmt.Sprintf("/%s", resourceRelativePath)
	}

	if basePath != "" && basePath != "/" {
		if strings.Index(basePath, "/") == 0 {
			return fmt.Sprintf("%s://%s%s%s", defaultScheme, host, basePath, path), nil
		}
		return fmt.Sprintf("%s://%s/%s%s", defaultScheme, host, basePath, path), nil
	}
	return fmt.Sprintf("%s://%s%s", defaultScheme, host, path), nil
}

// TODO: Expand (o ProviderClient) getResourceIDURL(resource SpecResource, id string) (string, error) function to handle
// TODO: subresourc URLs including parent ids. The expectation here is that given a spec resource that is subresource,
// TODO: the URL returned should have the path's path params already resolved with the right parent IDs. In order to achieve this,
// TODO: the function will need to accept a new parameter being an array of IDs (array so we support not just one level subresource but multiple).
func (o ProviderClient) getResourceIDURL(resource SpecResource, ids []string) (string, error) {
	if len(ids) <= 0 {
		return "", fmt.Errorf("getResourceIDURL cannot be called without the instance id")
	}

	url, err := o.getResourceURL(resource, ids)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(url, "/") {
		return fmt.Sprintf("%s%s", url, ids[len(ids)-1]), nil
	}
	return fmt.Sprintf("%s/%s", url, ids[len(ids)-1]), nil
}
