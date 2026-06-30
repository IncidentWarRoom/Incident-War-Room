import { colorForUsername, initials } from '../utils.js';

export default function Avatar({ username, size = 22 }) {
  return (
    <span
      className="wr-avatar"
      title={username}
      style={{
        width: size,
        height: size,
        fontSize: size * 0.43,
        background: colorForUsername(username),
      }}
    >
      {initials(username)}
    </span>
  );
}
