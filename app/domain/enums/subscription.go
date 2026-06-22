package enums

// SubscriptionStatus is the status of the subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
)

// BillingCycle is the billing cycle of the subscription
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
	BillingCycleOneTime BillingCycle = "onetime"
)
