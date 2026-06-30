import { statusStyle } from '../utils.js';

export default function StatusBadge({ status }) {
  const { color, dot } = statusStyle(status);
  return (
    <span className="wr-status">
      <span className="wr-status-dot" style={{ background: dot }} />
      <span className="wr-status-label" style={{ color }}>
        {status}
      </span>
    </span>
  );
}
