package provider

import (
	"net/http"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers/privateip"
)

type ConfigurationData struct {
	httpClient      *http.Client
	localInterfaces privateip.LocalInterfaces
}
