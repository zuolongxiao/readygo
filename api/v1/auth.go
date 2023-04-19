package api

import (
	"database/sql"
	"time"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/pkg/settings"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// Auth get JWT token
func Auth(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	au := models.Auth{}
	if err := w.Bind(&au); err != nil {
		w.Respond(err, nil)
		return
	}

	if !captchaStore.Verify(au.CaptchaID, au.CaptchaCode, true) {
		w.Respond(errs.ValidationError("invalid captcha code"), nil)
		return
	}

	admin := models.Admin{}
	s := services.New(&admin)
	if err := s.LoadByKey("username", au.Username); err != nil {
		// w.Respond(err, nil)
		w.Respond(errs.ValidationError("incorrect username or password"), nil)
		return
	}

	if !utils.VerifyPassword(admin.Password, au.Password) {
		w.Respond(errs.ValidationError("incorrect username or password"), nil)
		return
	}

	if admin.IsLocked == "Y" {
		w.Respond(errs.LockedError(admin.Username), nil)
		return
	}

	permissions := make([]string, 0, 10)
	if admin.ID == settings.App.SuperAdminID {
		permissions = append(permissions, "*")
	} else {
		if admin.RoleID > 0 {
			ao := []models.AuthorizationPermission{}
			am := models.Authorization{}
			roleID := map[string]interface{}{
				"role_id": admin.RoleID,
			}
			as := services.New(&am)
			if err := as.GetRows(&ao, roleID); err != nil {
				w.Respond(err, nil)
				return
			}
			permissionIDs := make([]uint64, 0, 10)
			for _, v := range ao {
				permissionIDs = append(permissionIDs, v.PermissionID)
			}

			if len(permissionIDs) > 0 {
				po := []models.PermissionView{}
				pm := models.Permission{}
				ps := services.New(&pm)
				if err := ps.GetRows(&po, permissionIDs); err != nil {
					w.Respond(err, nil)
					return
				}
				for _, v := range po {
					if v.IsEnabled == "Y" {
						permissions = append(permissions, v.Name)
					}
				}
			}
		}
	}

	token, err := utils.GenerateToken(admin.Username, permissions)
	if err != nil {
		w.Respond(errs.InternalServerError(err.Error()), nil)
		return
	}

	admin.LastLoginIP = c.ClientIP()
	admin.LastLoginAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	_ = s.Update("LastLoginIP", "LastLoginAt")

	nowTime := time.Now()
	expireTime := nowTime.Add(settings.JWT.Expires)
	data := map[string]string{
		"token":   token,
		"expires": expireTime.Format(time.RFC3339),
	}
	w.Respond(nil, data)
}
