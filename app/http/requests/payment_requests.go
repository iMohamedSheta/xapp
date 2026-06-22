package requests

import "github.com/imohamedsheta/xapp/app/domain/requests"

type CreateOrUpdatePaymentPaymobRequest struct {
	PaymentGateway string `json:"payment_gateway" validate:"required,eq=paymob"`
	PaymentMethod  string `json:"payment_method" validate:"required,oneof=credit mobile"`
	ApiKey         string `json:"api_key" validate:"required"`
	ApiSecret      string `json:"api_secret" validate:"required"`
	HmacSecret     string `json:"hmac_secret" validate:"required"`
	Currency       string `json:"currency" validate:"required,oneof=EGP"`
	IntegrationId  int64  `json:"integration_id" validate:"required"`
	IframeId       int64  `json:"iframe_id" validate:"required_if=PaymentMethod credit"`
	requests.Request
}

func (r *CreateOrUpdatePaymentPaymobRequest) Messages() map[string]string {
	return map[string]string{
		"payment_gateway.required": "بوابة الدفع مطلوبة",
		"payment_method.required":  "طريقة الدفع مطلوبة",
		"payment_method.oneof":     "طريقة الدفع يجب أن تكون بطاقة أو موبايل",
		"api_secret.required":      "api secret مطلوب",
		"hmac_secret.required":     "hmac secret مطلوب",
		"api_key.required":         "api key مطلوب",
		"currency.required":        "العملة مطلوبة",
		"currency.oneof":           "العملة يجب أن تكون EGP",
		"integration_id.required":  "integration id مطلوب",
		"iframe_id.required":       "iframe id مطلوب",
		"iframe_id.required_if":    "iframe id مطلوب إذا كانت طريقة الدفع بطاقة",
	}
}
