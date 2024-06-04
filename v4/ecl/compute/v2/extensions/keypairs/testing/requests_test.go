package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nttcom/eclcloud/v4/ecl/compute/v2/extensions/keypairs"
	"github.com/nttcom/eclcloud/v4/pagination"

	th "github.com/nttcom/eclcloud/v4/testhelper"
	fakeclient "github.com/nttcom/eclcloud/v4/testhelper/client"
)

func TestListKeyPair(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, listOutput)
	})

	count := 0
	err := keypairs.List(fakeclient.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := keypairs.ExtractKeyPairs(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, expectedKeyPairSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreateKeyPair(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestJSONRequest(t, r, createRequest)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, createResponse)
	})

	actual, err := keypairs.Create(fakeclient.ServiceClient(), keypairs.CreateOpts{
		Name: "createdkey",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &createdKeyPair, actual)
}

func TestImportKeypair(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestJSONRequest(t, r, importRequest)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, importResponse)
	})

	actual, err := keypairs.Create(fakeclient.ServiceClient(), keypairs.CreateOpts{
		Name:      "importedkey",
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDx8nkQv/zgGgB4rMYmIf+6A4l6Rr+o/6lHBQdW5aYd44bd8JttDCE/F/pNRr0lRE+PiqSPO8nDPHw0010JeMH9gYgnnFlyY3/OcJ02RhIPyyxYpv9FhY+2YiUkpwFOcLImyrxEsYXpD/0d3ac30bNH6Sw9JD9UZHYcpSxsIbECHw== Generated by Nova",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &importedKeyPair, actual)
}

func TestGetKeyPair(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/os-keypairs/firstkey", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, getResponse)
	})

	actual, err := keypairs.Get(fakeclient.ServiceClient(), "firstkey").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &firstKeyPair, actual)
}

func TestDeleteKeyPair(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/os-keypairs/deletedkey", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	err := keypairs.Delete(fakeclient.ServiceClient(), "deletedkey").ExtractErr()
	th.AssertNoErr(t, err)
}
