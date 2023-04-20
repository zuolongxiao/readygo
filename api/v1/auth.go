package api

import (
	"database/sql"
	"time"

	"readygo/models"
	"readygo/pkg/errs"
	"readygo/pkg/settings"
	"readygo/pkg/store"
	"readygo/services"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// Auth get JWT token
func Auth(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	au := models.Auth{}
	if err := cw.Bind(&au); err != nil {
		cw.Respond(err, nil)
		return
	}

	if settings.Captcha.Enabled {
		if !store.CaptchaStore.Verify(au.CaptchaID, au.CaptchaCode, true) {
			cw.Respond(errs.ValidationError("invalid captcha code"), nil)
			return
		}
	}

	admin := models.Admin{}
	s := services.New(&admin)
	if err := s.LoadByKey("username", au.Username); err != nil {
		// w.Respond(err, nil)
		cw.Respond(errs.ValidationError("incorrect username or password"), nil)
		return
	}

	if !utils.VerifyPassword(admin.Password, au.Password) {
		cw.Respond(errs.ValidationError("incorrect username or password"), nil)
		return
	}

	if admin.IsLocked == "Y" {
		cw.Respond(errs.LockedError(admin.Username), nil)
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
				cw.Respond(err, nil)
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
					cw.Respond(err, nil)
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
		cw.Respond(errs.InternalServerError(err.Error()), nil)
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
	cw.Respond(nil, data)
}
