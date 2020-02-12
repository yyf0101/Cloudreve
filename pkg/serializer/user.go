package serializer

import (
	"fmt"
	"github.com/HFO4/cloudreve/models"
	"github.com/HFO4/cloudreve/pkg/hashid"
)

// CheckLogin 检查登录
func CheckLogin() Response {
	return Response{
		Code: CodeCheckLogin,
		Msg:  "未登录",
	}
}

// User 用户序列化器
type User struct {
	ID             uint   `json:"id"`
	Email          string `json:"user_name"`
	Nickname       string `json:"nickname"`
	Status         int    `json:"status"`
	Avatar         string `json:"avatar"`
	CreatedAt      int64  `json:"created_at"`
	PreferredTheme string `json:"preferred_theme"`
	Score          int    `json:"score"`
	Policy         policy `json:"policy"`
	Group          group  `json:"group"`
	Tags           []tag  `json:"tags"`
}

type policy struct {
	SaveType       string   `json:"saveType"`
	MaxSize        string   `json:"maxSize"`
	AllowedType    []string `json:"allowedType"`
	UploadURL      string   `json:"upUrl"`
	AllowGetSource bool     `json:"allowSource"`
}

type group struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	AllowShare           bool   `json:"allowShare"`
	AllowRemoteDownload  bool   `json:"allowRemoteDownload"`
	AllowArchiveDownload bool   `json:"allowArchiveDownload"`
	ShareFreeEnabled     bool   `json:"shareFree"`
	ShareDownload        bool   `json:"shareDownload"`
	CompressEnabled      bool   `json:"compress"`
}

type tag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Color      string `json:"color"`
	Type       int    `json:"type"`
	Expression string `json:"expression"`
}

type storage struct {
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	tags, _ := model.GetTagsByUID(user.ID)
	return User{
		ID:             user.ID,
		Email:          user.Email,
		Nickname:       user.Nick,
		Status:         user.Status,
		Avatar:         user.Avatar,
		CreatedAt:      user.CreatedAt.Unix(),
		PreferredTheme: user.OptionsSerialized.PreferredTheme,
		Score:          user.Score,
		Policy: policy{
			SaveType:       user.Policy.Type,
			MaxSize:        fmt.Sprintf("%.2fmb", float64(user.Policy.MaxSize)/(1024*1024)),
			AllowedType:    user.Policy.OptionsSerialized.FileType,
			UploadURL:      user.Policy.GetUploadURL(),
			AllowGetSource: user.Policy.IsOriginLinkEnable,
		},
		Group: group{
			ID:                   user.GroupID,
			Name:                 user.Group.Name,
			AllowShare:           user.Group.ShareEnabled,
			AllowRemoteDownload:  user.Group.OptionsSerialized.Aria2,
			AllowArchiveDownload: user.Group.OptionsSerialized.ArchiveDownload,
			ShareFreeEnabled:     user.Group.OptionsSerialized.ShareFree,
			ShareDownload:        user.Group.OptionsSerialized.ShareDownload,
			CompressEnabled:      user.Group.OptionsSerialized.ArchiveTask,
		},
		Tags: buildTagRes(tags),
	}
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user model.User) Response {
	return Response{
		Data: BuildUser(user),
	}
}

// BuildUserStorageResponse 序列化用户存储概况响应
func BuildUserStorageResponse(user model.User) Response {
	total := user.Group.MaxStorage + user.GetAvailablePackSize()
	storageResp := storage{
		Used:  user.Storage,
		Free:  total - user.Storage,
		Total: total,
	}

	if total < user.Storage {
		storageResp.Free = 0
	}

	return Response{
		Data: storageResp,
	}
}

// buildTagRes 构建标签列表
func buildTagRes(tags []model.Tag) []tag {
	res := make([]tag, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		newTag := tag{
			ID:    hashid.HashID(tags[i].ID, hashid.TagID),
			Name:  tags[i].Name,
			Icon:  tags[i].Icon,
			Color: tags[i].Color,
			Type:  tags[i].Type,
		}
		if newTag.Type != 0 {
			newTag.Expression = tags[i].Expression

		}
		res = append(res, newTag)
	}

	return res
}
