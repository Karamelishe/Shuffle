# Internationalization Setup for Shuffle Frontend

This document outlines the internationalization (i18n) setup implemented for the Shuffle frontend application, including Russian translation support.

## Overview

The frontend has been configured with React i18next to support multiple languages, with initial support for English (default) and Russian.

## Dependencies Added

- `react-i18next@^12.0.0` - React integration for i18next
- Existing i18next packages were already present:
  - `i18next@^22.3.0`
  - `i18next-browser-languagedetector@^7.0.1`
  - `i18next-chained-backend@^4.2.0`
  - `i18next-http-backend@^2.1.1`
  - `i18next-localstorage-backend@^4.1.0`
  - `i18next-xhr-backend@^3.2.2`

## Files Created/Modified

### 1. `frontend/src/i18n.js` (New)
- Main i18n configuration file
- Sets up language detection and storage
- Contains inline translations for English and Russian
- Configures fallback language as English

### 2. `frontend/src/components/LanguageSwitcher.jsx` (New)
- React component for switching between languages
- Displays language selector with flags (🇺🇸 for English, 🇷🇺 for Russian)
- Saves language preference to localStorage
- Styled to match the existing theme

### 3. `frontend/src/index.js` (Modified)
- Added import for `./i18n` to initialize i18n on app startup

### 4. `frontend/package.json` (Modified)
- Added `react-i18next@^12.0.0` dependency

### 5. Component Updates
The following components have been updated with translation support:

#### `frontend/src/components/NewHeader.jsx`
- Added LanguageSwitcher component to both desktop and mobile headers
- Added for both logged-in and non-logged-in states

#### `frontend/src/views/LoginPage.jsx`
- Added `useTranslation` hook
- Translated form labels, placeholders, and button text
- Key translations include:
  - Email/Username field labels
  - Password field labels
  - Login button text
  - Form titles ("Welcome Back", "Sign in")

#### `frontend/src/views/Workflows.jsx`
- Added `useTranslation` hook
- Translated placeholders for:
  - Workflow name field
  - Description field
  - Tags field
  - Search/filter field

## Translation Keys Structure

The translations are organized into logical groups:

```javascript
{
  "common": {
    // Common UI elements (save, cancel, delete, etc.)
  },
  "nav": {
    // Navigation items
  },
  "login": {
    // Login page specific text
  },
  "dashboard": {
    // Dashboard specific text
  },
  "workflows": {
    // Workflows page specific text
  },
  "apps": {
    // Apps page specific text
  },
  "validation": {
    // Form validation messages
  },
  "messages": {
    // Success/error messages
  }
}
```

## Russian Translations Implemented

Key Russian translations include:

- **Common UI**: Сохранить (Save), Отмена (Cancel), Удалить (Delete), etc.
- **Navigation**: Панель управления (Dashboard), Рабочие процессы (Workflows), Приложения (Apps)
- **Login**: Вход в Shuffle (Login to Shuffle), Эл. почта (Email), Пароль (Password)
- **Workflows**: Рабочие процессы (Workflows), Создать новый рабочий процесс (Create New Workflow)
- **Forms**: Имя (Name), Описание (Description), Теги (Tags), Поиск (Search)

## Usage

### For Developers

To add translations to a component:

1. Import the translation hook:
```javascript
import { useTranslation } from 'react-i18next';
```

2. Use the hook in your component:
```javascript
const { t } = useTranslation();
```

3. Replace hardcoded text with translation keys:
```javascript
// Instead of: <Button>Save</Button>
<Button>{t('common.save')}</Button>

// Instead of: placeholder="Enter name"
placeholder={t('common.name')}
```

### For Users

Users can switch languages using the language selector in the header (both desktop and mobile versions). The language preference is automatically saved and will persist across sessions.

## Future Enhancements

1. **External Translation Files**: Move translations from inline configuration to separate JSON files in `public/locales/[lang]/translation.json`

2. **Additional Languages**: Easy to add more languages by:
   - Adding translation files
   - Updating the LanguageSwitcher component
   - Adding the language to the i18n configuration

3. **Translation Management**: Consider using translation management tools like Crowdin or Lokalise for collaborative translation

4. **RTL Support**: Add right-to-left language support if needed

5. **Pluralization**: Implement proper pluralization rules for languages that require it

6. **Date/Number Formatting**: Add locale-specific formatting for dates, numbers, and currencies

## Testing

To test the implementation:

1. Start the development server: `npm start`
2. Look for the language switcher in the header (flag icons with dropdown)
3. Switch between English (🇺🇸) and Russian (🇷🇺)
4. Verify that translated text appears correctly
5. Check that the language preference persists after page refresh

## Technical Notes

- The i18n configuration uses browser language detection as a fallback
- Language preference is stored in localStorage
- All existing i18next packages were already installed, only `react-i18next` was added
- The setup is designed to be performant with lazy loading capabilities for future external translation files