import { Routes, Route, useNavigate } from 'react-router-dom';
import IncidentList from './pages/IncidentList.jsx';
import IncidentDetail from './pages/IncidentDetail.jsx';

function Brand() {
  const navigate = useNavigate();
  return (
    <button className="wr-brand" onClick={() => navigate('/')}>
      <span className="wr-brand-mark">
        <span className="wr-brand-dot" />
      </span>
      <span className="wr-brand-name">War Room</span>
    </button>
  );
}

export default function App() {
  return (
    <div className="wr-app">
      <div className="wr-topbar">
        <Brand />
      </div>
      <div className="wr-main">
        <Routes>
          <Route path="/" element={<IncidentList />} />
          <Route path="/incidents/:id" element={<IncidentDetail />} />
          <Route path="/incidents/:id/:tab" element={<IncidentDetail />} />
        </Routes>
      </div>
    </div>
  );
}
