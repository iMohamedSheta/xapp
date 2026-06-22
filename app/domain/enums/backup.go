package enums

type BackupSchedule string

const (
	BackupScheduleDaily   BackupSchedule = "daily"
	BackupScheduleWeekly  BackupSchedule = "weekly"
	BackupScheduleMonthly BackupSchedule = "monthly"
	BackupScheduleCustom  BackupSchedule = "custom"
	BackupSchedule12h     BackupSchedule = "12h"
)
