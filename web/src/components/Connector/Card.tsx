import { CubeIcon } from "@heroicons/react/24/outline";

const logos: Record<string, string> = {
  pg: "🐘", mysql: "🐬", s3: "☁️", excel: "📊", gsheets: "📝", sf: "☁️", rest: "🔌",
};

export default function ConnectorCard({ conn, onClick }: any) {
  return (
    <div
      onClick={onClick}
      className="group relative p-6 bg-white/30 dark:bg-slate-800/40 backdrop-blur-lg border border-white/20 dark:border-slate-700 rounded-2xl shadow-lg hover:shadow-indigo-500/20 transition cursor-pointer"
    >
      <div className="flex items-center justify-between">
        <div className="text-4xl">{logos[conn.type] ?? "🔧"}</div>
        <span className="px-2 py-1 text-xs font-medium bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 rounded-full">
          {conn.type}
        </span>
      </div>
      <h3 className="mt-4 text-lg font-semibold text-slate-800 dark:text-slate-100">{conn.name}</h3>
      <p className="text-sm text-slate-500 dark:text-slate-400">Created {new Date(conn.created_at).toLocaleDateString()}</p>
      <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition">
        <CubeIcon className="w-5 h-5 text-slate-400" />
      </div>
    </div>
  );
}