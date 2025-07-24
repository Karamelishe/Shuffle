import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpApi from 'i18next-http-backend';

i18n
  .use(HttpApi)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    lng: 'en', // default language
    fallbackLng: 'en',
    debug: false, // Set to true for development
    
    interpolation: {
      escapeValue: false, // React already does escaping
    },

    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },

    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage'],
    },

    resources: {
      en: {
        translation: {
          // Common UI elements
          "common": {
            "save": "Save",
            "cancel": "Cancel",
            "delete": "Delete",
            "edit": "Edit",
            "create": "Create",
            "update": "Update",
            "close": "Close",
            "back": "Back",
            "next": "Next",
            "previous": "Previous",
            "loading": "Loading...",
            "search": "Search",
            "filter": "Filter",
            "name": "Name",
            "description": "Description",
            "tags": "Tags",
            "status": "Status",
            "actions": "Actions",
            "settings": "Settings",
            "help": "Help",
            "documentation": "Documentation",
            "logout": "Logout",
            "login": "Login",
            "register": "Register",
            "password": "Password",
            "email": "Email",
            "username": "Username",
            "organization": "Organization",
            "dashboard": "Dashboard",
            "workflows": "Workflows",
            "apps": "Apps",
            "admin": "Admin",
            "profile": "Profile",
            "notifications": "Notifications",
            "home": "Home"
          },
          
          // Navigation
          "nav": {
            "dashboard": "Dashboard",
            "workflows": "Workflows", 
            "apps": "Apps",
            "docs": "Documentation",
            "admin": "Admin",
            "settings": "Settings",
            "getting_started": "Getting Started",
            "usecases": "Use Cases",
            "search": "Search"
          },

          // Login page
          "login": {
            "title": "Login to Shuffle",
            "email_placeholder": "Email",
            "password_placeholder": "Password",
            "login_button": "Login",
            "forgot_password": "Forgot Password?",
            "no_account": "Don't have an account?",
            "register_link": "Register here",
            "welcome_back": "Welcome back",
            "sign_in": "Sign in to your account"
          },

          // Dashboard
          "dashboard": {
            "title": "Dashboard",
            "welcome": "Welcome to Shuffle",
            "recent_workflows": "Recent Workflows",
            "recent_executions": "Recent Executions",
            "statistics": "Statistics",
            "quick_actions": "Quick Actions",
            "create_workflow": "Create Workflow",
            "explore_apps": "Explore Apps"
          },

          // Workflows
          "workflows": {
            "title": "Workflows",
            "create_new": "Create New Workflow",
            "import": "Import Workflow",
            "export": "Export Workflow",
            "run": "Run Workflow",
            "edit": "Edit Workflow",
            "delete_confirm": "Are you sure you want to delete this workflow?",
            "name_placeholder": "Workflow Name",
            "description_placeholder": "Workflow Description",
            "tags_placeholder": "Tags (comma separated)",
            "search_placeholder": "Search Workflows",
            "no_workflows": "No workflows found",
            "execution_history": "Execution History",
            "triggers": "Triggers",
            "actions": "Actions"
          },

          // Apps
          "apps": {
            "title": "Apps",
            "explore": "Explore Apps",
            "my_apps": "My Apps",
            "create_app": "Create App",
            "install": "Install",
            "uninstall": "Uninstall",
            "configure": "Configure",
            "documentation": "Documentation",
            "search_placeholder": "Search Apps",
            "categories": "Categories",
            "featured": "Featured",
            "popular": "Popular"
          },

          // Forms and validation
          "validation": {
            "required": "This field is required",
            "invalid_email": "Please enter a valid email address", 
            "password_too_short": "Password must be at least 8 characters",
            "passwords_dont_match": "Passwords don't match"
          },

          // Messages and notifications
          "messages": {
            "success": "Operation completed successfully",
            "error": "An error occurred",
            "saved": "Saved successfully",
            "deleted": "Deleted successfully",
            "updated": "Updated successfully",
            "created": "Created successfully",
            "loading": "Loading...",
            "no_results": "No results found",
            "confirm_delete": "Are you sure you want to delete this item?"
          }
        }
      },
      ru: {
        translation: {
          // Common UI elements
          "common": {
            "save": "Сохранить",
            "cancel": "Отмена",
            "delete": "Удалить",
            "edit": "Редактировать",
            "create": "Создать",
            "update": "Обновить",
            "close": "Закрыть",
            "back": "Назад",
            "next": "Далее",
            "previous": "Предыдущий",
            "loading": "Загрузка...",
            "search": "Поиск",
            "filter": "Фильтр",
            "name": "Имя",
            "description": "Описание",
            "tags": "Теги",
            "status": "Статус",
            "actions": "Действия",
            "settings": "Настройки",
            "help": "Помощь",
            "documentation": "Документация",
            "logout": "Выйти",
            "login": "Войти",
            "register": "Регистрация",
            "password": "Пароль",
            "email": "Эл. почта",
            "username": "Имя пользователя",
            "organization": "Организация",
            "dashboard": "Панель управления",
            "workflows": "Рабочие процессы",
            "apps": "Приложения",
            "admin": "Администратор",
            "profile": "Профиль",
            "notifications": "Уведомления",
            "home": "Главная"
          },
          
          // Navigation
          "nav": {
            "dashboard": "Панель управления",
            "workflows": "Рабочие процессы",
            "apps": "Приложения", 
            "docs": "Документация",
            "admin": "Администратор",
            "settings": "Настройки",
            "getting_started": "Начало работы",
            "usecases": "Сценарии использования",
            "search": "Поиск"
          },

          // Login page
          "login": {
            "title": "Вход в Shuffle",
            "email_placeholder": "Эл. почта",
            "password_placeholder": "Пароль",
            "login_button": "Войти",
            "forgot_password": "Забыли пароль?",
            "no_account": "Нет аккаунта?",
            "register_link": "Зарегистрироваться здесь",
            "welcome_back": "Добро пожаловать",
            "sign_in": "Войдите в свой аккаунт"
          },

          // Dashboard
          "dashboard": {
            "title": "Панель управления",
            "welcome": "Добро пожаловать в Shuffle",
            "recent_workflows": "Недавние рабочие процессы",
            "recent_executions": "Недавние выполнения",
            "statistics": "Статистика",
            "quick_actions": "Быстрые действия",
            "create_workflow": "Создать рабочий процесс",
            "explore_apps": "Исследовать приложения"
          },

          // Workflows
          "workflows": {
            "title": "Рабочие процессы",
            "create_new": "Создать новый рабочий процесс",
            "import": "Импортировать рабочий процесс",
            "export": "Экспортировать рабочий процесс",
            "run": "Запустить рабочий процесс",
            "edit": "Редактировать рабочий процесс",
            "delete_confirm": "Вы уверены, что хотите удалить этот рабочий процесс?",
            "name_placeholder": "Название рабочего процесса",
            "description_placeholder": "Описание рабочего процесса",
            "tags_placeholder": "Теги (через запятую)",
            "search_placeholder": "Поиск рабочих процессов",
            "no_workflows": "Рабочие процессы не найдены",
            "execution_history": "История выполнения",
            "triggers": "Триггеры",
            "actions": "Действия"
          },

          // Apps
          "apps": {
            "title": "Приложения",
            "explore": "Исследовать приложения",
            "my_apps": "Мои приложения",
            "create_app": "Создать приложение",
            "install": "Установить",
            "uninstall": "Удалить",
            "configure": "Настроить",
            "documentation": "Документация",
            "search_placeholder": "Поиск приложений",
            "categories": "Категории",
            "featured": "Рекомендуемые",
            "popular": "Популярные"
          },

          // Forms and validation
          "validation": {
            "required": "Это поле обязательно для заполнения",
            "invalid_email": "Пожалуйста, введите действительный адрес электронной почты",
            "password_too_short": "Пароль должен содержать не менее 8 символов",
            "passwords_dont_match": "Пароли не совпадают"
          },

          // Messages and notifications
          "messages": {
            "success": "Операция выполнена успешно",
            "error": "Произошла ошибка",
            "saved": "Сохранено успешно",
            "deleted": "Удалено успешно",
            "updated": "Обновлено успешно",
            "created": "Создано успешно",
            "loading": "Загрузка...",
            "no_results": "Результаты не найдены",
            "confirm_delete": "Вы уверены, что хотите удалить этот элемент?"
          },

          // License management
          "license": {
            "title": "Управление лицензиями",
            "current_license": "Текущая лицензия",
            "no_license": "Лицензия не активирована",
            "no_license_description": "Для доступа к расширенному функционалу необходимо активировать лицензию.",
            "activate_license": "Активировать лицензию",
            "change_license": "Изменить лицензию",
            "license_key": "Лицензионный ключ",
            "key_required": "Введите лицензионный ключ",
            "activation_success": "Лицензия успешно активирована",
            "activation_failed": "Не удалось активировать лицензию",
            "activation_error": "Ошибка при активации лицензии",
            "activation_description": "Введите лицензионный ключ для активации расширенного функционала.",
            "activate": "Активировать",
            "expires_at": "Действует до",
            "days_remaining": "Осталось дней: {{days}}",
            "expired": "Истекла",
            "unlimited": "Без ограничений",
            "limits": "Лимиты",
            "max_users": "Максимум пользователей",
            "max_workflows": "Максимум workflow'ов",
            "max_executions": "Максимум выполнений",
            "features": "Функции",
            "type_basic": "Базовая",
            "type_professional": "Профессиональная",
            "type_enterprise": "Корпоративная",
            "status_active": "Активна",
            "status_expired": "Истекла",
            "status_revoked": "Отозвана",
            "feature_workflows_5": "До 5 workflow'ов",
            "feature_workflows_50": "До 50 workflow'ов",
            "feature_workflows_unlimited": "Неограниченные workflow'ы",
            "feature_users_3": "До 3 пользователей",
            "feature_users_10": "До 10 пользователей",
            "feature_users_unlimited": "Неограниченные пользователи",
            "feature_executions_1000": "До 1000 выполнений",
            "feature_executions_10000": "До 10000 выполнений",
            "feature_executions_unlimited": "Неограниченные выполнения",
            "feature_apps_basic": "Базовые приложения",
            "feature_apps_all": "Все приложения",
            "feature_support_community": "Поддержка сообщества",
            "feature_support_email": "Email поддержка",
            "feature_support_priority": "Приоритетная поддержка",
            "feature_integrations_advanced": "Расширенные интеграции",
            "feature_integrations_all": "Все интеграции",
            "feature_reporting_basic": "Базовая отчетность",
            "feature_reporting_advanced": "Расширенная отчетность",
            "feature_sso_enabled": "Single Sign-On",
            "feature_audit_enabled": "Аудит действий",
            "feature_custom_branding_enabled": "Кастомный брендинг",
            "license_error": "Ошибка лицензии",
            "license_warning": "Предупреждение о лицензии",
            "expired_message": "Ваша лицензия истекла. Некоторые функции могут быть недоступны.",
            "expiring_soon": "Лицензия скоро истечет",
            "expiring_soon_message": "Ваша лицензия истечет через {{days}} дней.",
            "renew_license": "Обновить лицензию",
            "upgrade_license": "Улучшить лицензию",
            "usage_warning": "Предупреждение об использовании",
            "manage_license": "Управление лицензией",
            "default_warning_message": "Для использования этой функции требуется активная лицензия.",
            "current_license_info": "Информация о текущей лицензии",
            "type": "Тип"
          }
        }
      }
    }
  });

export default i18n;