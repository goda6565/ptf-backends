package ui

import (
	"context"
	"encoding/json"
	"time"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"github.com/goda6565/ptf-backends/applications/auth/infrastructure/repositoryimpl"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/logger"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/utils"
	"github.com/goda6565/ptf-backends/applications/auth/service"
	"github.com/goda6565/ptf-backends/applications/auth/ui/gen"
	"github.com/goda6565/ptf-backends/applications/auth/ui/handler"
	"github.com/goda6565/ptf-backends/applications/auth/ui/middleware"
)

// swagger設定
func setUpSwagger(router *gin.Engine) (*openapi3.T, error) {
	// OpenAPI (Swagger) 定義を取得
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}

	// 環境変数 ENV を取得（デフォルトは "development"）
	env := utils.GetEnvDefault("ENV", "development")
	if env == "development" {
		// Swagger 定義を JSON に変換
		swaggerJson, _ := json.Marshal(swagger)
		// swag 用の Spec オブジェクトを作成
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",           // インスタンス名
			SwaggerTemplate:  string(swaggerJson), // Swagger のテンプレート（JSON 文字列）
		}
		// swag ライブラリに登録
		swag.Register(SwaggerInfo.InfoInstanceName, SwaggerInfo)
		// /swagger/*any で Swagger UI を提供
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	// Swagger 定義を返す
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	// Gin Engine を作成
	router := gin.New()

	// CORS 設定
	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	// Swagger の設定（開発環境の場合は Swagger UI を有効化）
	swagger, err := setUpSwagger(router)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// ミドルウェアの設定
	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// health check
	router.GET("/health", handler.Health)

	
	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(10 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			// OapiRequestValidatorWithOptions を利用して、認証関数付きのバリデーションミドルウェアを作成
			v1.Use(ginMiddleware.OapiRequestValidatorWithOptions(swagger, &ginMiddleware.Options{
				Options: openapi3filter.Options{
					AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
						bearer := input.RequestValidationInput.Request.Header.Get("Authorization")
						if bearer == "" {
							logger.Error("Authorization header is required")
							return fmt.Errorf("Authorization header is required")
						}
						token := bearer[len("Bearer "):]
						claims, err := utils.ValidateToken(token)
						if err != nil {
							logger.Error(err.Error())
							return err
						}
						ctx := context.WithValue(c, "extraInfo", claims.Email)
						input.RequestValidationInput.Request = input.RequestValidationInput.Request.WithContext(ctx)
						return nil
					},
				},
			}))
			userRepositoryImpl := repositoryimpl.NewUserRepository(db)
			userService := service.NewUserService(userRepositoryImpl)
			userHandler := handler.NewUserHandler(userService)

			api.RegisterHandlers(v1, userHandler)
		}
	}
	// ルーターを返す
	return router, nil
}
