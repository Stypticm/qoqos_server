package models

import (
	"time"

	"gorm.io/datatypes"
	"github.com/lib/pq"
)

// Enums
type SkupkaStatus string

const (
	SkupkaDraft      SkupkaStatus = "draft"
	SkupkaOnTheWay   SkupkaStatus = "on_the_way"
	SkupkaInProgress SkupkaStatus = "in_progress"
	SkupkaAccepted   SkupkaStatus = "accepted"
	SkupkaPaid       SkupkaStatus = "paid"
	SkupkaCompleted  SkupkaStatus = "completed"
	SkupkaSubmitted  SkupkaStatus = "submitted"
	SkupkaInspected  SkupkaStatus = "inspected"
)

type MarketplaceLotStatus string

const (
	LotDraft     MarketplaceLotStatus = "draft"
	LotAvailable MarketplaceLotStatus = "available"
	LotReserved  MarketplaceLotStatus = "reserved"
	LotSold       MarketplaceLotStatus = "sold"
	LotArchived  MarketplaceLotStatus = "archived"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderInDelivery OrderStatus = "in_delivery"
	OrderCompleted  OrderStatus = "completed"
	OrderCancelled  OrderStatus = "cancelled"
)

type ChatStatus string

const (
	ChatActive   ChatStatus = "active"
	ChatArchived ChatStatus = "archived"
)

type Role string

const (
	RoleUser    Role = "USER"
	RoleAdmin   Role = "ADMIN"
	RoleMaster  Role = "MASTER"
	RoleManager Role = "MANAGER"
	RoleCourier Role = "COURIER"
)

type SenderType string

const (
	SenderUser  SenderType = "user"
	SenderAdmin SenderType = "admin"
)

type RepairStatus string

const (
	RepairCreated        RepairStatus = "created"
	RepairCourierAssigned RepairStatus = "courier_assigned"
	RepairInTransit      RepairStatus = "in_transit"
	RepairReceived        RepairStatus = "received"
	RepairUnpacked       RepairStatus = "unpacked"
	RepairDiagnosing     RepairStatus = "diagnosing"
	RepairPriceApproval  RepairStatus = "price_approval"
	RepairRepairing      RepairStatus = "repairing"
	RepairCompleted      RepairStatus = "completed"
	RepairReadyForPickup RepairStatus = "ready_for_pickup"
	RepairDelivered      RepairStatus = "delivered"
	RepairCancelled      RepairStatus = "cancelled"
)

// Models

type Skupka struct {
	ID                   string         `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId           string         `gorm:"column:telegramId" json:"telegramId"`
	Username             string         `gorm:"column:username" json:"username"`
	Modelname            *string        `gorm:"column:modelname" json:"modelname"`
	Imei                 *string        `gorm:"column:imei" json:"imei"`
	SN                   *string        `gorm:"column:sn" json:"sn"`
	Status               SkupkaStatus   `gorm:"type:text;default:draft;column:status" json:"status"`
	CreatedAt            time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	PhotoUrls            pq.StringArray `gorm:"type:text[];column:photoUrls" json:"photoUrls"`
	VideoUrls            pq.StringArray `gorm:"type:text[];column:videoUrls" json:"videoUrls"`
	DeviceData           datatypes.JSON `gorm:"column:deviceData" json:"deviceData"`
	AIAnalysis           datatypes.JSON `gorm:"column:aiAnalysis" json:"aiAnalysis"`
	PriceRange           datatypes.JSON `gorm:"column:priceRange" json:"priceRange"`
	AIModelUsed          *string        `gorm:"column:aiModelUsed" json:"aiModelUsed"`
	AnalysisConfidence   *float64       `gorm:"column:analysisConfidence" json:"analysisConfidence"`
	ChatHistory          datatypes.JSON `gorm:"column:chatHistory" json:"chatHistory"`
	DeviceConditions     datatypes.JSON `gorm:"column:deviceConditions" json:"deviceConditions"`
	Price                *float64       `gorm:"column:price" json:"price"`
	FinalPrice           *float64       `gorm:"column:finalPrice" json:"finalPrice"`
	PriceAgreed          *bool          `gorm:"column:priceAgreed" json:"priceAgreed"`
	DamagePercent        float64        `gorm:"default:0;column:damagePercent" json:"damagePercent"`
	FunctionDiscount     *float64       `gorm:"default:0;column:functionDiscount" json:"functionDiscount"`
	Courier              datatypes.JSON `gorm:"column:courier" json:"courier"`
	DeliveryMethod       *string        `gorm:"column:deliveryMethod" json:"deliveryMethod"`
	PickupPointPoint     *string        `gorm:"column:pickupPoint" json:"pickupPoint"`
	CurrentStep          *string        `gorm:"column:currentStep" json:"currentStep"`
	Inspection           datatypes.JSON `gorm:"column:inspection" json:"inspection"`
	InspectionCompleted  bool           `gorm:"default:false;column:inspectionCompleted" json:"inspectionCompleted"`
	StatusNote           *string        `gorm:"column:statusNote" json:"statusNote"`
	SubmittedAt          *time.Time     `gorm:"column:submittedAt" json:"submittedAt"`
	PriceConfirmed       bool           `gorm:"default:false;column:priceConfirmed" json:"priceConfirmed"`
	CourierReminderSent  bool           `gorm:"default:false;column:courierReminderSent" json:"courierReminderSent"`
	CourierUserConfirmed bool           `gorm:"default:false;column:courierUserConfirmed" json:"courierUserConfirmed"`
	Comment              *string        `gorm:"column:comment" json:"comment"`
	Feedback             *string        `gorm:"column:feedback" json:"feedback"`
	UserEvaluation       *string        `gorm:"column:userEvaluation" json:"userEvaluation"`
	AdditionalConditions datatypes.JSON `gorm:"column:additionalConditions" json:"additionalConditions"`
	AssignedMasterId     *string        `gorm:"column:assignedMasterId" json:"assignedMasterId"`
	CourierId            *string        `gorm:"column:courierId" json:"courierId"`

	// Relations
	User           *User   `gorm:"foreignKey:TelegramId;references:TelegramId" json:"user,omitempty"`
	AssignedMaster *Master `gorm:"foreignKey:AssignedMasterId;references:ID" json:"assignedMaster,omitempty"`
}

func (Skupka) TableName() string {
	return "Skupka"
}

type Master struct {
	ID         string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId string    `gorm:"unique;column:telegramId" json:"telegramId"`
	Username   string    `gorm:"unique;column:username" json:"username"`
	Name       *string   `gorm:"column:name" json:"name"`
	IsActive   bool      `gorm:"default:true;column:isActive" json:"isActive"`
	CreatedAt  time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	PointId    *int      `gorm:"column:pointId" json:"pointId"`

	// Relations
	Point *Point `gorm:"foreignKey:PointId;references:ID" json:"point,omitempty"`
}

func (Master) TableName() string {
	return "Master"
}

type Point struct {
	ID           int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Address      string `gorm:"unique;column:address" json:"address"`
	WorkingHours string `gorm:"default:'10:00 - 22:00';column:workingHours" json:"workingHours"`
	Name         string `gorm:"default:'Точка приёма';column:name" json:"name"`
}

func (Point) TableName() string {
	return "Point"
}

type DeviceInspection struct {
	ID              string         `gorm:"primaryKey;type:text;column:id" json:"id"`
	SkupkaId        string         `gorm:"column:skupkaId" json:"skupkaId"`
	InspectionToken string         `gorm:"column:inspectionToken" json:"inspectionToken"`
	TokenExpiresAt  time.Time      `gorm:"column:tokenExpiresAt" json:"tokenExpiresAt"`
	TestsResults    datatypes.JSON `gorm:"column:testsResults" json:"testsResults"`
	FinalPrice      *float64       `gorm:"column:finalPrice" json:"finalPrice"`
	InspectionNotes *string        `gorm:"column:inspectionNotes" json:"inspectionNotes"`
	CompletedAt     *time.Time     `gorm:"column:completedAt" json:"completedAt"`
	CreatedAt       time.Time      `gorm:"column:createdAt" json:"createdAt"`
	MasterUsername  string         `gorm:"column:masterUsername" json:"masterUsername"`
}

func (DeviceInspection) TableName() string {
	return "DeviceInspection"
}

type Device struct {
	ID        string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	Model     string    `gorm:"column:model" json:"model"`
	Variant   string    `gorm:"column:variant" json:"variant"`
	Storage   string    `gorm:"column:storage" json:"storage"`
	Color     string    `gorm:"column:color" json:"color"`
	Country   string    `gorm:"column:country" json:"country"`
	SimType   string    `gorm:"column:simType" json:"simType"`
	BasePrice float64   `gorm:"column:basePrice" json:"basePrice"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Device) TableName() string {
	return "Device"
}

type MarketPrice struct {
	ID          string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	DeviceId    string    `gorm:"column:deviceId" json:"deviceId"`
	Source      string    `gorm:"column:source" json:"source"`
	Price       float64   `gorm:"column:price" json:"price"`
	URL         *string   `gorm:"column:url" json:"url"`
	Title       *string   `gorm:"column:title" json:"title"`
	Description *string   `gorm:"column:description" json:"description"`
	Location    *string   `gorm:"column:location" json:"location"`
	Condition   *string   `gorm:"column:condition" json:"condition"`
	SellerType  *string   `gorm:"column:sellerType" json:"sellerType"`
	ParsedAt    time.Time `gorm:"default:now();column:parsedAt" json:"parsedAt"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (MarketPrice) TableName() string {
	return "MarketPrice"
}

type MarketplaceLot struct {
	ID          string               `gorm:"primaryKey;type:text;column:id" json:"id"`
	SkupkaId    *string              `gorm:"unique;column:skupkaId" json:"skupkaId"`
	Title       string               `gorm:"column:title" json:"title"`
	Model       *string              `gorm:"column:model" json:"model"`
	Storage     *string              `gorm:"column:storage" json:"storage"`
	Color       *string              `gorm:"column:color" json:"color"`
	Brand       *string              `gorm:"column:brand" json:"brand"`
	SKU         *string              `gorm:"unique;column:sku" json:"sku"`
	Condition   *string              `gorm:"column:condition" json:"condition"`
	Description *string              `gorm:"column:description" json:"description"`
	Price       float64              `gorm:"column:price" json:"price"`
	OldPrice    *float64             `gorm:"column:oldPrice" json:"oldPrice"`
	Photos      pq.StringArray       `gorm:"type:text[];column:photos" json:"photos"`
	CoverPhoto  *string              `gorm:"column:coverPhoto" json:"coverPhoto"`
	Status      MarketplaceLotStatus `gorm:"type:text;default:available;column:status" json:"status"`
	TelegramId  string               `gorm:"column:telegramId" json:"telegramId"`
	SellerName  *string              `gorm:"column:sellerName" json:"sellerName"`
	ViewsCount  int                  `gorm:"default:0;column:viewsCount" json:"viewsCount"`
	IsAccessory bool                 `gorm:"default:false;column:isAccessory" json:"isAccessory"`
	TargetBrand *string              `gorm:"column:targetBrand" json:"targetBrand"`
	TargetModel *string              `gorm:"column:targetModel" json:"targetModel"`
	CreatedAt   time.Time            `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `gorm:"column:updatedAt" json:"updatedAt"`
	PublishedAt *time.Time           `gorm:"column:publishedAt" json:"publishedAt"`
	SoldAt      *time.Time           `gorm:"column:soldAt" json:"soldAt"`
}

func (MarketplaceLot) TableName() string {
	return "MarketplaceLot"
}

type CartItem struct {
	ID         string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId string    `gorm:"uniqueIndex:idx_user_lot;column:telegramId" json:"telegramId"`
	LotId      string    `gorm:"uniqueIndex:idx_user_lot;column:lotId" json:"lotId"`
	Quantity   int       `gorm:"default:1;column:quantity" json:"quantity"`
	AddedAt    time.Time `gorm:"column:addedAt" json:"addedAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (CartItem) TableName() string {
	return "CartItem"
}

type FavoriteItem struct {
	ID         string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId string    `gorm:"uniqueIndex:idx_fav_user_lot;column:telegramId" json:"telegramId"`
	LotId      string    `gorm:"uniqueIndex:idx_fav_user_lot;column:lotId" json:"lotId"`
	AddedAt    time.Time `gorm:"column:addedAt" json:"addedAt"`
}

func (FavoriteItem) TableName() string {
	return "FavoriteItem"
}

type Order struct {
	ID              string      `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId      string      `gorm:"column:telegramId" json:"telegramId"`
	DeliveryMethod  string      `gorm:"column:deliveryMethod" json:"deliveryMethod"`
	DeliveryAddress *string     `gorm:"column:deliveryAddress" json:"deliveryAddress"`
	PickupPointId   *int        `gorm:"column:pickupPointId" json:"pickupPointId"`
	DeliveryDate    *time.Time  `gorm:"column:deliveryDate" json:"deliveryDate"`
	DeliveryTime    *string     `gorm:"column:deliveryTime" json:"deliveryTime"`
	Status          OrderStatus `gorm:"type:text;default:pending;column:status" json:"status"`
	TotalPrice      float64     `gorm:"column:totalPrice" json:"totalPrice"`
	ConfirmedAt     *time.Time  `gorm:"column:confirmedAt" json:"confirmedAt"`
	InDeliveryAt    *time.Time  `gorm:"column:inDeliveryAt" json:"inDeliveryAt"`
	CompletedAt     *time.Time  `gorm:"column:completedAt" json:"completedAt"`
	CourierId       *string     `gorm:"column:courierId" json:"courierId"`
	CourierName     *string     `gorm:"column:courierName" json:"courierName"`
	CourierPhone    *string     `gorm:"column:courierPhone" json:"courierPhone"`
	TrackingNotes   *string     `gorm:"column:trackingNotes" json:"trackingNotes"`
	CreatedAt       time.Time   `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time   `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Order) TableName() string {
	return "Order"
}

type OrderItem struct {
	ID      string  `gorm:"primaryKey;type:text;column:id" json:"id"`
	OrderId string  `gorm:"column:orderId" json:"orderId"`
	LotId   string  `gorm:"column:lotId" json:"lotId"`
	Title   string  `gorm:"column:title" json:"title"`
	Price   float64 `gorm:"column:price" json:"price"`
}

func (OrderItem) TableName() string {
	return "OrderItem"
}

type OperatorChat struct {
	ID           string     `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId   string     `gorm:"unique;column:telegramId" json:"telegramId"`
	UserNickname *string    `gorm:"column:userNickname" json:"userNickname"`
	Status       ChatStatus `gorm:"type:text;default:active;column:status" json:"status"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

func (OperatorChat) TableName() string {
	return "OperatorChat"
}

type OperatorMessage struct {
	ID         string     `gorm:"primaryKey;type:text;column:id" json:"id"`
	ChatId     string     `gorm:"column:chatId" json:"chatId"`
	SenderId   string     `gorm:"column:senderId" json:"senderId"`
	SenderType SenderType `gorm:"type:text;column:senderType" json:"senderType"`
	Text       string     `gorm:"column:text" json:"text"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
}

func (OperatorMessage) TableName() string {
	return "OperatorMessage"
}

type AuthRequest struct {
	ID               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id" json:"id"`
	Status           string         `gorm:"default:pending;column:status" json:"status"`
	TelegramId       *string        `gorm:"column:telegram_id" json:"telegramId"`
	TelegramUsername *string        `gorm:"column:telegram_username" json:"telegramUsername"`
	TelegramData     datatypes.JSON `gorm:"column:telegram_data" json:"telegramData"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"createdAt"`
}

func (AuthRequest) TableName() string {
	return "auth_requests"
}

type User struct {
	ID           string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId   string    `gorm:"unique;column:telegramId" json:"telegramId"`
	PasswordHash string    `gorm:"column:passwordHash" json:"-"`
	Role         Role      `gorm:"type:text;default:USER;column:role" json:"role"`
	CreatedAt    time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (User) TableName() string {
	return "User"
}

type QuickLead struct {
	ID           string     `gorm:"primaryKey;type:text;column:id" json:"id"`
	Name         string     `gorm:"column:name" json:"name"`
	Phone        string     `gorm:"column:phone" json:"phone"`
	ProductId    *string    `gorm:"column:productId" json:"productId"`
	ProductTitle *string    `gorm:"column:productTitle" json:"productTitle"`
	Price        *float64   `gorm:"column:price" json:"price"`
	TelegramId   *string    `gorm:"column:telegramId" json:"telegramId"`
	Status       string     `gorm:"default:new;column:status" json:"status"`
	Address      *string    `gorm:"column:address" json:"address"`
	DeliveryDate *time.Time `gorm:"column:deliveryDate" json:"deliveryDate"`
	DeliveryTime *string    `gorm:"column:deliveryTime" json:"deliveryTime"`
	IsRead       bool       `gorm:"default:false;column:isRead" json:"isRead"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

func (QuickLead) TableName() string {
	return "QuickLead"
}

type BlogPost struct {
	ID        string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	Content   string    `gorm:"type:text;column:content" json:"content"`
	Excerpt   *string   `gorm:"type:text;column:excerpt" json:"excerpt"`
	Image     *string   `gorm:"column:image" json:"image"`
	Category  string    `gorm:"default:'Новости';column:category" json:"category"`
	Author    *string   `gorm:"column:author" json:"author"`
	Published bool      `gorm:"default:false;column:published" json:"published"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (BlogPost) TableName() string {
	return "BlogPost"
}

type TradeInEvaluation struct {
	ID              string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId      string    `gorm:"column:telegramId" json:"telegramId"`
	Category        string    `gorm:"default:'iPhone';column:category" json:"category"`
	Model           string    `gorm:"column:model" json:"model"`
	Variant         *string   `gorm:"column:variant" json:"variant"`
	Storage         string    `gorm:"column:storage" json:"storage"`
	Color           string    `gorm:"column:color" json:"color"`
	IsOriginal      bool      `gorm:"default:true;column:isOriginal" json:"isOriginal"`
	IsReset         bool      `gorm:"default:true;column:isReset" json:"isReset"`
	ScreenCondition string    `gorm:"column:screenCondition" json:"screenCondition"`
	BodyCondition   string    `gorm:"column:bodyCondition" json:"bodyCondition"`
	IsRostest       bool      `gorm:"default:true;column:isRostest" json:"isRostest"`
	BatteryHealth   string    `gorm:"column:batteryHealth" json:"batteryHealth"`
	HasFullSet      bool      `gorm:"default:true;column:hasFullSet" json:"hasFullSet"`
	WasRepaired     bool      `gorm:"default:false;column:wasRepaired" json:"wasRepaired"`
	HasReceipt      bool      `gorm:"default:false;column:hasReceipt" json:"hasReceipt"`
	IsFunctional    bool      `gorm:"default:true;column:isFunctional" json:"isFunctional"`
	IsBatterySafe   bool      `gorm:"default:true;column:isBatterySafe" json:"isBatterySafe"`
	IsHardwareOk    bool      `gorm:"default:true;column:isHardwareOk" json:"isHardwareOk"`
	IsClean         bool      `gorm:"default:true;column:isClean" json:"isClean"`
	CalculatedPrice float64   `gorm:"column:calculatedPrice" json:"calculatedPrice"`
	MinPrice        *float64  `gorm:"column:minPrice" json:"minPrice"`
	MaxPrice        *float64  `gorm:"column:maxPrice" json:"maxPrice"`
	Status          string    `gorm:"default:draft;column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (TradeInEvaluation) TableName() string {
	return "TradeInEvaluation"
}

type PushSubscription struct {
	ID         string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId *string   `gorm:"column:telegramId" json:"telegramId"`
	Endpoint   string    `gorm:"unique;column:endpoint" json:"endpoint"`
	P256dh     string    `gorm:"column:p256dh" json:"p256dh"`
	Auth       string    `gorm:"column:auth" json:"auth"`
	CreatedAt  time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (PushSubscription) TableName() string {
	return "PushSubscription"
}

type BotAccess struct {
	ID              string     `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId      string     `gorm:"unique;column:telegramId" json:"telegramId"`
	IsAuthenticated bool       `gorm:"default:false;column:isAuthenticated" json:"isAuthenticated"`
	Attempts        int        `gorm:"default:0;column:attempts" json:"attempts"`
	BlockedUntil    *time.Time `gorm:"column:blockedUntil" json:"blockedUntil"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
}

func (BotAccess) TableName() string {
	return "BotAccess"
}

type RepairRequest struct {
	ID               string         `gorm:"primaryKey;type:text;column:id" json:"id"`
	TelegramId       string         `gorm:"column:telegramId" json:"telegramId"`
	DeviceModel      string         `gorm:"column:deviceModel" json:"deviceModel"`
	SerialNumber     *string        `gorm:"column:serialNumber" json:"serialNumber"`
	Category         string         `gorm:"column:category" json:"category"`
	IssueDescription *string        `gorm:"column:issueDescription" json:"issueDescription"`
	IssuePhotos      pq.StringArray `gorm:"type:text[];column:issuePhotos" json:"issuePhotos"`
	EstimatedMin     *float64       `gorm:"column:estimatedMin" json:"estimatedMin"`
	EstimatedMax     *float64       `gorm:"column:estimatedMax" json:"estimatedMax"`
	FinalPrice       *float64       `gorm:"column:finalPrice" json:"finalPrice"`
	PriceOriginal    *float64       `gorm:"column:priceOriginal" json:"priceOriginal"`
	PriceNonOriginal *float64       `gorm:"column:priceNonOriginal" json:"priceNonOriginal"`
	ClientPriceChoice *string       `gorm:"column:clientPriceChoice" json:"clientPriceChoice"`
	ReturnMethod     *string        `gorm:"column:returnMethod" json:"returnMethod"`
	MasterNotes      *string        `gorm:"column:masterNotes" json:"masterNotes"`
	ClientContact    *string        `gorm:"column:clientContact" json:"clientContact"`
	ClientAddress    *string        `gorm:"column:clientAddress" json:"clientAddress"`
	DeliveryMethod   string         `gorm:"column:deliveryMethod" json:"deliveryMethod"`
	AppointmentDate  *time.Time     `gorm:"column:appointmentDate" json:"appointmentDate"`
	AppointmentTime  *string        `gorm:"column:appointmentTime" json:"appointmentTime"`
	Status           RepairStatus   `gorm:"type:text;default:created;column:status" json:"status"`
	UnpackVideoUrl   *string        `gorm:"column:unpackVideoUrl" json:"unpackVideoUrl"`
	EnvelopeNumber   *string        `gorm:"column:envelopeNumber" json:"envelopeNumber"`
	AssignedMasterId *string        `gorm:"column:assignedMasterId" json:"assignedMasterId"`
	CourierId        *string        `gorm:"column:courierId" json:"courierId"`
	CourierNotes     *string        `gorm:"column:courierNotes" json:"courierNotes"`
	CustomerNotes    *string        `gorm:"column:customerNotes" json:"customerNotes"`
	CreatedAt        time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"column:updatedAt" json:"updatedAt"`

	// Relations
	User           *User   `gorm:"foreignKey:TelegramId;references:TelegramId" json:"user,omitempty"`
	AssignedMaster *Master `gorm:"foreignKey:AssignedMasterId;references:ID" json:"assignedMaster,omitempty"`
}

func (RepairRequest) TableName() string {
	return "RepairRequest"
}

type AgentAuditLog struct {
	ID        string         `gorm:"primaryKey;type:text;column:id" json:"id"`
	RequestId *string        `gorm:"column:requestId" json:"requestId"`
	ChatId    string         `gorm:"column:chatId" json:"chatId"`
	AgentName string         `gorm:"column:agentName" json:"agentName"`
	Action    string         `gorm:"column:action" json:"action"`
	Input     datatypes.JSON `gorm:"column:input" json:"input"`
	Output    datatypes.JSON `gorm:"column:output" json:"output"`
	Status    string         `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:createdAt" json:"createdAt"`
}

func (AgentAuditLog) TableName() string {
	return "AgentAuditLog"
}

type IdempotencyKey struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
}

func (IdempotencyKey) TableName() string {
	return "IdempotencyKey"
}
