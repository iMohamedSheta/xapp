package enums

type PlanExpireAction string

const (
	PlanExpireActionBlock     PlanExpireAction = "block"
	PlanExpireActionDownGrade PlanExpireAction = "downgrade"
)
