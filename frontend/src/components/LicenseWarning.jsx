import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Alert,
  AlertTitle,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
  LinearProgress,
  Chip,
  Stack,
} from '@mui/material';
import {
  Warning as WarningIcon,
  Error as ErrorIcon,
  Info as InfoIcon,
  VpnKey as VpnKeyIcon,
  Upgrade as UpgradeIcon,
} from '@mui/icons-material';

const LicenseWarning = ({ 
  globalUrl, 
  userdata, 
  onLicenseActivate, 
  severity = 'warning',
  message = '',
  showDialog = false,
  onClose = () => {},
  licenseInfo = null 
}) => {
  const { t } = useTranslation();
  const [open, setOpen] = useState(showDialog);

  useEffect(() => {
    setOpen(showDialog);
  }, [showDialog]);

  const handleClose = () => {
    setOpen(false);
    onClose();
  };

  const handleActivateLicense = () => {
    if (onLicenseActivate) {
      onLicenseActivate();
    }
    handleClose();
  };

  const getIcon = () => {
    switch (severity) {
      case 'error':
        return <ErrorIcon />;
      case 'info':
        return <InfoIcon />;
      default:
        return <WarningIcon />;
    }
  };

  const getProgressColor = (percentage) => {
    if (percentage >= 90) return 'error';
    if (percentage >= 70) return 'warning';
    return 'primary';
  };

  const calculateUsagePercentage = (current, max) => {
    if (max === -1) return 0; // unlimited
    return Math.min(100, (current / max) * 100);
  };

  // Если есть информация о лицензии, показываем детальные предупреждения
  if (licenseInfo) {
    const daysRemaining = Math.ceil((new Date(licenseInfo.expires_at) - new Date()) / (1000 * 60 * 60 * 24));
    const isExpiringSoon = daysRemaining <= 30 && daysRemaining > 0;
    const isExpired = daysRemaining <= 0;

    if (isExpired) {
      return (
        <Alert severity="error" sx={{ mb: 2 }}>
          <AlertTitle>{t('license.expired')}</AlertTitle>
          <Typography variant="body2">
            {t('license.expired_message')}
          </Typography>
          <Box sx={{ mt: 1 }}>
            <Button
              size="small"
              variant="outlined"
              startIcon={<VpnKeyIcon />}
              onClick={handleActivateLicense}
            >
              {t('license.renew_license')}
            </Button>
          </Box>
        </Alert>
      );
    }

    if (isExpiringSoon) {
      return (
        <Alert severity="warning" sx={{ mb: 2 }}>
          <AlertTitle>{t('license.expiring_soon')}</AlertTitle>
          <Typography variant="body2">
            {t('license.expiring_soon_message', { days: daysRemaining })}
          </Typography>
          <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
            <LinearProgress
              variant="determinate"
              value={Math.max(0, (daysRemaining / 30) * 100)}
              sx={{ flexGrow: 1, height: 6, borderRadius: 3 }}
              color={getProgressColor(100 - (daysRemaining / 30) * 100)}
            />
            <Chip
              label={t('license.days_remaining', { days: daysRemaining })}
              size="small"
              color="warning"
            />
          </Box>
          <Box sx={{ mt: 1 }}>
            <Button
              size="small"
              variant="outlined"
              startIcon={<UpgradeIcon />}
              onClick={handleActivateLicense}
            >
              {t('license.renew_license')}
            </Button>
          </Box>
        </Alert>
      );
    }

    // Проверяем лимиты использования
    const usageWarnings = [];
    
    // Здесь можно добавить логику проверки текущего использования
    // Пока что показываем только если лимиты близки к исчерпанию
    
    if (usageWarnings.length > 0) {
      return (
        <Alert severity="warning" sx={{ mb: 2 }}>
          <AlertTitle>{t('license.usage_warning')}</AlertTitle>
          <Stack spacing={1}>
            {usageWarnings.map((warning, index) => (
              <Typography key={index} variant="body2">
                {warning}
              </Typography>
            ))}
          </Stack>
          <Box sx={{ mt: 1 }}>
            <Button
              size="small"
              variant="outlined"
              startIcon={<UpgradeIcon />}
              onClick={handleActivateLicense}
            >
              {t('license.upgrade_license')}
            </Button>
          </Box>
        </Alert>
      );
    }
  }

  // Базовое предупреждение
  if (message) {
    return (
      <Alert severity={severity} sx={{ mb: 2 }}>
        <AlertTitle>
          {severity === 'error' ? t('license.license_error') : t('license.license_warning')}
        </AlertTitle>
        <Typography variant="body2">
          {message}
        </Typography>
        <Box sx={{ mt: 1 }}>
          <Button
            size="small"
            variant="outlined"
            startIcon={<VpnKeyIcon />}
            onClick={handleActivateLicense}
          >
            {t('license.activate_license')}
          </Button>
        </Box>
      </Alert>
    );
  }

  // Диалог с предупреждением
  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        {getIcon()}
        {severity === 'error' ? t('license.license_error') : t('license.license_warning')}
      </DialogTitle>
      <DialogContent>
        <Typography variant="body1" paragraph>
          {message || t('license.default_warning_message')}
        </Typography>
        
        {licenseInfo && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="h6" gutterBottom>
              {t('license.current_license_info')}
            </Typography>
            <Stack spacing={1}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2" color="textSecondary">
                  {t('license.type')}:
                </Typography>
                <Chip
                  label={t(`license.type_${licenseInfo.type}`)}
                  size="small"
                  color="primary"
                  variant="outlined"
                />
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2" color="textSecondary">
                  {t('license.expires_at')}:
                </Typography>
                <Typography variant="body2">
                  {new Date(licenseInfo.expires_at).toLocaleDateString('ru-RU')}
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2" color="textSecondary">
                  {t('license.status')}:
                </Typography>
                <Chip
                  label={t(`license.status_${licenseInfo.status}`)}
                  size="small"
                  color={licenseInfo.status === 'active' ? 'success' : 'error'}
                  variant="outlined"
                />
              </Box>
            </Stack>
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>
          {t('common.close')}
        </Button>
        <Button
          onClick={handleActivateLicense}
          variant="contained"
          startIcon={<VpnKeyIcon />}
        >
          {licenseInfo ? t('license.manage_license') : t('license.activate_license')}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default LicenseWarning;