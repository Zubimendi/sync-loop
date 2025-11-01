import { Routes, Route, Navigate } from "react-router-dom";
import Sidebar from "./components/Layout/Sidebar";
import Connectors from "./pages/Connectors";
import useAuth from "./hooks/useAuth";

function App() {
  const { user, loading } = useAuth();
  if (loading) return <div className="p-6">Loading…</div>;
  if (!user) return <Navigate to="/login" replace />; // now safe inside Router
  return (
    <div className="flex h-screen bg-gray-50 dark:bg-gray-900">
      <Sidebar />
      <main className="flex-1 overflow-auto">
        <Routes>
          <Route path="/connectors" element={<Connectors />} />
          <Route path="/" element={<Navigate to="/connectors" replace />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;