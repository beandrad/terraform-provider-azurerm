package azuresdkhacks

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"

func CreateSecurityCenterContact(client *security.ContactsClient, ctx context.Context, securityContactName string, securityContact security.Contact) (result security.Contact, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ContactsClient.Create")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: securityContact,
			Constraints: []validation.Constraint{{Target: "securityContact.ContactProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "securityContact.ContactProperties.Email", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("security.ContactsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, securityContactName, securityContact)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.ContactsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "security.ContactsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = createResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.ContactsClient", "Create", resp, "Failure responding to request")
	}

	return
}

func createResponder(resp *http.Response) (result security.Contact, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
