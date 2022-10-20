package account

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// DataLakeStoreAccountsClient is the creates an Azure Data Lake Analytics account management client.
type DataLakeStoreAccountsClient struct {
	BaseClient
}

// NewDataLakeStoreAccountsClient creates an instance of the DataLakeStoreAccountsClient client.
func NewDataLakeStoreAccountsClient(subscriptionID string) DataLakeStoreAccountsClient {
	return NewDataLakeStoreAccountsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDataLakeStoreAccountsClientWithBaseURI creates an instance of the DataLakeStoreAccountsClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewDataLakeStoreAccountsClientWithBaseURI(baseURI string, subscriptionID string) DataLakeStoreAccountsClient {
	return DataLakeStoreAccountsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Add updates the specified Data Lake Analytics account to include the additional Data Lake Store account.
// Parameters:
// resourceGroupName - the name of the Azure resource group.
// accountName - the name of the Data Lake Analytics account.
// dataLakeStoreAccountName - the name of the Data Lake Store account to add.
// parameters - the details of the Data Lake Store account.
func (client DataLakeStoreAccountsClient) Add(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string, parameters *AddDataLakeStoreParameters) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DataLakeStoreAccountsClient.Add")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AddPreparer(ctx, resourceGroupName, accountName, dataLakeStoreAccountName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Add", nil, "Failure preparing request")
		return
	}

	resp, err := client.AddSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Add", resp, "Failure sending request")
		return
	}

	result, err = client.AddResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Add", resp, "Failure responding to request")
	}

	return
}

// AddPreparer prepares the Add request.
func (client DataLakeStoreAccountsClient) AddPreparer(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string, parameters *AddDataLakeStoreParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":              autorest.Encode("path", accountName),
		"dataLakeStoreAccountName": autorest.Encode("path", dataLakeStoreAccountName),
		"resourceGroupName":        autorest.Encode("path", resourceGroupName),
		"subscriptionId":           autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/dataLakeStoreAccounts/{dataLakeStoreAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if parameters != nil {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithJSON(parameters))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AddSender sends the Add request. The method will close the
// http.Response Body if it receives an error.
func (client DataLakeStoreAccountsClient) AddSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// AddResponder handles the response to the Add request. The method always
// closes the http.Response Body.
func (client DataLakeStoreAccountsClient) AddResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete updates the Data Lake Analytics account specified to remove the specified Data Lake Store account.
// Parameters:
// resourceGroupName - the name of the Azure resource group.
// accountName - the name of the Data Lake Analytics account.
// dataLakeStoreAccountName - the name of the Data Lake Store account to remove
func (client DataLakeStoreAccountsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DataLakeStoreAccountsClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, accountName, dataLakeStoreAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client DataLakeStoreAccountsClient) DeletePreparer(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":              autorest.Encode("path", accountName),
		"dataLakeStoreAccountName": autorest.Encode("path", dataLakeStoreAccountName),
		"resourceGroupName":        autorest.Encode("path", resourceGroupName),
		"subscriptionId":           autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/dataLakeStoreAccounts/{dataLakeStoreAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client DataLakeStoreAccountsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client DataLakeStoreAccountsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the specified Data Lake Store account details in the specified Data Lake Analytics account.
// Parameters:
// resourceGroupName - the name of the Azure resource group.
// accountName - the name of the Data Lake Analytics account.
// dataLakeStoreAccountName - the name of the Data Lake Store account to retrieve
func (client DataLakeStoreAccountsClient) Get(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string) (result DataLakeStoreAccountInformation, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DataLakeStoreAccountsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, accountName, dataLakeStoreAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DataLakeStoreAccountsClient) GetPreparer(ctx context.Context, resourceGroupName string, accountName string, dataLakeStoreAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":              autorest.Encode("path", accountName),
		"dataLakeStoreAccountName": autorest.Encode("path", dataLakeStoreAccountName),
		"resourceGroupName":        autorest.Encode("path", resourceGroupName),
		"subscriptionId":           autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/dataLakeStoreAccounts/{dataLakeStoreAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DataLakeStoreAccountsClient) GetSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DataLakeStoreAccountsClient) GetResponder(resp *http.Response) (result DataLakeStoreAccountInformation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAccount gets the first page of Data Lake Store accounts linked to the specified Data Lake Analytics account.
// The response includes a link to the next page, if any.
// Parameters:
// resourceGroupName - the name of the Azure resource group.
// accountName - the name of the Data Lake Analytics account.
// filter - oData filter. Optional.
// top - the number of items to return. Optional.
// skip - the number of items to skip over before returning elements. Optional.
// selectParameter - oData Select statement. Limits the properties on each entry to just those requested, e.g.
// Categories?$select=CategoryName,Description. Optional.
// orderby - orderBy clause. One or more comma-separated expressions with an optional "asc" (the default) or
// "desc" depending on the order you'd like the values sorted, e.g. Categories?$orderby=CategoryName desc.
// Optional.
// count - the Boolean value of true or false to request a count of the matching resources included with the
// resources in the response, e.g. Categories?$count=true. Optional.
func (client DataLakeStoreAccountsClient) ListByAccount(ctx context.Context, resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool) (result DataLakeStoreAccountInformationListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DataLakeStoreAccountsClient.ListByAccount")
		defer func() {
			sc := -1
			if result.dlsailr.Response.Response != nil {
				sc = result.dlsailr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil}}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("account.DataLakeStoreAccountsClient", "ListByAccount", err.Error())
	}

	result.fn = client.listByAccountNextResults
	req, err := client.ListByAccountPreparer(ctx, resourceGroupName, accountName, filter, top, skip, selectParameter, orderby, count)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "ListByAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.dlsailr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "ListByAccount", resp, "Failure sending request")
		return
	}

	result.dlsailr, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "ListByAccount", resp, "Failure responding to request")
	}

	return
}

// ListByAccountPreparer prepares the ListByAccount request.
func (client DataLakeStoreAccountsClient) ListByAccountPreparer(ctx context.Context, resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if len(orderby) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderby)
	}
	if count != nil {
		queryParameters["$count"] = autorest.Encode("query", *count)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/dataLakeStoreAccounts", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByAccountSender sends the ListByAccount request. The method will close the
// http.Response Body if it receives an error.
func (client DataLakeStoreAccountsClient) ListByAccountSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListByAccountResponder handles the response to the ListByAccount request. The method always
// closes the http.Response Body.
func (client DataLakeStoreAccountsClient) ListByAccountResponder(resp *http.Response) (result DataLakeStoreAccountInformationListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByAccountNextResults retrieves the next set of results, if any.
func (client DataLakeStoreAccountsClient) listByAccountNextResults(ctx context.Context, lastResults DataLakeStoreAccountInformationListResult) (result DataLakeStoreAccountInformationListResult, err error) {
	req, err := lastResults.dataLakeStoreAccountInformationListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "listByAccountNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "listByAccountNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.DataLakeStoreAccountsClient", "listByAccountNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByAccountComplete enumerates all values, automatically crossing page boundaries as required.
func (client DataLakeStoreAccountsClient) ListByAccountComplete(ctx context.Context, resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool) (result DataLakeStoreAccountInformationListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DataLakeStoreAccountsClient.ListByAccount")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByAccount(ctx, resourceGroupName, accountName, filter, top, skip, selectParameter, orderby, count)
	return
}
