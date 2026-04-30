import { t } from '../../../lib/i18n';

export default function StatusBadge(props: { connected: boolean; authValid: boolean; onClick?: () => void }) {
  const status = () => {
    if (!props.connected) return 'offline';
    if (!props.authValid) return 'warning';
    return 'online';
  };

  const text = () => {
    if (!props.connected) return t('Offline');
    if (!props.authValid) return t('Invalid');
    return t('Connected');
  };

  return (
    <div class={`status-badge ${status()}`} onClick={() => props.onClick?.()}>
      <span class={`status-dot ${status()}`} />
      <span class="status-text">{text()}</span>
    </div>
  );
}
