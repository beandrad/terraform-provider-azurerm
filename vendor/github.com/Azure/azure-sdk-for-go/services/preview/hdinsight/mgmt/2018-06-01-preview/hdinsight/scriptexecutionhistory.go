package hdinsight

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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ScriptExecutionHistoryClient is the hDInsight Management Client
type ScriptExecutionHistoryClient struct {
	BaseClient
}

// NewScriptExecutionHistoryClient creates an instance of the ScriptExecutionHistoryClient client.
func NewScriptExecutionHistoryClient(subscriptionID string) ScriptExecutionHistoryClient {
	return NewScriptExecutionHistoryClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewScriptExecutionHistoryClientWithBaseURI creates an instance of the ScriptExecutionHistoryClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewScriptExecutionHistoryClientWithBaseURI(baseURI string, subscriptionID string) ScriptExecutionHistoryClient {
	return ScriptExecutionHistoryClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByCluster lists all scripts' execution history for the specified cluster.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the cluster.
func (client ScriptExecutionHistoryClient) ListByCluster(ctx context.Context, resourceGroupName string, clusterName string) (result ScriptActionExecutionHistoryListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ScriptExecutionHistoryClient.ListByCluster")
		defer func() {
			sc := -1
			if result.saehl.Response.Response != nil {
				sc = result.saehl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByClusterNextResults
	req, err := client.ListByClusterPreparer(ctx, resourceGroupName, clusterName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "ListByCluster", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByClusterSender(req)
	if err != nil {
		result.saehl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "ListByCluster", resp, "Failure sending request")
		return
	}

	result.saehl, err = client.ListByClusterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "ListByCluster", resp, "Failure responding to request")
	}
	if result.saehl.hasNextLink() && result.saehl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByClusterPreparer prepares the ListByCluster request.
func (client ScriptExecutionHistoryClient) ListByClusterPreparer(ctx context.Context, resourceGroupName string, clusterName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":       autorest.Encode("path", clusterName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HDInsight/clusters/{clusterName}/scriptExecutionHistory", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByClusterSender sends the ListByCluster request. The method will close the
// http.Response Body if it receives an error.
func (client ScriptExecutionHistoryClient) ListByClusterSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByClusterResponder handles the response to the ListByCluster request. The method always
// closes the http.Response Body.
func (client ScriptExecutionHistoryClient) ListByClusterResponder(resp *http.Response) (result ScriptActionExecutionHistoryList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByClusterNextResults retrieves the next set of results, if any.
func (client ScriptExecutionHistoryClient) listByClusterNextResults(ctx context.Context, lastResults ScriptActionExecutionHistoryList) (result ScriptActionExecutionHistoryList, err error) {
	req, err := lastResults.scriptActionExecutionHistoryListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "listByClusterNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByClusterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "listByClusterNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByClusterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "listByClusterNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByClusterComplete enumerates all values, automatically crossing page boundaries as required.
func (client ScriptExecutionHistoryClient) ListByClusterComplete(ctx context.Context, resourceGroupName string, clusterName string) (result ScriptActionExecutionHistoryListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ScriptExecutionHistoryClient.ListByCluster")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByCluster(ctx, resourceGroupName, clusterName)
	return
}

// Promote promotes the specified ad-hoc script execution to a persisted script.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the cluster.
// scriptExecutionID - the script execution Id
func (client ScriptExecutionHistoryClient) Promote(ctx context.Context, resourceGroupName string, clusterName string, scriptExecutionID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ScriptExecutionHistoryClient.Promote")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PromotePreparer(ctx, resourceGroupName, clusterName, scriptExecutionID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "Promote", nil, "Failure preparing request")
		return
	}

	resp, err := client.PromoteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "Promote", resp, "Failure sending request")
		return
	}

	result, err = client.PromoteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hdinsight.ScriptExecutionHistoryClient", "Promote", resp, "Failure responding to request")
	}

	return
}

// PromotePreparer prepares the Promote request.
func (client ScriptExecutionHistoryClient) PromotePreparer(ctx context.Context, resourceGroupName string, clusterName string, scriptExecutionID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":       autorest.Encode("path", clusterName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"scriptExecutionId": autorest.Encode("path", scriptExecutionID),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HDInsight/clusters/{clusterName}/scriptExecutionHistory/{scriptExecutionId}/promote", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PromoteSender sends the Promote request. The method will close the
// http.Response Body if it receives an error.
func (client ScriptExecutionHistoryClient) PromoteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// PromoteResponder handles the response to the Promote request. The method always
// closes the http.Response Body.
func (client ScriptExecutionHistoryClient) PromoteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
