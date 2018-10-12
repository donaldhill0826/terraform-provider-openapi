package openapi

import (
	"github.com/go-openapi/spec"
)

const extTfHeader = "x-terraform-header"

type parameterGroups [][]spec.Parameter

// getHeaderConfigurations gets all the header configurations for a specific
func getHeaderConfigurations(parameters []spec.Parameter) SpecHeaderParameters {
	return getHeaderConfigurationsForParameterGroups(parameterGroups{parameters})
}

// getHeaderConfigurationsForParameterGroups loops through the provided parametersGroup (collection of parameters per operation) and
// returns a map containing all the header configurations; the key will either be the value specified in the extTfHeader
// or if not present the default value will be the name of the header. In any case, the key name will be translated to
// a terraform compliant field name if needed (more details in convertToTerraformCompliantFieldName method)
func getHeaderConfigurationsForParameterGroups(parametersGroup parameterGroups) SpecHeaderParameters {
	headerParameters := SpecHeaderParameters{}
	for _, parameters := range parametersGroup {
		for _, parameter := range parameters {
			switch parameter.In {
			case "header":
				if preferredName, exists := parameter.Extensions.GetString(extTfHeader); exists {
					headerParameters = append(headerParameters, SpecHeaderParam{Name: parameter.Name, TerraformName: preferredName})
				} else {
					headerParameters = append(headerParameters, SpecHeaderParam{Name: parameter.Name})
				}
			}
		}
	}
	return headerParameters
}

// getPathHeaderParams aggregates all header type parameters found in the given path and returns the corresponding
// header configurations
func getPathHeaderParams(path spec.PathItem) SpecHeaderParameters {
	parametersGroup := parameterGroups{}
	parametersGroup = appendOperationParametersIfPresent(parametersGroup, path.Post)
	parametersGroup = appendOperationParametersIfPresent(parametersGroup, path.Get)
	parametersGroup = appendOperationParametersIfPresent(parametersGroup, path.Put)
	parametersGroup = appendOperationParametersIfPresent(parametersGroup, path.Delete)
	return getHeaderConfigurationsForParameterGroups(parametersGroup)
}

// appendOperationParametersIfPresent is a helper function that checks whether the given operation is not nil and if so
// appends its parameters to the parametersGroups
func appendOperationParametersIfPresent(parametersGroups parameterGroups, operation *spec.Operation) parameterGroups {
	if operation != nil {
		parametersGroups = append(parametersGroups, operation.Parameters)
	}
	return parametersGroups
}
