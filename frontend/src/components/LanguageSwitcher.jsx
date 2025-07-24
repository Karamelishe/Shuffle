import React from 'react';
import { useTranslation } from 'react-i18next';
import { 
  Select, 
  MenuItem, 
  FormControl, 
  Box,
  Typography 
} from '@mui/material';
import { LanguageOutlined as LanguageIcon } from '@mui/icons-material';

const LanguageSwitcher = ({ theme }) => {
  const { i18n, t } = useTranslation();

  const handleLanguageChange = (event) => {
    const newLanguage = event.target.value;
    i18n.changeLanguage(newLanguage);
    localStorage.setItem('i18nextLng', newLanguage);
  };

  const languages = [
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' }
  ];

  const currentLanguage = languages.find(lang => lang.code === i18n.language) || languages[0];

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', minWidth: 120 }}>
      <LanguageIcon 
        sx={{ 
          mr: 1, 
          color: theme?.palette?.text?.primary || '#fff',
          fontSize: '1.2rem'
        }} 
      />
      <FormControl size="small" variant="outlined">
        <Select
          value={i18n.language}
          onChange={handleLanguageChange}
          sx={{
            color: theme?.palette?.text?.primary || '#fff',
            '& .MuiOutlinedInput-notchedOutline': {
              borderColor: theme?.palette?.divider || 'rgba(255, 255, 255, 0.23)',
            },
            '& .MuiSvgIcon-root': {
              color: theme?.palette?.text?.primary || '#fff',
            },
            '&:hover .MuiOutlinedInput-notchedOutline': {
              borderColor: theme?.palette?.primary?.main || '#ff8544',
            },
            '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
              borderColor: theme?.palette?.primary?.main || '#ff8544',
            },
            minWidth: 100,
            height: 35,
          }}
          MenuProps={{
            PaperProps: {
              sx: {
                bgcolor: theme?.palette?.backgroundColor || '#1f2023',
                '& .MuiMenuItem-root': {
                  color: theme?.palette?.text?.primary || '#fff',
                  '&:hover': {
                    backgroundColor: theme?.palette?.hoverColor || 'rgba(255, 133, 68, 0.1)',
                  },
                },
              },
            },
          }}
        >
          {languages.map((language) => (
            <MenuItem key={language.code} value={language.code}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <span style={{ fontSize: '1.1em' }}>{language.flag}</span>
                <Typography variant="body2">{language.name}</Typography>
              </Box>
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
};

export default LanguageSwitcher;