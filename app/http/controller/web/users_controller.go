package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/request/web/users"
	"goskeleton/app/model"
	"goskeleton/app/service/users/service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
	"strconv"
	"time"
)

type Users struct {
}

// Logout 退出登录
func (u *Users) Logout(ctx *gin.Context) {
	token, ok := cur_userinfo.GetCurrentToken(ctx)
	if ok {
		if success := service.UserServiceFactory().Logout(token); success {
			response.Success(ctx, consts.CurdStatusOkMsg, nil)
			return
		}
	}
	response.Fail(ctx, consts.NotAuthorize, consts.UserIdNotExist, nil)
}

// CurrentUser 获取当前用户
func (u *Users) CurrentUser(context *gin.Context) {
	if userId, exist := cur_userinfo.GetCurrentUserId(context); exist {
		currentUser := service.UserServiceFactory().FindById(userId)
		response.Success(context, consts.CurdStatusOkMsg, currentUser)
	} else {
		response.Fail(context, consts.NotAuthorize, consts.UserIdNotExist, nil)
	}
}

// UserDetail 1.用户详情
func (u *Users) UserDetail(context *gin.Context) {
	userId := context.Param("userId")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		response.Fail(context, consts.InternalServerErrorCode, consts.InternalServerErrorMsg, "")
	}
	user := service.UserServiceFactory().FindById(id)
	response.Success(context, consts.CurdStatusOkMsg, user)
}

// Register 1.用户注册
func (u *Users) Register(context *gin.Context) {
	var r users.Register
	if err := context.ShouldBindJSON(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}
	userIp := context.ClientIP()
	if service.UserServiceFactory().Register(r.Username, r.Password, userIp) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdRegisterFailCode, consts.CurdRegisterFailMsg, "")
	}
}

// Login 2.用户登录
func (u *Users) Login(context *gin.Context) {
	var r users.Login
	if err := context.ShouldBindJSON(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}
	userModelFact := model.UserModelFactory("")
	userModel, err := userModelFact.Login(r.Username, r.Password)
	if err != nil {
		response.Fail(context, 200, err.Error(), nil)
		return
	}

	if userModel != nil {
		token, err := service.UserServiceFactory().LoginSuccess(userModel)
		if err != nil {
			response.Fail(context, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, "")
		}
		data := gin.H{
			"userId":    userModel.Id,
			"username":  r.Username,
			"realName":  userModel.RealName,
			"phone":     r.Phone,
			"roles":     nil,
			"token":     token,
			"updatedAt": time.Now().Format(variable.DateFormat),
		}
		response.Success(context, consts.CurdStatusOkMsg, data)
		go userModel.UpdateUserLoginInfo(context.ClientIP(), userModel.Id)
		return
	}
	response.Fail(context, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, "")
}

// Show 3.用户查询（show）
func (u *Users) Show(context *gin.Context) {

	var userParam users.UserParam
	if err := context.ShouldBindJSON(&userParam); err != nil {
		response.ValidatorError(context, err)
		return
	}

	userName := userParam.Username
	page := userParam.PageNum
	limit := userParam.Limit

	limitStart := (page - 1) * limit
	counts, showlist := model.UserModelFactory("").Show(userName, limitStart, limit)
	if counts > 0 && showlist != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"counts": counts, "list": showlist})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// Store 4.用户新增(store)
func (u *Users) Store(context *gin.Context) {
	var r users.Register
	if err := context.ShouldBindJSON(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}
	if service.UserServiceFactory().Store(r.Username, r.Password, r.RealName, r.Phone) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "")
	}
}

// Update 5.用户更新(update)
func (u *Users) Update(context *gin.Context) {
	var r users.Register
	if err := context.ShouldBindJSON(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}
	userId := r.UserId
	userName := r.Username
	pass := r.Password
	realName := r.RealName
	phone := r.Phone
	// userIp := context.ClientIP()

	// 检查正在修改的用户名是否被其他人使用
	if model.UserModelFactory("").UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+", "+userName+" 已经被其他人使用", "")
		return
	}

	if service.UserServiceFactory().Update(int(userId), userName, pass, realName, phone) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}

}

// Destroy 6.删除记录
func (u *Users) Destroy(context *gin.Context) {
	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
		return
	}
	if service.UserServiceFactory().Destroy(id) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}

// EditPersonalInfo 编辑自己的信息
/*func (u *Users) EditPersonalInfo(context *gin.Context) {
    // 获取当前请求用户id
    userId, exists := cur_userinfo.GetCurrentUserId(context)
    if !exists {
        response.Fail(context, consts.CurdTokenFailCode, consts.CurdTokenFailMsg, "")
    } else {

        userName := context.GetString(consts.ValidatorPrefix + "user_name")
        usersModel := users.CreateUserFactory("")

        // 检查正在修改的用户名是否被其他站使用
        if usersModel.UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
            response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+",该用户名: "+userName+" 已经被其他人占用", "")
            return
        }
        // 这里使用token解析的id更新表单参数里面的id，加固安全
        context.Set(consts.ValidatorPrefix+"id", float64(userId))

        if usersModel.UpdateData(context) {
            response.Success(context, consts.CurdStatusOkMsg, "")
        } else {
            response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
        }
    }
}*/
