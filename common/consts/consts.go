package consts

const DbService = "db"
const RavenClient = "raven_client"
const LoggerService = "logger"
const DefaultCountryCode = "BD"

var IsGinInitialized = false

const RedisService = "redis"
const RedisV8DB = "redis_client_v8"
const Cache = "cache"

const GormDbService = "gorm_db"

const ValidAttachmentType1 = "stock"
const ValidAttachmentType2 = "warehouse"
const ValidAttachmentType3 = "groupitem"

const CommonService = "service.common"
const CommonController = "controller.common"

const AccountTypePersonal = 1           //"PERSONAL"
const AccountTypeBusinessMainBook = 2   //"BUSINESS_MAIN_BOOK"
const AccountTypeBusinessBranchBook = 3 //"BUSINESS_BRANCH_BOOK"
const AccountTypeBusinessSubBook = 4    //"BUSINESS_SUB_BOOK"

var AccountTypesToString = map[int8]string{
	AccountTypePersonal:           "PERSONAL",
	AccountTypeBusinessMainBook:   "BUSINESS_MAIN_BOOK",
	AccountTypeBusinessBranchBook: "BUSINESS_BRANCH_BOOK",
	AccountTypeBusinessSubBook:    "BUSINESS_SUB_BOOK",
}

const AccountPermissionAdmin = "ADMIN"
const AccountPermissionOperator = "OPERATOR"
const AccountPermissionApprover = "APPROVER"
const AccountPermissionSuspended = "SUSPENDED"

var AccountTypes = map[string]int8{
	"PERSONAL":             AccountTypePersonal,
	"BUSINESS_MAIN_BOOK":   AccountTypeBusinessMainBook,
	"BUSINESS_BRANCH_BOOK": AccountTypeBusinessBranchBook,
	"BUSINESS_SUB_BOOK":    AccountTypeBusinessSubBook,
}

const RedisCommonRepository = "common.redis_repository"
const RedisDB = "redis_client"
const RedisCommonService = "common.redis_service"
const DefaultNotificationLimit = 50

const ActivityLogRepository = "activity_log_repository"
const ActivityLogService = "activity_log_service"
const ActivityLogHandler = "activity_log_handler"
const NotificationService = "service.notification"
const NotificationController = "controller.notification"

const ActionCreate = "Create"
const ActionBatchCreate = "BatchCreate"
const ActionDelete = "Delete"
const ActionBatchDelete = "BatchDelete"
const ActionUpdate = "Update"
const ActionNotification = "Notification"
