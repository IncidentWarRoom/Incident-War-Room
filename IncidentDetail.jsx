import { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { getIncident, getTimeline, getImages, ApiError } from '../api.js';
import { formatDate, formatDuration, shortRef } from '../utils.js';
import SeverityBadge from '../components/SeverityBadge.jsx';
import StatusBadge from '../components/StatusBadge.jsx';
import Avatar from '../components/Avatar.jsx';
import Skeleton from '../components/Skeleton.jsx';
import ErrorState from '../components/ErrorState.jsx';
import TimelineTab from './TimelineTab.jsx';
import PhotosTab from './PhotosTab.jsx';

export default function IncidentDetail() {
  const { id, tab: tabParam } = useParams();
  const navigate = useNavigate();
  const tab = tabParam === 'photos' ? 'photos' : 'timeline';

  const [status, setStatus] = useState('loading'); // loading | ready | error
  const [incident, setIncident] = useState(null);
  const [timeline, setTimeline] = useState(null);
  const [images, setImages] = useState([]);
  const [error, setError] = useState(null);

  const load = useCallback(async () => {
    setStatus('loading');
    setError(null);
    try {
      const [incidentData, timelineData, imagesData] = await Promise.all([
        getIncident(id),
        getTimeline(id),
        getImages(id),
      ]);
      setIncident(incidentData);
      setTimeline(timelineData);
      setImages(imagesData || []);
      setStatus('ready');
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setError("This incident doesn't exist.");
      } else if (err instanceof ApiError && err.status === 400) {
        setError('Invalid incident id.');
      } else {
        setError(err instanceof ApiError ? err.message : 'Unexpected error');
      }
      setStatus('error');
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  if (status === 'loading') return <DetailSkeleton />;
  if (status === 'error') {
    return (
      <div style={{ maxWidth: 760, margin: '0 auto', padding: '28px 24px' }}>
        <BackLink />
        <ErrorState message={error} onRetry={load} />
      </div>
    );
  }

  const created = formatDate(incident.createdAt);
  const closed = formatDate(incident.closedAt);
  const duration = formatDuration(timeline?.durationSeconds);
  const responders = timeline?.responders || [];

  return (
    <div style={{ maxWidth: 760, margin: '0 auto', padding: '28px 24px' }}>
      <BackLink />

      {/* ---- header ---- */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 6 }}>
        <span
          style={{
            fontFamily: 'ui-monospace, SF Mono, Menlo, Consolas, monospace',
            fontSize: 11,
            color: '#8a8983',
          }}
        >
          {shortRef(incident.id)}
        </span>
        <SeverityBadge severity={incident.severity} />
        <StatusBadge status={incident.status} />
      </div>
      <h1 style={{ fontSize: 21, fontWeight: 600, margin: '0 0 14px' }}>{incident.title}</h1>

      {/* ---- time span ---- */}
      <div style={{ fontSize: 13, color: '#6f6e69', marginBottom: 14 }}>
        Opened {created ? created.date : '-'} at {created ? created.time : '-'}
        {incident.closedAt ? (
          <>
            {' '}- closed {closed.date} at {closed.time}
            {duration ? <> (took {duration})</> : null}
          </>
        ) : (
          ' - ongoing'
        )}
      </div>

      {/* ---- links row ---- */}
      {(incident.telegraphUrls?.length > 0 || incident.reportUrl) && (
        <div style={{ display: 'flex', gap: 16, marginBottom: 16, fontSize: 13 }}>
          {incident.telegraphUrls?.length > 0 && (
            <a href={incident.telegraphUrls[0]} target="_blank" rel="noreferrer" style={{ color: '#0e7490' }}>
              View full timeline {'->'}
            </a>
          )}
          {incident.reportUrl && (
            <a href={incident.reportUrl} target="_blank" rel="noreferrer" style={{ color: '#0e7490' }}>
              Download report {'->'}
            </a>
          )}
        </div>
      )}

      {/* ---- responders ---- */}
      {responders.length > 0 && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 22 }}>
          <span style={{ fontSize: 12, color: '#a8a69f' }}>Responders</span>
          <div style={{ display: 'flex', gap: -4 }}>
            {responders.map((name) => (
              <Avatar key={name} username={name} />
            ))}
          </div>
        </div>
      )}

      {/* ---- tabs ---- */}
      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid #ebe9e6', marginBottom: 18 }}>
        <TabButton
          active={tab === 'timeline'}
          onClick={() => navigate(`/incidents/${id}/timeline`)}
        >
          Timeline
        </TabButton>
        <TabButton
          active={tab === 'photos'}
          onClick={() => navigate(`/incidents/${id}/photos`)}
        >
          Photos {images.length > 0 ? `(${images.length})` : ''}
        </TabButton>
      </div>

      {tab === 'timeline' ? (
        <TimelineTab events={timeline?.events} />
      ) : (
        <PhotosTab images={images} />
      )}
    </div>
  );
}

function BackLink() {
  return (
    <Link
      to="/"
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 6,
        fontSize: 12.5,
        color: '#6f6e69',
        marginBottom: 18,
      }}
    >
      <svg width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6">
        <path d="M10 13 5 8l5-5" />
      </svg>
      All incidents
    </Link>
  );
}

function TabButton({ active, onClick, children }) {
  return (
    <button
      onClick={onClick}
      style={{
        border: 'none',
        background: 'transparent',
        cursor: 'pointer',
        padding: '8px 4px',
        marginRight: 18,
        fontSize: 13,
        fontWeight: 500,
        color: active ? '#1d1d1b' : '#a8a69f',
        borderBottom: active ? '2px solid #1d1d1b' : '2px solid transparent',
        marginBottom: -1,
      }}
    >
      {children}
    </button>
  );
}

function DetailSkeleton() {
  return (
    <div style={{ maxWidth: 760, margin: '0 auto', padding: '28px 24px' }}>
      <Skeleton width={100} height={12} style={{ marginBottom: 18 }} />
      <Skeleton width={300} height={22} style={{ marginBottom: 10 }} />
      <Skeleton width={220} height={13} style={{ marginBottom: 24 }} />
      <Skeleton width="100%" height={80} />
    </div>
  );
}
