package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zuolongxiao/readygo/models"
	"github.com/zuolongxiao/readygo/pkg/utils"
	"github.com/zuolongxiao/readygo/services"
)

// ListTags list tags
func ListTags(c *gin.Context) {
	w := utils.NewContextWrapper(c)
	s := services.New(&models.Tag{})

	var tags []models.TagView
	if err := s.Find(&tags, c); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"list":   tags,
		"offset": s.GetOffset(),
	}

	w.Respond(nil, data)
}

// GetTag get tag
func GetTag(c *gin.Context) {
	w := utils.NewContextWrapper(c)
	s := services.New(&models.Tag{})

	var tag models.TagView
	if err := s.GetByID(&tag, c.Param("id")); err != nil {
		w.Respond(err, nil)
		return
	}

	w.Respond(nil, tag)
}

// CreateTag create tag
func CreateTag(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	binding := models.TagCreate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Tag{}
	s := services.New(&m)
	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m.CreatedBy = w.GetUsername()
	if err := s.Create(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"created_at": m.CreatedAt.Time,
	}

	w.Respond(nil, data)
}

// UpdateTag update tag
func UpdateTag(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	binding := models.TagUpdate{}
	if err := w.Bind(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m := models.Tag{}
	s := services.New(&m)
	if err := s.LoadByID(c.Param("id")); err != nil {
		w.Respond(err, nil)
		return
	}

	if err := s.Fill(&binding); err != nil {
		w.Respond(err, nil)
		return
	}

	m.UpdatedBy = w.GetUsername()
	if err := s.Save(); err != nil {
		w.Respond(err, nil)
		return
	}

	data := map[string]interface{}{
		"id":         m.ID,
		"updated_at": m.UpdatedAt.Time,
	}

	w.Respond(nil, data)
}

// DeleteTag delete tag
func DeleteTag(c *gin.Context) {
	w := utils.NewContextWrapper(c)

	s := services.New(&models.Tag{})
	if err := s.LoadByID(c.Param("id")); err != nil {
		w.Respond(err, nil)
		return
	}

	if err := s.Delete(); err != nil {
		w.Respond(err, nil)
		return
	}

	w.Respond(nil, nil)
}
