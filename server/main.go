
package main

import (
	"crypto/x509"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"net/http"

	"io/ioutil"

	"encoding/base64"
	"encoding/xml"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

func main() {

	// Get metadata from Okta
	res, err := http.Get("https://wawhal.okta.com/app/exksglejxVvOlrFYY5d6/sso/saml/metadata")
	if err != nil {
		panic(err)
	}
	rawMetadata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Read and validate the key descriptors
	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		panic(err)
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				panic(fmt.Errorf("metadata certificate(%d) must not be empty", idx))
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				panic(err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				panic(err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// Sign the AuthnRequest with a random key for testing
	randomKeyStore := dsig.RandomKeyStoreForTest()

	// Create a service provider with appropriate metadata
	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       "http://ui.saml.test",
		AssertionConsumerServiceURL: "http://server.saml.test/v1/_saml_callback",
		SignAuthnRequests:           true,
		AudienceURI:                 "http://server.saml.test/v1/_saml_callback",
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}

	// URL for performing authentication. UI redirects to this URL for SSO.	
	authURL, err := sp.BuildAuthURL("")
	if err != nil {
		panic(err)
	}

	// HTTP Handlers

	// Handler to receive assertion
	http.HandleFunc("/v1/_saml_callback", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse the assertion sent by Okta
		assertionInfo, err := sp.RetrieveAssertionInfo(req.FormValue("SAMLResponse"))
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		// Check for request request/response compatibility
		if assertionInfo.WarningInfo.InvalidTime {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		if assertionInfo.WarningInfo.NotInAudience {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		// Set the appropriate session cookie (this example just sets logged_in=true&<email_address>)
		cookie := http.Cookie {
			Name: "logged_in",
			Value: "true&" + assertionInfo.NameID,
			Expires: time.Now().Add(10*time.Minute),
			Domain: "saml.test",
		}
		http.SetCookie(rw, &cookie)

		// Redirect to the main app
		http.Redirect(rw, req, "http://ui.saml.test", http.StatusSeeOther)
	})

	// handler for the UI to get the login URI
	http.HandleFunc("/v1/login_uri", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		response := map[string]interface{}{
		    "uri": authURL,
		}
		json.NewEncoder(rw).Encode(response)	
	})

	// handler to check if the user is logged in; dummy code just checks for 
	http.HandleFunc("/v1/is_logged_in", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
		    "logged_in": false,
		}
		loggedInCookie, err := req.Cookie("logged_in")
		if err != nil {
			fmt.Println("err1")
			fmt.Println(err)
			response = map[string]interface{}{
			    "logged_in": false,
			}
		} else {
			cookieValue := loggedInCookie.Value
			response = map[string]interface{}{
			    "logged_in": strings.Contains(string(cookieValue), "true"),
			}
		}
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		fmt.Println("returning")
		json.NewEncoder(rw).Encode(response)	
	})

	// handler to logout user, remove the cookie (this example just sets logged_in=false)
	http.HandleFunc("/v1/logout", func(rw http.ResponseWriter, req *http.Request) {
		cookie := http.Cookie {
			Name: "logged_in",
			Value: "false",
			Domain: "saml.test",
		}
		http.SetCookie(rw, &cookie)
		fmt.Fprintf(rw, "OK")
	})

	println("Starting server on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
