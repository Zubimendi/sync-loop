import { useEffect, useState } from "react";
import { PlayIcon, PauseIcon, ArrowPathIcon } from "@heroicons/react/24/solid";
import axios from "axios";
import { Link } from "react-router-dom";

type Job = {
  id: string;
  type: string;
  status: string;
  start_time: string;
  table: string;
  schedule?: ScheduleInfo;
};

type ScheduleInfo = {
  id: string;
  cron_expr: string;
  is_active: boolean;
  last_run_time: string;
  next_run_time: string;
};

export default function Jobs() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [running, setRunning] = useState(false);
  const [stopping, setStopping] = useState<Set<string>>(new Set());
  const [showScheduleModal, setShowScheduleModal] = useState(false);
  const [scheduleForm, setScheduleForm] = useState({
    table: "users",
    cron_expr: "* * * * *",
    is_active: true,
  });

  const cancel = async (id: string) => {
    setStopping((s) => new Set(s).add(id));
    await axios.post("/api/v1/jobs/cancel", { workflow_id: id });

    // poll every 500 ms until status changes
    const poll = setInterval(async () => {
      const { data } = await axios.get("/api/v1/jobs");
      const updated = data.jobs.find((j: Job) => j.id === id);
      if (!updated || updated.status !== "RUNNING") {
        clearInterval(poll);
        setStopping((s) => {
          const ns = new Set(s);
          ns.delete(id);
          return ns;
        });
        setJobs(data.jobs); // final list
      }
    }, 500);
  };

  const load = async () => {
    const { data } = await axios.get("/api/v1/jobs");
    setJobs(data.jobs);
    setLoading(false);
  };

  useEffect(() => {
    load();
    const t = setInterval(load, 5000);
    return () => clearInterval(t);
  }, []);

  const runNow = async (incremental = false) => {
    setRunning(true);
    try {
      await axios.post("/api/v1/jobs/run-now", { table: "users", incremental });
      await load(); // wait for list to update
    } catch {
      alert("Could not start job");
    } finally {
      setRunning(false);
    }
  };

  const retryJob = async (job: Job) => {
    try {
      await axios.post("/api/v1/jobs/retry", {
        workflow_id: job.id,
        run_id: job.id.split("-").pop(),
      });
      await load();
    } catch (error) {
      alert("Could not retry job");
    }
  };

  const createSchedule = async () => {
    try {
      await axios.post("/api/v1/jobs/schedule", scheduleForm);
      setShowScheduleModal(false);
      await load();
    } catch (error) {
      alert("Could not create schedule");
    }
  };

  const toggleSchedule = async (job: Job, pause: boolean) => {
    try {
      await axios.post("/api/v1/jobs/schedule/toggle", {
        schedule_id: job.schedule?.id,
        pause,
      });
      await load();
    } catch (error) {
      alert("Could not update schedule");
    }
  };

  const statusColor = (s: string) => {
    if (s === "RUNNING")
      return "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300";
    if (s === "COMPLETED")
      return "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300";
    if (s === "FAILED")
      return "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300";
    return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
  };

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-800 dark:text-slate-100">
          Jobs
        </h1>
        <div className="flex gap-2">
          <button
            onClick={() => setShowScheduleModal(true)}
            className="btn-secondary"
            disabled={running}
          >
            <PauseIcon className="w-5 h-5 inline mr-2" />
            Schedule
          </button>
          <button
            onClick={() => runNow(false)}
            disabled={running}
            className="btn-indigo"
          >
            {running ? (
              <>
                <svg
                  className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                  ></path>
                </svg>
                Starting...
              </>
            ) : (
              <>
                <PlayIcon className="w-5 h-5 inline mr-2" />
                Run Full Sync
              </>
            )}
          </button>
          <button
            onClick={() => runNow(true)}
            disabled={running}
            className="btn-green"
          >
            <PlayIcon className="w-5 h-5 inline mr-2" />
            Run Incremental
          </button>
        </div>
      </div>

      {loading && <p className="text-slate-500">Loading…</p>}
      {!loading && jobs.length === 0 && (
        <div className="text-center py-16">
          <div className="text-6xl mb-4">⏱️</div>
          <h3 className="text-xl font-semibold text-slate-700 dark:text-slate-300">
            No jobs yet
          </h3>
          <p className="text-slate-500 dark:text-slate-400 mt-2">
            Click "Run now" to start your first job.
          </p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {jobs
          .filter((j) => j.status !== "CANCELED" && j.status !== "TERMINATED")
          .map((j) => (
            <Link
              key={j.id}
              to={`/jobs/${j.id}`}
              className="block p-4 bg-white/30 dark:bg-slate-800/40 backdrop-blur-lg border border-white/20 dark:border-slate-700 rounded-2xl shadow-lg hover:shadow-xl transition-all duration-200 hover:border-indigo-400 dark:hover:border-indigo-500"
            >
              {" "}
              <div
                key={j.id}
                className="p-4 bg-white/30 dark:bg-slate-800/40 backdrop-blur-lg border border-white/20 dark:border-slate-700 rounded-2xl shadow-lg"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-slate-500 dark:text-slate-400">
                    {j.type}
                  </span>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${statusColor(
                      j.status
                    )}`}
                  >
                    {j.status}
                  </span>
                </div>
                <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
                  Table: {j.table}
                </h3>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  Started {new Date(j.start_time).toLocaleString()}
                </p>

                {j.status === "RUNNING" && (
                  <button
                    onClick={() => cancel(j.id)}
                    disabled={stopping.has(j.id)}
                    className="mt-3 bg-red-600 hover:bg-red-700 disabled:bg-red-400 text-white font-semibold py-2 px-4 rounded-md shadow-md transition-colors focus:outline-none focus:ring-2 focus:ring-red-400"
                  >
                    {stopping.has(j.id) ? "Cancelling…" : "Cancel"}
                  </button>
                )}

                {j.status === "FAILED" && (
                  <button
                    onClick={() => retryJob(j)}
                    className="mt-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md shadow-md transition-colors focus:outline-none focus:ring-2 focus:ring-blue-400"
                  >
                    <ArrowPathIcon className="w-4 h-4 inline mr-1" />
                    Retry
                  </button>
                )}

                {j.schedule && (
                  <div className="mt-3 pt-3 border-t border-slate-200 dark:border-slate-700">
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-slate-500">
                        Schedule: {j.schedule.cron_expr}
                      </span>
                      <button
                        onClick={() => toggleSchedule(j, j.schedule!.is_active)}
                        className={`text-xs px-2 py-1 rounded ${
                          j.schedule.is_active
                            ? "bg-red-100 text-red-700 hover:bg-red-200"
                            : "bg-green-100 text-green-700 hover:bg-green-200"
                        }`}
                      >
                        {j.schedule.is_active ? (
                          <>
                            <PauseIcon className="w-3 h-3 inline mr-1" />
                            Pause
                          </>
                        ) : (
                          <>
                            <PlayIcon className="w-3 h-3 inline mr-1" />
                            Resume
                          </>
                        )}
                      </button>
                    </div>
                    <p className="text-xs text-slate-500 mt-1">
                      Next run:{" "}
                      {new Date(j.schedule.next_run_time).toLocaleString()}
                    </p>
                  </div>
                )}
              </div>
            </Link>
          ))}
      </div>

      <button
        onClick={async () => {
          if (!confirm("Terminate ALL running jobs?")) return;
          await axios.post("/api/v1/jobs/terminate-all");
          load();
        }}
        className="bg-red-600 hover:bg-red-700 text-white font-semibold py-2 px-4 rounded-md shadow-md transition-colors focus:outline-none focus:ring-2 focus:ring-red-400 mt-5"
      >
        Terminate All Running
      </button>

      {/* Schedule Modal */}
      {showScheduleModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-slate-800 p-6 rounded-lg shadow-xl max-w-md w-full mx-4">
            <h2 className="text-xl font-bold mb-4">Create Schedule</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">Table</label>
                <input
                  type="text"
                  value={scheduleForm.table}
                  onChange={(e) =>
                    setScheduleForm({ ...scheduleForm, table: e.target.value })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                  placeholder="users"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">
                  Cron Expression
                </label>
                <input
                  type="text"
                  value={scheduleForm.cron_expr}
                  onChange={(e) =>
                    setScheduleForm({
                      ...scheduleForm,
                      cron_expr: e.target.value,
                    })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                  placeholder="* * * * *"
                />
                <p className="text-xs text-slate-500 mt-1">
                  * * * * * = every minute
                </p>
              </div>
              <div className="flex items-center">
                <input
                  type="checkbox"
                  checked={scheduleForm.is_active}
                  onChange={(e) =>
                    setScheduleForm({
                      ...scheduleForm,
                      is_active: e.target.checked,
                    })
                  }
                  className="mr-2"
                />
                <label className="text-sm">Active immediately</label>
              </div>
            </div>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setShowScheduleModal(false)}
                className="px-4 py-2 text-slate-600 hover:text-slate-800"
              >
                Cancel
              </button>
              <button
                onClick={createSchedule}
                className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
              >
                Create Schedule
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
