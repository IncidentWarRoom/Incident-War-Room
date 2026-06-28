import EmptyState from '../components/EmptyState.jsx';

export default function PhotosTab({ images }) {
  if (!images || images.length === 0) {
    return <EmptyState message="No photos have been attached to this incident." />;
  }

  return (
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))',
        gap: 14,
      }}
    >
      {images.map((img) => (
        <figure
          key={img.eventId}
          style={{
            margin: 0,
            border: '1px solid #ebe9e6',
            borderRadius: 10,
            overflow: 'hidden',
            background: '#fff',
          }}
        >
          <img
            src={img.url}
            alt={img.message || 'incident photo'}
            style={{ width: '100%', height: 130, objectFit: 'cover', display: 'block' }}
          />
          <figcaption style={{ padding: '8px 10px' }}>
            <div style={{ fontSize: 12.5, color: '#33332f', marginBottom: 2 }}>
              {img.message || 'No caption'}
            </div>
            <div style={{ fontSize: 11, color: '#a8a69f' }}>{img.username}</div>
          </figcaption>
        </figure>
      ))}
    </div>
  );
}
