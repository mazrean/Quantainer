package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/pkg/common"
	"github.com/mazrean/Quantainer/service"
)

const (
	basePath            = "https://q.trap.jp/api/v3"
	limit               = 200
	embURLRegexFragment = `/files/([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})`
)

var (
	imageRegex  = regexp.MustCompile(`^image/*`)
	embURLRegex = regexp.MustCompile(strings.ReplaceAll("https://q.trap.jp", ".", `\.`) + embURLRegexFragment)
)

type user struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type fileMeta struct {
	ID            string `json:"id,omitempty"`
	UserID        string `json:"uploaderId,omitempty"`
	Name          string `json:"name,omitempty"`
	Mime          string `json:"mime,omitempty"`
	Md5           string `json:"md5,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	timeCreatedAt time.Time
}

type Bot struct {
	accessToken       string
	verificationToken string
	defaultChannels   []string
	fileService       service.File
	resourceService   service.Resource
}

func NewBot(
	accessToken common.AccessToken,
	verificationToken common.VerificationToken,
	defaultChannels common.DefaultChannels,
	updatedAt common.UpdatedAt,
	fileService service.File,
	resourceService service.Resource,
) (*Bot, error) {
	bot := &Bot{
		accessToken:       string(accessToken),
		verificationToken: string(verificationToken),
		defaultChannels:   []string(defaultChannels),
		fileService:       fileService,
		resourceService:   resourceService,
	}

	go func(){
		err := bot.setupBot(context.Background(), time.Time(updatedAt))
		if err != nil {
			log.Printf("failed to setup bot: %v\n", err)
		}
	}()

	return bot, nil
}

func (b *Bot) setupBot(ctx context.Context, updatedAt time.Time) error {
	log.Println("info: init start")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users?include-suspended=true", basePath), nil)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+b.accessToken)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do HTTP request: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to get messages:(Status:%d %s)", res.StatusCode, res.Status)
	}

	users := []user{}
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		return fmt.Errorf("failed to decode HTTP response: %w", err)
	}

	userMap := map[string]user{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	strUpdatedAt := updatedAt.Format("2006-01-02T15:04:05.999999Z")

	for _, defaultChannel := range b.defaultChannels {
		for {
			time.Sleep(time.Second)

			fmt.Println("init start", strUpdatedAt)
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/files?channelId=%s&limit=%d&order=asc&since=%s", basePath, defaultChannel, limit, strUpdatedAt), nil)
			if err != nil {
				return fmt.Errorf("failed to make HTTP request: %w", err)
			}
			req.Header.Add("Authorization", "Bearer "+b.accessToken)

			res, err := httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to do HTTP request: %w", err)
			}
			if res.StatusCode != 200 {
				return fmt.Errorf("failed to get file metas:(Status:%d %s)", res.StatusCode, res.Status)
			}

			fileMetas := []fileMeta{}
			err = json.NewDecoder(res.Body).Decode(&fileMetas)
			if err != nil {
				return fmt.Errorf("failed to decode HTTP response: %w", err)
			}

			for _, meta := range fileMetas {
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/files/%s", basePath, meta.ID), nil)
				if err != nil {
					return fmt.Errorf("failed to make HTTP request: %w", err)
				}
				req.Header.Add("Authorization", "Bearer "+b.accessToken)

				res, err := httpClient.Do(req)
				if err != nil {
					return fmt.Errorf("failed to do HTTP request: %w", err)
				}
				if res.StatusCode != 200 {
					return fmt.Errorf("failed to get file metas:(Status:%d %s)", res.StatusCode, res.Status)
				}

				user := service.NewUserInfo(
					values.NewTrapMemberID(uuid.MustParse(userMap[meta.UserID].ID)),
					values.NewTrapMemberName(userMap[meta.UserID].Name),
					values.TrapMemberStatusActive,
				)

				file, err := b.fileService.UploadBotFile(ctx, user, res.Body)
				if err != nil {
					return fmt.Errorf("failed to upload file: %w", err)
				}

				var resourceType values.ResourceType
				if file.File.GetType() == values.FileTypeOther {
					resourceType = values.ResourceTypeOther
				} else {
					resourceType = values.ResourceTypeImage
				}

				createdAt, err := time.Parse("2006-01-02T15:04:05.999999Z", meta.CreatedAt)
				if err != nil {
					return fmt.Errorf("failed to parse createdAt: %w", err)
				}

				_, err = b.resourceService.CreateBotResource(ctx,
					user,
					file.File.GetID(),
					values.NewResourceName("traQ File"),
					resourceType,
					values.NewResourceComment(""),
					createdAt,
				)
				if err != nil {
					return fmt.Errorf("failed to create resource: %w", err)
				}
			}

			if len(fileMetas) == limit {
				updatedAt, err = time.Parse(time.RFC3339, fileMetas[limit-1].CreatedAt)
				if err != nil {
					return fmt.Errorf("failed to parse time: %w", err)
				}
				strUpdatedAt = updatedAt.Format("2006-01-02T15:04:05.999999Z")
			} else {
				break
			}
		}
	}

	fmt.Println("init end")
	return nil
}
