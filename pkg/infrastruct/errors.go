package infrastruct

import "net/http"

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

var (
	ErrorInternalServerError        = NewError("внутренняя ошибка сервера", http.StatusInternalServerError)
	ErrorBadRequest                 = NewError("плохие входные данные запроса", http.StatusBadRequest)
	ErrorPhoneIsExist               = NewError("номер уже зарегистрирован", http.StatusBadRequest)
	ErrorEmailIsExist               = NewError("email уже зарегистрирован", http.StatusBadRequest)
	ErrorPhoneIsIncorrect           = NewError("неверный номер телефона", http.StatusForbidden)
	ErrorEmailIsIncorrect           = NewError("неверный email", http.StatusForbidden)
	ErrorJWTIsBroken                = NewError("jwt испорчен", http.StatusForbidden)
	ErrorPermissionDenied           = NewError("у вас недостаточно прав", http.StatusForbidden)
	ErrorPhoneIsExistManager        = NewError("зарегистрирован менеджером", http.StatusBadRequest)
	ErrorCodeIsIncorrect            = NewError("неверный код", http.StatusForbidden)
	ErrorOldPhoneIsIncorrect        = NewError("неверный старый номер телефона", http.StatusForbidden)
	ErrorPasswordOrEmailIsIncorrect = NewError("Неверный пароль или логин", http.StatusForbidden)
	ErrorPhoneOrEmailIsIncorrect    = NewError("the phone number or email matches the existing one", http.StatusConflict)
	ErrorTeapot                     = NewError("неверная дата этапа", http.StatusTeapot)
	ErrorStageDate                  = NewError("identical stage date in the project", http.StatusBadRequest)

	//ТОЛЬКО ДЛЯ ПРОЕКТА
	ErrorProjectBudgetNegative = NewError("budget cannot be negative", http.StatusBadRequest)
	ErrorProjectTimeOut        = NewError("error start date or end date in project creation", http.StatusBadRequest)
	ErrorProjectHotCoffee      = NewError("wrong project dates", http.StatusBadRequest)
	ErrorProject               = NewError("error creating project", http.StatusBadRequest)
	ErrorProjectStageTimeOut   = NewError("error start date or end date in stage creation", http.StatusBadRequest)
	ErrorProjectStageHotCoffee = NewError("wrong stage date", http.StatusBadRequest)
	ErrorProjectStage          = NewError("error creating stage", http.StatusBadRequest)
	ErrorProjectCardTimeOut    = NewError("error start date or end date in card creation", http.StatusBadRequest)
	ErrorProjectCardHotCoffee  = NewError("wrong card deadline", http.StatusBadRequest)
	ErrorProjectCard           = NewError("error creating card", http.StatusBadRequest)
	ErrorProjectTask           = NewError("error creating task", http.StatusBadRequest)
	ErrorProjectCheques        = NewError("error creating cheques", http.StatusBadRequest)
	ErrorProjectCardDelete     = NewError("error tasks or cheques have not been deleted", http.StatusBadRequest)
	ErrorProjectEmpty          = NewError("error wrong project", http.StatusBadRequest)
	ErrorWrongStage            = NewError("error wrong stage", http.StatusBadRequest)
	ErrorWrongCard             = NewError("error wrong card", http.StatusBadRequest)
	ErrorWrongProject          = NewError("error wrong project", http.StatusBadRequest)
	ErrorProjectStageDelete    = NewError("error cards have not been deleted", http.StatusBadRequest)
)
