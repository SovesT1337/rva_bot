package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"x.localhost/rvabot/config"
	"x.localhost/rvabot/internal/database"
	"x.localhost/rvabot/internal/errors"
	"x.localhost/rvabot/internal/handler"
	"x.localhost/rvabot/internal/health"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/metrics"
	"x.localhost/rvabot/internal/ratelimit"
	"x.localhost/rvabot/internal/recovery"
	"x.localhost/rvabot/internal/shutdown"
	"x.localhost/rvabot/internal/state"
	"x.localhost/rvabot/internal/telegram"

	"github.com/joho/godotenv"
)

// BotService представляет основной сервис бота
type BotService struct {
	config          *config.Config
	database        *database.Database
	repo            database.ContentRepositoryInterface
	rateLimiter     *ratelimit.UserRateLimiter
	stateManager    *state.Manager
	healthManager   *health.Manager
	shutdownManager *shutdown.Manager
	server          *http.Server
}

// NewBotService создает новый экземпляр сервиса бота
func NewBotService(cfg *config.Config) *BotService {
	return &BotService{
		config: cfg,
	}
}

// Initialize инициализирует сервис
func (bs *BotService) Initialize() error {
	// Инициализируем базу данных
	var err error
	bs.database, err = database.NewDatabase(bs.config.Database.Path)
	if err != nil {
		appErr := errors.WrapError(err, errors.ErrorTypeDatabase, "Ошибка инициализации БД")
		logger.BotError("Ошибка инициализации БД: %v", appErr)
		return appErr
	}

	logger.DatabaseInfo("База данных инициализирована успешно")

	// Создаем репозиторий
	bs.repo = database.NewContentRepository(bs.database)

	// Инициализируем rate limiter
	bs.rateLimiter = ratelimit.NewUserRateLimiter(ratelimit.DefaultConfig())

	// Инициализируем state manager
	bs.stateManager = state.NewManager(30*time.Minute, 5*time.Minute)

	// Запускаем метрики
	metrics.StartMetricsLogger(5 * time.Minute) // Логируем метрики каждые 5 минут

	// Инициализируем health checks
	bs.healthManager = health.NewManager(5 * time.Second)
	bs.setupHealthChecks()

	// Инициализируем HTTP сервер для health checks
	bs.setupHTTPServer()

	// Инициализируем shutdown manager
	bs.shutdownManager = shutdown.NewManager(30 * time.Second)
	bs.setupShutdownHandlers()

	return nil
}

// setupHealthChecks настраивает health checks
func (bs *BotService) setupHealthChecks() {
	// Database health check
	dbCheck := health.NewDatabaseCheck(func(ctx context.Context) error {
		// Простая проверка подключения к БД
		_, err := bs.repo.GetTrainers()
		return err
	})
	bs.healthManager.RegisterCheck(dbCheck)

	// Telegram API health check
	telegramCheck := health.NewTelegramCheck(func(ctx context.Context) error {
		// Проверяем доступность Telegram API
		_, err := telegram.GetUpdates(bs.config.GetBotURL(), 0)
		return err
	})
	bs.healthManager.RegisterCheck(telegramCheck)

	// Memory health check
	memoryCheck := health.NewMemoryCheck(512) // 512MB лимит
	bs.healthManager.RegisterCheck(memoryCheck)
}

// setupHTTPServer настраивает HTTP сервер для health checks
func (bs *BotService) setupHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", bs.healthManager.HTTPHandler())
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	bs.server = &http.Server{
		Addr:         ":" + bs.config.Server.Port,
		Handler:      mux,
		ReadTimeout:  bs.config.Server.ReadTimeout,
		WriteTimeout: bs.config.Server.WriteTimeout,
	}
}

// setupShutdownHandlers настраивает обработчики shutdown
func (bs *BotService) setupShutdownHandlers() {
	// HTTP сервер
	bs.shutdownManager.RegisterHandler(&httpShutdownHandler{server: bs.server})

	// База данных
	bs.shutdownManager.RegisterHandler(&databaseShutdownHandler{database: bs.database})

	// State manager
	bs.shutdownManager.RegisterHandler(&stateShutdownHandler{stateManager: bs.stateManager})

	// Rate limiter (если нужен cleanup)
	bs.shutdownManager.RegisterHandler(&rateLimiterShutdownHandler{rateLimiter: bs.rateLimiter})
}

// Start запускает сервис
func (bs *BotService) Start() error {
	botUrl := bs.config.GetBotURL()
	logger.BotInfo("Запуск бота...")
	logger.BotInfo("URL бота: %s", botUrl)

	// Запускаем HTTP сервер в горутине с recovery
	recovery.RecoverGoroutine(context.Background(), "http_server", func() {
		logger.BotInfo("Запуск HTTP сервера на :%s", bs.config.Server.Port)
		if err := bs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.BotError("Ошибка HTTP сервера: %v", err)
		}
	})

	// Запускаем основной цикл бота в горутине с recovery
	recovery.RecoverGoroutine(context.Background(), "bot_loop", func() {
		logger.BotInfo("Запуск основного цикла бота...")
		handler.BotLoopWithComponents(botUrl, bs.repo, bs.rateLimiter, bs.stateManager, nil)
	})

	// Запускаем shutdown manager
	bs.shutdownManager.Start()
	logger.BotInfo("Бот запущен и готов к работе. Нажмите Ctrl+C для завершения...")

	// Ждем сигнал завершения
	bs.shutdownManager.Wait()

	return nil
}

// HTTP shutdown handler
type httpShutdownHandler struct {
	server *http.Server
}

func (h *httpShutdownHandler) Name() string {
	return "http_server"
}

func (h *httpShutdownHandler) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

// Database shutdown handler
type databaseShutdownHandler struct {
	database *database.Database
}

func (h *databaseShutdownHandler) Name() string {
	return "database"
}

func (h *databaseShutdownHandler) Shutdown(ctx context.Context) error {
	return h.database.Close()
}

// State manager shutdown handler
type stateShutdownHandler struct {
	stateManager *state.Manager
}

func (h *stateShutdownHandler) Name() string {
	return "state_manager"
}

func (h *stateShutdownHandler) Shutdown(ctx context.Context) error {
	h.stateManager.Shutdown()
	return nil
}

// Rate limiter shutdown handler
type rateLimiterShutdownHandler struct {
	rateLimiter *ratelimit.UserRateLimiter
}

func (h *rateLimiterShutdownHandler) Name() string {
	return "rate_limiter"
}

func (h *rateLimiterShutdownHandler) Shutdown(ctx context.Context) error {
	// Rate limiter не требует специального shutdown
	logger.BotInfo("Rate limiter shutdown completed")
	return nil
}

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		logger.Warn("MAIN", "Не удалось загрузить .env файл: %v", err)
	}

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Валидируем конфигурацию
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфигурации: %v", err)
	}

	// Устанавливаем уровень логирования
	logger.SetLevel(logger.INFO)

	// Создаем и инициализируем сервис
	botService := NewBotService(cfg)
	if err := botService.Initialize(); err != nil {
		log.Fatalf("Ошибка инициализации сервиса: %v", err)
	}

	// Запускаем сервис
	if err := botService.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервиса: %v", err)
	}
}
