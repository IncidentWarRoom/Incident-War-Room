import { severityStyle } from '../utils.js';

export default function SeverityBadge({ severity }) {
  const { color, bg } = severityStyle(severity);
  return (
    <span className="wr-badge" style={{ color, background: bg }}>
      {severity}
    </span>
  );
}
