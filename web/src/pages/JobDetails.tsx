import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import axios from "axios";
import { ArrowLeftIcon, ArrowPathIcon, ClockIcon, XCircleIcon, CheckCircleIcon } from "@heroicons/react/24/solid";

type JobStatus = {
  workflow_id: string;
  status: string;
  start_time: string;
  close_time?: string;
  failure_reason?: string;
  table?: string;
};

type HistoryEvent = {
  event_id: number;
  event_time: string;
  event_type: string;
  details?: any;
};

export default function JobDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [status, setStatus] = useState<JobStatus | null>(null);
  const [history, setHistory] = useState<HistoryEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);

  const loadStatus = async () => {
    if (!id) return;
    
    try {
      const { data } = await axios.get(`/api/v1/jobs/${id}/status`);
      setStatus(data);
      
      // Convert history iterator to array if it exists
      if (data.history) {
        const events: HistoryEvent[] = [];
        const historyIter = data.history;
        if (historyIter && typeof historyIter === 'object') {
          // Handle different history formats
          if (Array.isArray(historyIter)) {
            events.push(...historyIter);
          } else if (historyIter.events) {
            events.push(...historyIter.events);
          }
        }
        setHistory(events);
      }
      
      // Stop auto-refresh if job is completed/failed
      if (data.status === "COMPLETED" || data.status === "FAILED" || data.status === "CANCELLED") {
        setAutoRefresh(false);
      }
    } catch (error) {
      console.error("Failed to load job status:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStatus();
    
    if (autoRefresh) {
      const interval = setInterval(loadStatus, 3000); // Refresh every 3 seconds
      return () => clearInterval(interval);
    }
  }, [id, autoRefresh]);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "RUNNING":
        return <ClockIcon className="w-6 h-6 text-blue-500 animate-pulse" />;
      case "COMPLETED":
        return <CheckCircleIcon className="w-6 h-6 text-green-500" />;
      case "FAILED":
        return <XCircleIcon className="w-6 h-6 text-red-500" />;
      default:
        return <ClockIcon className="w-6 h-6 text-gray-500" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "RUNNING":
        return "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300";
      case "COMPLETED":
        return "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300";
      case "FAILED":
        return "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300";
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
    }
  };

  const formatDuration = (start: string, end?: string) => {
    const startTime = new Date(start);
    const endTime = end ? new Date(end) : new Date();
    const duration = endTime.getTime() - startTime.getTime();
    
    const seconds = Math.floor(duration / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  };

  if (loading) {
    return (
      <div className="p-6 flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
          <p className="mt-4 text-slate-600">Loading job details...</p>
        </div>
      </div>
    );
  }

  if (!status) {
    return (
      <div className="p-6">
        <div className="text-center py-16">
          <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-4">Job Not Found</h2>
          <p className="text-slate-600 mb-6">The job you're looking for doesn't exist or has been deleted.</p>
          <Link to="/jobs" className="btn-indigo">
            <ArrowLeftIcon className="w-4 h-4 inline mr-2" />
            Back to Jobs
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center">
          <button 
            onClick={() => navigate("/jobs")} 
            className="mr-4 p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700"
          >
            <ArrowLeftIcon className="w-5 h-5" />
          </button>
          <h1 className="text-2xl font-bold text-slate-800 dark:text-slate-100">
            Job Details
          </h1>
        </div>
        
        <div className="flex items-center gap-3">
          {status.status === "RUNNING" && (
            <button
              onClick={() => setAutoRefresh(!autoRefresh)}
              className={`px-3 py-1 text-sm rounded-lg ${
                autoRefresh 
                  ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300' 
                  : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300'
              }`}
            >
              {autoRefresh ? 'Auto-refresh ON' : 'Auto-refresh OFF'}
            </button>
          )}
          
          <span className={`px-3 py-1 text-sm font-medium rounded-full ${getStatusColor(status.status)}`}>
            {status.status}
          </span>
        </div>
      </div>

      {/* Status Card */}
      <div className="bg-white/30 dark:bg-slate-800/40 backdrop-blur-lg border border-white/20 dark:border-slate-700 rounded-2xl shadow-lg p-6 mb-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center">
            {getStatusIcon(status.status)}
            <div className="ml-3">
              <h2 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
                {status.workflow_id}
              </h2>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                {status.table && `Table: ${status.table}`}
              </p>
            </div>
          </div>
          
          <div className="text-right">
            <p className="text-sm text-slate-500 dark:text-slate-400">Duration</p>
            <p className="text-lg font-semibold text-slate-800 dark:text-slate-100">
              {formatDuration(status.start_time, status.close_time)}
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
          <div>
            <p className="text-sm text-slate-500 dark:text-slate-400">Started</p>
            <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
              {new Date(status.start_time).toLocaleString()}
            </p>
          </div>
          
          {status.close_time && (
            <div>
              <p className="text-sm text-slate-500 dark:text-slate-400">Finished</p>
              <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                {new Date(status.close_time).toLocaleString()}
              </p>
            </div>
          )}
          
          {status.failure_reason && (
            <div className="md:col-span-3">
              <p className="text-sm text-slate-500 dark:text-slate-400 mb-2">Failure Reason</p>
              <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
                <p className="text-sm text-red-800 dark:text-red-200">
                  {status.failure_reason}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* History Timeline */}
      {history.length > 0 && (
        <div className="bg-white/30 dark:bg-slate-800/40 backdrop-blur-lg border border-white/20 dark:border-slate-700 rounded-2xl shadow-lg p-6">
          <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-4">
            Execution History
          </h3>
          
          <div className="space-y-3">
            {history.map((event, index) => (
              <div key={index} className="flex items-start">
                <div className="flex-shrink-0 w-8 h-8 bg-indigo-100 dark:bg-indigo-900/30 rounded-full flex items-center justify-center">
                  <span className="text-xs font-medium text-indigo-600 dark:text-indigo-400">
                    {index + 1}
                  </span>
                </div>
                
                <div className="ml-4 flex-1">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                      {event.event_type}
                    </p>
                    <p className="text-xs text-slate-500 dark:text-slate-400">
                      {new Date(event.event_time).toLocaleTimeString()}
                    </p>
                  </div>
                  
                  {event.details && (
                    <div className="mt-2 text-xs text-slate-600 dark:text-slate-400 bg-slate-50 dark:bg-slate-700/50 rounded p-2">
                      <pre className="whitespace-pre-wrap">{JSON.stringify(event.details, null, 2)}</pre>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {history.length === 0 && (
        <div className="text-center py-12">
          <div className="text-4xl mb-4">📋</div>
          <h3 className="text-lg font-semibold text-slate-700 dark:text-slate-300 mb-2">
            No History Available
          </h3>
          <p className="text-slate-500 dark:text-slate-400">
            Workflow history will appear here as the job progresses.
          </p>
        </div>
      )}
    </div>
  );
}