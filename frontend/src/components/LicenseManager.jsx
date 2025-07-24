import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  Card,
  CardContent,
  Grid,
  Alert,
  Chip,
  LinearProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Divider,
  Tooltip,
  IconButton,
  Snackbar,
} from '@mui/material';
import {
  Security as SecurityIcon,
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Warning as WarningIcon,
  Info as InfoIcon,
  Refresh as RefreshIcon,
  VpnKey as VpnKeyIcon,
  Business as BusinessIcon,
  People as PeopleIcon,
  AccountTree as WorkflowIcon,
  PlayArrow as ExecutionIcon,
  Schedule as ScheduleIcon,
  ContentCopy as CopyIcon,
} from '@mui/icons-material';

const LicenseManager = ({ globalUrl, userdata, theme }) => {
  const { t } = useTranslation();
  const [licenseKey, setLicenseKey] = useState('');
  const [licenseInfo, setLicenseInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [activationDialog, setActivationDialog] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');

  // Загружаем информацию о лицензии при монтировании компонента
  useEffect(() => {
    loadLicenseInfo();
  }, []);

  const loadLicenseInfo = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${globalUrl}/api/v1/license/info`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${userdata?.apikey}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        if (data.success) {
          setLicenseInfo(data.license);
        }
      } else if (response.status !== 404) {
        console.error('Failed to load license info');
      }
    } catch (err) {
      console.error('Error loading license info:', err);
    } finally {
      setLoading(false);
    }
  };

  const activateLicense = async () => {
    if (!licenseKey.trim()) {
      setError(t('license.key_required'));
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const response = await fetch(`${globalUrl}/api/v1/license/activate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${userdata?.apikey}`,
        },
        body: JSON.stringify({
          key: licenseKey.trim(),
        }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        setSuccess(t('license.activation_success'));
        setLicenseInfo(data.license);
        setLicenseKey('');
        setActivationDialog(false);
        showSnackbar(t('license.activation_success'));
      } else {
        setError(data.reason || t('license.activation_failed'));
      }
    } catch (err) {
      setError(t('license.activation_error'));
      console.error('License activation error:', err);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    showSnackbar(t('common.copied_to_clipboard'));
  };

  const showSnackbar = (message) => {
    setSnackbarMessage(message);
    setSnackbarOpen(true);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatLimit = (limit) => {
    return limit === -1 ? t('license.unlimited') : limit.toString();
  };

  const getLicenseStatusColor = (status) => {
    switch (status) {
      case 'active':
        return 'success';
      case 'expired':
        return 'error';
      case 'revoked':
        return 'error';
      default:
        return 'default';
    }
  };

  const getLicenseTypeColor = (type) => {
    switch (type) {
      case 'basic':
        return 'info';
      case 'professional':
        return 'warning';
      case 'enterprise':
        return 'success';
      default:
        return 'default';
    }
  };

  const getDaysRemaining = (expiresAt) => {
    const now = new Date();
    const expiry = new Date(expiresAt);
    const diffTime = expiry - now;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  const getProgressValue = (expiresAt, createdAt) => {
    const now = new Date();
    const created = new Date(createdAt);
    const expiry = new Date(expiresAt);
    const total = expiry - created;
    const elapsed = now - created;
    return Math.min(100, Math.max(0, (elapsed / total) * 100));
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <SecurityIcon />
        {t('license.title')}
      </Typography>

      {/* Статус лицензии */}
      {licenseInfo ? (
        <Grid container spacing={3}>
          {/* Основная информация о лицензии */}
          <Grid item xs={12} md={8}>
            <Card elevation={3}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                  <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <VpnKeyIcon />
                    {t('license.current_license')}
                  </Typography>
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    <Chip
                      label={t(`license.type_${licenseInfo.type}`)}
                      color={getLicenseTypeColor(licenseInfo.type)}
                      variant="outlined"
                    />
                    <Chip
                      label={t(`license.status_${licenseInfo.status}`)}
                      color={getLicenseStatusColor(licenseInfo.status)}
                      icon={licenseInfo.status === 'active' ? <CheckCircleIcon /> : <ErrorIcon />}
                    />
                  </Box>
                </Box>

                <Grid container spacing={2}>
                  <Grid item xs={12} sm={6}>
                    <Typography variant="body2" color="textSecondary">
                      {t('license.license_key')}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
                        {licenseInfo.key}
                      </Typography>
                      <Tooltip title={t('common.copy')}>
                        <IconButton size="small" onClick={() => copyToClipboard(licenseInfo.key)}>
                          <CopyIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </Grid>

                  <Grid item xs={12} sm={6}>
                    <Typography variant="body2" color="textSecondary">
                      {t('license.expires_at')}
                    </Typography>
                    <Typography variant="body1">
                      {formatDate(licenseInfo.expires_at)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {getDaysRemaining(licenseInfo.expires_at) > 0
                        ? t('license.days_remaining', { days: getDaysRemaining(licenseInfo.expires_at) })
                        : t('license.expired')}
                    </Typography>
                  </Grid>

                  <Grid item xs={12}>
                    <Box sx={{ mt: 1 }}>
                      <LinearProgress
                        variant="determinate"
                        value={getProgressValue(licenseInfo.expires_at, licenseInfo.created_at)}
                        sx={{
                          height: 8,
                          borderRadius: 4,
                          backgroundColor: theme.palette.grey[300],
                        }}
                      />
                    </Box>
                  </Grid>
                </Grid>

                <Divider sx={{ my: 2 }} />

                {/* Лимиты лицензии */}
                <Typography variant="h6" gutterBottom>
                  {t('license.limits')}
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={4}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <PeopleIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="textSecondary">
                          {t('license.max_users')}
                        </Typography>
                        <Typography variant="h6">
                          {formatLimit(licenseInfo.max_users)}
                        </Typography>
                      </Box>
                    </Box>
                  </Grid>

                  <Grid item xs={12} sm={4}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <WorkflowIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="textSecondary">
                          {t('license.max_workflows')}
                        </Typography>
                        <Typography variant="h6">
                          {formatLimit(licenseInfo.max_workflows)}
                        </Typography>
                      </Box>
                    </Box>
                  </Grid>

                  <Grid item xs={12} sm={4}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <ExecutionIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="textSecondary">
                          {t('license.max_executions')}
                        </Typography>
                        <Typography variant="h6">
                          {formatLimit(licenseInfo.max_executions)}
                        </Typography>
                      </Box>
                    </Box>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>

          {/* Функции лицензии */}
          <Grid item xs={12} md={4}>
            <Card elevation={3}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {t('license.features')}
                </Typography>
                <List dense>
                  {licenseInfo.features?.map((feature, index) => (
                    <ListItem key={index} sx={{ px: 0 }}>
                      <ListItemIcon sx={{ minWidth: 32 }}>
                        <CheckCircleIcon color="success" fontSize="small" />
                      </ListItemIcon>
                      <ListItemText
                        primary={
                          <Typography variant="body2">
                            {t(`license.feature_${feature.replace(':', '_')}`)}
                          </Typography>
                        }
                      />
                    </ListItem>
                  ))}
                </List>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      ) : (
        // Нет активной лицензии
        <Card elevation={3}>
          <CardContent sx={{ textAlign: 'center', py: 4 }}>
            <WarningIcon sx={{ fontSize: 64, color: 'warning.main', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              {t('license.no_license')}
            </Typography>
            <Typography variant="body1" color="textSecondary" paragraph>
              {t('license.no_license_description')}
            </Typography>
            <Button
              variant="contained"
              size="large"
              onClick={() => setActivationDialog(true)}
              sx={{ mt: 2 }}
            >
              {t('license.activate_license')}
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Кнопки управления */}
      <Box sx={{ mt: 3, display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
        <Button
          variant="outlined"
          startIcon={<RefreshIcon />}
          onClick={loadLicenseInfo}
          disabled={loading}
        >
          {t('common.refresh')}
        </Button>
        <Button
          variant="contained"
          startIcon={<VpnKeyIcon />}
          onClick={() => setActivationDialog(true)}
        >
          {licenseInfo ? t('license.change_license') : t('license.activate_license')}
        </Button>
      </Box>

      {/* Диалог активации лицензии */}
      <Dialog
        open={activationDialog}
        onClose={() => {
          setActivationDialog(false);
          setError('');
          setLicenseKey('');
        }}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          {t('license.activate_license')}
        </DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="textSecondary" paragraph>
            {t('license.activation_description')}
          </Typography>
          
          <TextField
            fullWidth
            label={t('license.license_key')}
            value={licenseKey}
            onChange={(e) => setLicenseKey(e.target.value)}
            placeholder="XXXX-XXXX-XXXX-XXXX"
            variant="outlined"
            sx={{ mb: 2 }}
            inputProps={{
              style: { fontFamily: 'monospace' }
            }}
          />

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setActivationDialog(false);
              setError('');
              setLicenseKey('');
            }}
          >
            {t('common.cancel')}
          </Button>
          <Button
            onClick={activateLicense}
            variant="contained"
            disabled={loading || !licenseKey.trim()}
          >
            {loading ? t('common.loading') : t('license.activate')}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar для уведомлений */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={3000}
        onClose={() => setSnackbarOpen(false)}
        message={snackbarMessage}
      />
    </Box>
  );
};

export default LicenseManager;