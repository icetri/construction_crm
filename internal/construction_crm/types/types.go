package types

import (
	"database/sql"
	pq "github.com/lib/pq"
	"io"
	"time"
)

const (
	RoleUser    string = "USER"
	RoleManager string = "MANAGER"
	RoleAdmin   string = "ADMIN"
)

type Token struct {
	Token string `json:"token"`
}

type Register struct {
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	RegManager bool   `json:"registered_manager"`
}

type User struct {
	ID         int          `json:"id" db:"id"`
	Image      string       `json:"image" db:"image"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime `json:"-" db:"updated_at"`
	Phone      string       `json:"phone" db:"phone"`
	Email      string       `json:"email" db:"email"`
	FirstName  string       `json:"first_name" db:"first_name"`
	LastName   string       `json:"last_name" db:"last_name"`
	MiddleName string       `json:"middle_name" db:"middle_name"`
	Role       string       `json:"role"`
	RegManager bool         `json:"registered_manager" db:"registered_manager"`
	Code       string       `json:"-" db:"code"`
}

type AllUsers struct {
	ID         int          `json:"id" db:"id"`
	Image      string       `json:"image" db:"image"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime `json:"-" db:"updated_at"`
	Address    []Address    `json:"address"`
	Phone      string       `json:"phone" db:"phone"`
	Email      string       `json:"email" db:"email"`
	FirstName  string       `json:"first_name" db:"first_name"`
	LastName   string       `json:"last_name" db:"last_name"`
	MiddleName string       `json:"middle_name" db:"middle_name"`
	Role       string       `json:"role"`
	RegManager bool         `json:"registered_manager" db:"registered_manager"`
	Code       string       `json:"-" db:"code"`
}

type PaginationListUsers struct {
	ID         int          `json:"id" db:"id"`
	Image      string       `json:"image" db:"image"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime `json:"-" db:"updated_at"`
	Phone      string       `json:"phone" db:"phone"`
	Email      string       `json:"email" db:"email"`
	Address    []Address    `json:"address"`
	FirstName  string       `json:"first_name" db:"first_name"`
	LastName   string       `json:"last_name" db:"last_name"`
	MiddleName string       `json:"middle_name" db:"middle_name"`
	Role       string       `json:"role"`
	RegManager bool         `json:"registered_manager" db:"registered_manager"`
	Code       string       `json:"-" db:"code"`
	Hits       int64        `json:"projects"`
}

type Manager struct {
	ID         int          `json:"id" db:"id"`
	Image      string       `json:"image" db:"image"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime `json:"-" db:"updated_at"`
	Phone      string       `json:"phone" db:"phone"`
	Email      string       `json:"email" db:"email"`
	FirstName  string       `json:"first_name" db:"first_name"`
	LastName   string       `json:"last_name" db:"last_name"`
	MiddleName string       `json:"middle_name" db:"middle_name"`
	Role       string       `json:"role" db:"role"`
	Password   string       `json:"password" db:"password"`
	Country    string       `json:"country" db:"country"`
	City       string       `json:"city" db:"city"`
}

///PUT CABINET

type PutUserInfo struct {
	ID         int    `json:"-" db:"id"`
	Image      string `json:"image" db:"image"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	Email      string `json:"email" db:"email"`
}

type PutUserCont struct {
	ID    int    `json:"-" db:"id"`
	Phone string `json:"phone" db:"phone"`
	Code  string `json:"code" db:"code"`
}

type PutUserPhone struct {
	ID       int    `json:"-" db:"id"`
	NewPhone string `json:"phone" db:"phone"`
}

//// AUTH //////
type Auth struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Code      string `json:"code" db:"code"`
	Role      string `json:"role"`
}

type AuthPhone struct {
	Phone string `json:"phone"`
}

type AuthEmail struct {
	Email string `json:"email"`
}

type AuthCodePhone struct {
	Phone string `json:"phone"`
	Code  string `json:"code" db:"code"`
}

type AuthCodeEmail struct {
	Email string `json:"email"`
	Code  string `json:"code" db:"code"`
}

type AuthEmailManager struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistredUserByManagerProject struct {
	Id         int    `json:"id" db:"id"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	RegManager bool   `json:"-" db:"registered_manager"`
}

///CABINET

type Cabinet struct {
	ID         int    `json:"id" db:"id"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Image      string `json:"image" db:"image"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
}

type CabinetNew struct {
	ID         int       `json:"id" db:"id"`
	Phone      string    `json:"phone" db:"phone"`
	Email      string    `json:"email" db:"email"`
	Image      string    `json:"image" db:"image"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	MiddleName string    `json:"middle_name" db:"middle_name"`
	Addresses  []Address `json:"address"`
}

type Address struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Address     string `json:"address" db:"address"`
	City        string `json:"city" db:"city"`
	Country     string `json:"country" db:"country"`
	Entrance    string `json:"entrance" db:"entrance"`
	Description string `json:"description" db:"description"`
	UserId      int    `json:"user_id" db:"user_id"`
	ProjectID   int    `json:"project_id" db:"project_id"`
}

///PROJECTS
type Projects struct {
	Id                   int    `json:"id" db:"id"`
	Image                string `json:"image" db:"image"`
	File                 string `json:"file" db:"file"`
	Address              string `json:"address" db:"address"`
	UserId               int    `json:"user_id" db:"user_id"`
	StartDate            string `json:"start_date" db:"start_date"`
	EndDate              string `json:"end_date" db:"end_date"`
	Active               bool   `json:"active" db:"active"`
	Maker                int    `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int    `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int    `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int    `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int    `json:"work_cost_spent" db:"work_cost_spent"`
}

type ProjectsWithClient struct {
	Id                   int    `json:"id" db:"id"`
	Image                string `json:"image" db:"image"`
	File                 string `json:"file" db:"file"`
	Address              string `json:"address" db:"address"`
	UserId               int    `json:"user_id" db:"user_id"`
	StartDate            string `json:"start_date" db:"start_date"`
	EndDate              string `json:"end_date" db:"end_date"`
	Active               bool   `json:"active" db:"active"`
	Maker                int    `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int    `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int    `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int    `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int    `json:"work_cost_spent" db:"work_cost_spent"`
	Client               *User  `json:"client"`
}

/////////ДОБАВЛЕНИЕ ПРОЕКТА////////
type AddProject struct {
	Id                   int    `json:"-" db:"id"`
	Image                string `json:"image" db:"image"`
	File                 string `json:"file" db:"file"`
	Address              string `json:"address" db:"address"`
	UserId               int    `json:"user_id" db:"user_id"`
	StartDate            string `json:"start_date" db:"start_date"`
	EndDate              string `json:"end_date" db:"end_date"`
	Active               bool   `json:"active" db:"active"`
	Maker                int    `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int    `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int    `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int    `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int    `json:"work_cost_spent" db:"work_cost_spent"`
}

type AddStage struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ProjectId int    `json:"project_id" db:"project_id"`
	Phase     bool   `json:"phase" db:"phase"`
	Date      string `json:"date" db:"date"`
}

type AddCard struct {
	Id               int            `json:"-" db:"id"`
	Image            pq.StringArray `json:"images" db:"images"`
	Title            string         `json:"title" db:"title"`
	Deadline         string         `json:"deadline" db:"deadline"`
	StagesId         int            `json:"stages_id" db:"stages_id"`
	Rating           string         `json:"rating" db:"rating"`
	Description      string         `json:"description" db:"description"`
	LeftToPay        int            `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int            `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string         `json:"status" db:"status"`
	State            string         `json:"state" db:"state"`
	ProjectId        int            `json:"project_id"`
}

type AddTask struct {
	Id       int            `json:"-" db:"id"`
	Image    pq.StringArray `json:"images" db:"images"`
	Title    string         `json:"title" db:"title"`
	CardId   int            `json:"card_id" db:"card_id"`
	Complete bool           `json:"complete" db:"complete"`
	Length   int            `json:"length" db:"length"`
}

type AddCheque struct {
	Id        int            `json:"_" db:"id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	File      pq.StringArray `json:"file" db:"file"`
	Name      string         `json:"name" db:"name"`
	Cost      int            `json:"cost" db:"cost"`
	Type      string         `json:"type" db:"type"`
	CardId    int            `json:"-" db:"card_id"`
	UserId    int            `json:"-" db:"user_id"`
	ProjectId int            `json:"-" db:"project_id"`
	Length    int            `json:"length" db:"length"`
	Weight    pq.Int64Array  `json:"weight" db:"weight"`
}

type DeleteCheque struct {
	Id        int `json:"id" db:"id"`
	CardId    int `json:"-" db:"card_id"`
	ProjectId int `json:"-" db:"project_id"`
}

/////////К ПРОЕКТУ/////////
type Project struct {
	Id                   int      `json:"id" db:"id"`
	Image                string   `json:"image" db:"image"`
	File                 string   `json:"file" db:"file"`
	Address              string   `json:"address" db:"address"`
	UserId               int      `json:"user_id" db:"user_id"`
	StartDate            string   `json:"start_date" db:"start_date"`
	EndDate              string   `json:"end_date" db:"end_date"`
	Active               bool     `json:"active" db:"active"`
	Maker                int      `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int      `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int      `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int      `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int      `json:"work_cost_spent" db:"work_cost_spent"`
	Stages               []Stages `json:"stages"`
}

type ProjectManager struct {
	Id                   int                  `json:"id" db:"id"`
	Image                string               `json:"image" db:"image"`
	File                 string               `json:"file" db:"file"`
	Address              string               `json:"address" db:"address"`
	UserId               int                  `json:"user_id" db:"user_id"`
	StartDate            string               `json:"start_date" db:"start_date"`
	EndDate              string               `json:"end_date" db:"end_date"`
	Active               bool                 `json:"active" db:"active"`
	Maker                int                  `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int                  `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int                  `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int                  `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int                  `json:"work_cost_spent" db:"work_cost_spent"`
	Client               *PaginationListUsers `json:"client"`
	Stages               []StagesWithCards    `json:"stages"`
	Cheques              []Cheque             `json:"cheques"`
}

type Stages struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ProjectId int    `json:"project_id" db:"project_id"`
	Phase     bool   `json:"phase" db:"phase"`
	Date      string `json:"date" db:"date"`
	Cards     []Card `json:"cards"`
}

type StagesWithCards struct {
	Id        int             `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	ProjectId int             `json:"project_id" db:"project_id"`
	Phase     bool            `json:"phase" db:"phase"`
	Date      string          `json:"date" db:"date"`
	Cards     []CardWithCount `json:"cards"`
}

type StagesCards struct {
	Id        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	ProjectId int     `json:"project_id" db:"project_id"`
	Phase     bool    `json:"phase" db:"phase"`
	Date      string  `json:"date" db:"date"`
	Cards     []Cards `json:"cards"`
}

type Stage struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ProjectId int    `json:"project_id" db:"project_id"`
	Phase     bool   `json:"phase" db:"phase"`
	Date      string `json:"date" db:"date"`
}

type Cards struct {
	Id               int            `json:"id" db:"id"`
	Images           pq.StringArray `json:"images" db:"images"`
	Title            string         `json:"title" db:"title"`
	Deadline         string         `json:"deadline" db:"deadline"`
	StagesId         int            `json:"stages_id" db:"stages_id"`
	Rating           string         `json:"rating" db:"rating"`
	Description      string         `json:"description" db:"description"`
	LeftToPay        int            `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int            `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string         `json:"status" db:"status"`
	State            string         `json:"state" db:"state"`
	Tasks            []Task         `json:"tasks"`
	Cheques          []Cheque       `json:"cheques"`
}

type Card struct {
	Id               int            `json:"id" db:"id"`
	Images           pq.StringArray `json:"images" db:"images"`
	Title            string         `json:"title" db:"title"`
	Deadline         string         `json:"deadline" db:"deadline"`
	StagesId         int            `json:"stages_id" db:"stages_id"`
	Rating           string         `json:"rating" db:"rating"`
	Description      string         `json:"description" db:"description"`
	LeftToPay        int            `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int            `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string         `json:"status" db:"status"`
	State            string         `json:"state" db:"state"`
}

type CardWithCount struct {
	Id               int            `json:"id" db:"id"`
	Images           pq.StringArray `json:"images" db:"images"`
	Title            string         `json:"title" db:"title"`
	Deadline         string         `json:"deadline" db:"deadline"`
	StagesId         int            `json:"stages_id" db:"stages_id"`
	Rating           string         `json:"rating" db:"rating"`
	Description      string         `json:"description" db:"description"`
	LeftToPay        int            `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int            `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string         `json:"status" db:"status"`
	State            string         `json:"state" db:"state"`
	Cheques          int            `json:"cheques"`
}

type Task struct {
	Id       int            `json:"id" db:"id"`
	Title    string         `json:"title" db:"title"`
	CardId   int            `json:"card_id" db:"card_id"`
	Complete bool           `json:"complete" db:"complete"`
	Images   pq.StringArray `json:"images" db:"images"`
	Length   int            `json:"length" db:"length"`
}

type Cheque struct {
	Id        int            `json:"id" db:"id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	File      pq.StringArray `json:"file" db:"file"`
	Name      string         `json:"name" db:"name"`
	Cost      int            `json:"cost" db:"cost"`
	Type      string         `json:"type" db:"type"`
	CardId    int            `json:"card_id" db:"card_id"`
	UserId    int            `json:"user_id" db:"user_id"`
	ProjectId int            `json:"project_id" db:"project_id"`
	Length    int            `json:"length" db:"length"`
	Weight    pq.Int64Array  `json:"weight" db:"weight"`
}

//////////Изменение прокета///////////////
type UpdateProject struct {
	Id                   int           `json:"id" db:"id"`
	Image                string        `json:"image" db:"image"`
	File                 string        `json:"file" db:"file"`
	Address              string        `json:"address" db:"address"`
	UserId               int           `json:"user_id" db:"user_id"`
	StartDate            string        `json:"start_date" db:"start_date"`
	EndDate              string        `json:"end_date" db:"end_date"`
	Active               bool          `json:"active" db:"active"`
	Maker                int           `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int           `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int           `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int           `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int           `json:"work_cost_spent" db:"work_cost_spent"`
	Client               *Client       `json:"client"`
	Stages               []UpdateStage `json:"stages"`
}

type ProjectManagerWithClient struct {
	Id                   int      `json:"id" db:"id"`
	Image                string   `json:"image" db:"image"`
	File                 string   `json:"file" db:"file"`
	Address              string   `json:"address" db:"address"`
	UserId               int      `json:"user_id" db:"user_id"`
	StartDate            string   `json:"start_date" db:"start_date"`
	EndDate              string   `json:"end_date" db:"end_date"`
	Active               bool     `json:"active" db:"active"`
	Maker                int      `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int      `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int      `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int      `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int      `json:"work_cost_spent" db:"work_cost_spent"`
	Client               *User    `json:"client"`
	Stages               []Stages `json:"stages"`
	Cheques              []Cheque `json:"cheques"`
}

type Client struct {
	Id         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
}

type UpdateStage struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ProjectId int    `json:"project_id" db:"project_id"`
	Phase     bool   `json:"phase" db:"phase"`
	Date      string `json:"date" db:"date"`
	Cards     []Card `json:"cards"`
}

type UpdateCard struct {
	Id               int            `json:"id" db:"id"`
	Images           pq.StringArray `json:"images" db:"images"`
	Title            string         `json:"title" db:"title"`
	Deadline         string         `json:"deadline" db:"deadline"`
	StagesId         int            `json:"stages_id" db:"stages_id"`
	Rating           string         `json:"rating" db:"rating"`
	Description      string         `json:"description" db:"description"`
	LeftToPay        int            `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int            `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string         `json:"status" db:"status"`
	State            string         `json:"state" db:"state"`
	ProjectId        int            `json:"project_id"`
	Tasks            []Task         `json:"tasks"`
	Cheques          []Cheque       `json:"cheques"`
}

/////////КАЛЕНДАРЬ//////////
type UserCalendar struct {
	ID         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
}

type ProjectCalendar struct {
	Id        int              `json:"id" db:"id"`
	Address   string           `json:"address" db:"address"`
	User      *UserCalendar    `json:"client"`
	StartDate string           `json:"start_date" db:"start_date"`
	EndDate   string           `json:"end_date" db:"end_date"`
	Active    bool             `json:"-" db:"active"`
	Maker     int              `json:"maker_id" db:"maker_id"`
	Stages    []StagesCalendar `json:"stages"`
}

type StagesCalendar struct {
	Id        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	ProjectId int            `json:"project_id" db:"project_id"`
	Phase     bool           `json:"phase" db:"phase"`
	Date      string         `json:"date" db:"date"`
	Cards     []CardCalendar `json:"cards"`
}

type CardCalendar struct {
	Id       int    `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	Deadline string `json:"deadline" db:"deadline"`
	StagesId int    `json:"stages_id" db:"stages_id"`
	Status   string `json:"status" db:"status"`
}

/////////ДЛЯ ЗАГРУЗКИ ФАЙЛОВ/////////////////
type UploadFile struct {
	ID   int
	Body io.Reader
	Head string
}

///////Добавление проекта и вывод//////
type AddProjectManager struct {
	Id                   int               `json:"id" db:"id"`
	Image                string            `json:"image" db:"image"`
	File                 string            `json:"file" db:"file"`
	Address              string            `json:"address" db:"address"`
	UserId               int               `json:"user_id" db:"user_id"`
	StartDate            string            `json:"start_date" db:"start_date"`
	EndDate              string            `json:"end_date" db:"end_date"`
	Active               bool              `json:"active" db:"active"`
	Maker                int               `json:"maker_id" db:"maker_id"`
	MaterialCostsOverall int               `json:"material_costs_over_all" db:"material_costs_over_all"`
	WorkCostsOverall     int               `json:"work_costs_over_all" db:"work_costs_over_all"`
	MaterialCostSpent    int               `json:"material_cost_spent" db:"material_cost_spent"`
	WorkCostSpent        int               `json:"work_cost_spent" db:"work_cost_spent"`
	AddStagesManager     []AddStageManager `json:"stages"`
}

type AddStageManager struct {
	Id              int              `json:"id" db:"id"`
	Name            string           `json:"name" db:"name"`
	ProjectId       int              `json:"project_id" db:"project_id"`
	Phase           bool             `json:"phase" db:"phase"`
	Date            string           `json:"date" db:"date"`
	AddCardsManager []AddCardManager `json:"cards"`
}

type AddCardManager struct {
	Id               int                `json:"id" db:"id"`
	Image            pq.StringArray     `json:"images" db:"images"`
	Title            string             `json:"title" db:"title"`
	Deadline         string             `json:"deadline" db:"deadline"`
	StagesId         int                `json:"stages_id" db:"stages_id"`
	Rating           string             `json:"rating" db:"rating"`
	Description      string             `json:"description" db:"description"`
	LeftToPay        int                `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int                `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string             `json:"status" db:"status"`
	State            string             `json:"state" db:"state"`
	AddTasks         []AddTaskManager   `json:"tasks"`
	AddCheque        []AddChequeManager `json:"cheques"`
}

type AddTaskManager struct {
	Id       int            `json:"id" db:"id"`
	Image    pq.StringArray `json:"images" db:"images"`
	Title    string         `json:"title" db:"title"`
	CardId   int            `json:"card_id" db:"card_id"`
	Complete bool           `json:"complete" db:"complete"`
	Length   int            `json:"length" db:"length"`
}

type CardsAdd struct {
	Id               int                `json:"id" db:"id"`
	Image            pq.StringArray     `json:"images" db:"images"`
	Title            string             `json:"title" db:"title"`
	Deadline         string             `json:"deadline" db:"deadline"`
	StagesId         int                `json:"stages_id" db:"stages_id"`
	Rating           string             `json:"rating" db:"rating"`
	Description      string             `json:"description" db:"description"`
	LeftToPay        int                `json:"left_to_pay" db:"left_to_pay"`
	SpentOnMaterials int                `json:"spent_on_materials" db:"spent_on_materials"`
	Status           string             `json:"status" db:"status"`
	State            string             `json:"state" db:"state"`
	ProjectId        int                `json:"project_id"`
	AddTasks         []AddTaskManager   `json:"tasks"`
	AddCheque        []AddChequeManager `json:"cheques"`
}

type AddChequeManager struct {
	Id        int            `json:"id" db:"id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	File      pq.StringArray `json:"file" db:"file"`
	Name      string         `json:"name" db:"name"`
	Cost      int            `json:"cost" db:"cost"`
	Type      string         `json:"type" db:"type"`
	CardId    int            `json:"card_id" db:"card_id"`
	UserId    int            `json:"user_id" db:"user_id"`
	ProjectId int            `json:"project_id" db:"project_id"`
	Length    int            `json:"length" db:"length"`
	Weight    pq.Int64Array  `json:"weight" db:"weight"`
}

type FileInfo struct {
	ID       int       `json:"id" db:"id"`
	Date     time.Time `json:"created_at" db:"created_at"`
	Name     string    `json:"url" db:"url"`
	Length   int64     `json:"length" db:"length"`
	MimeType string    `json:"mime" db:"mime"`
	Object   string    `json:"object" db:"object"`
	Role     string    `json:"role" db:"role"`
	UserId   int       `json:"user_id" db:"user_id"`
	Tag      string    `json:"tag" db:"tag"`
}

type DeviceToken struct {
	Token string `json:"token" db:"token"`
}

type Sms struct {
	Messages []SMSMessages `json:"messages"`
	Login    string        `json:"login"`
	Password string        `json:"password"`
}

type SMSMessages struct {
	Phone    string `json:"phone"`
	ClientID int    `json:"clientId"`
	Text     string `json:"text"`
	Sender   string `json:"sender"`
}

type DeleteCard struct {
	ProjectID int
	StageId   int
	CardId    int
	UserID    int
}

type DeleteStage struct {
	ProjectID int
	StageId   int
	UserID    int
}
