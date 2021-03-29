package tests

import (
	"bytes"
	"log"
	"net/http"

	"github.com/cloud-bulldozer/ocm-api-load/pkg/helpers"
	"github.com/cloud-bulldozer/ocm-api-load/pkg/tests/handlers"
	authv1 "github.com/openshift-online/ocm-sdk-go/authorizations/v1"
)

// Specify Test Cases
// They are written this way to re-use functionality where possible and
// hopefully make it easier to modify and/or extend given the declarative
// style.
var tests = []helpers.TestOptions{
	{
		TestName: "self-access-token",
		Path:     "/api/accounts_mgmt/v1/access_token",
		Method:   http.MethodPost,
		Handler:  handlers.TestStaticEndpoint,
	},
	{
		TestName: "list-subscriptions",
		Path:     "/api/accounts_mgmt/v1/subscriptions",
		Method:   http.MethodGet,
		Handler:  handlers.TestStaticEndpoint,
	},
	{
		TestName: "access-review",
		Path:     "/api/authorizations/v1/access_review",
		Method:   http.MethodPost,
		Handler:  handlers.TestStaticEndpoint,
		Body:     accessReviewBody(),
	},
	{
		TestName: "register-new-cluster",
		Path:     "/api/accounts_mgmt/v1/cluster_registrations",
		Method:   http.MethodPost,
		Handler:  handlers.TestRegisterNewCluster,
	},
	{
		TestName: "create-cluster",
		Path:     "/api/clusters_mgmt/v1/clusters",
		Method:   http.MethodPost,
		Handler:  handlers.TestCreateCluster,
	},
	{
		TestName: "list-clusters",
		Path:     "/api/clusters_mgmt/v1/clusters",
		Method:   http.MethodGet,
		Handler:  handlers.TestListClusters,
	},
	{
		TestName: "get-current-account",
		Path:     "/api/accounts_mgmt/v1/current_account",
		Method:   http.MethodGet,
		Handler:  handlers.TestStaticEndpoint,
	},
}

func accessReviewBody() []byte {
	buff := &bytes.Buffer{}
	resourceReviewReq, err := authv1.NewAccessReviewRequest().
		AccountUsername(helpers.AccountUsername).
		Action("get").
		ResourceType("Subscription").
		Build()
	if err != nil {
		log.Printf("building `access-review` reques: %s", err)
		return buff.Bytes()
	}
	err = authv1.MarshalAccessReviewRequest(resourceReviewReq, buff)
	if err != nil {
		log.Printf("marshaling `access-review` reques: %s", err)
	}
	return buff.Bytes()
}
