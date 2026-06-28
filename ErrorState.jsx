export default function ErrorState({ message, onRetry }) {
  return (
    <div className="wr-error">
      <div>{message || 'Something went wrong loading this page.'}</div>
      {onRetry && (
        <button className="wr-retry" onClick={onRetry}>
          Try again
        </button>
      )}
    </div>
  );
}
