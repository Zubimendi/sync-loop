import { useEffect, useState } from "react";
import { PlusIcon } from "@heroicons/react/24/solid";
import { createConnector, listConnectors } from "../../lib/api";

const logos: Record<string, string> = {
  pg: "🐘",
  mysql: "🐬",
  s3: "☁️",
  excel: "📊",
  gsheets: "📝",
  sf: "☁️",
  rest: "🔌",
};

export default function ConnectorGrid() {
  const [list, setList] = useState<any[]>([]);
  const [sheet, setSheet] = useState(false);

  useEffect(() => {
    listConnectors().then(setList).catch(console.error);
  }, []);

  const handleAdd = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const payload = {
      name: fd.get("name"),
      type: fd.get("type"),
      config: {
        host: fd.get("host"),
        port: Number(fd.get("port")),
        database: fd.get("database"),
        user: fd.get("user"),
        password: fd.get("password"),
      },
    };
    await createConnector(payload);
    setSheet(false);
    setList(await listConnectors());
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">
          Connectors
        </h1>
        <button onClick={() => setSheet(true)} className="btn-indigo">
          <PlusIcon className="w-5 h-5 inline mr-2" />
          Add Connector
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {list.map((c) => (
          <div
            key={c.id}
            className="border rounded-lg p-4 dark:border-gray-700"
          >
            <div className="text-3xl mb-2">{logos[c.type] ?? "🔧"}</div>
            <div className="font-medium text-gray-900 dark:text-white">
              {c.name}
            </div>
            <div className="text-sm text-gray-500 dark:text-gray-400">
              {c.type}
            </div>
          </div>
        ))}
      </div>

      {sheet && (
        <div className="fixed inset-0 bg-black bg-opacity-30 flex justify-end">
          <div className="w-96 bg-white dark:bg-gray-900 p-6 shadow-xl">
            <h2 className="text-lg font-semibold mb-4">New Postgres Source</h2>
            <form onSubmit={handleAdd} className="space-y-3">
              <input
                name="name"
                placeholder="Name"
                required
                className="input"
              />
              <select name="type" required className="input">
                <option value="pg">Postgres</option>
                <option value="mysql">MySQL</option>
                <option value="s3">S3</option>
              </select>
              <input
                name="host"
                placeholder="Host"
                required
                className="input"
              />
              <input
                name="port"
                type="number"
                placeholder="Port"
                required
                className="input"
              />
              <input
                name="database"
                placeholder="Database"
                required
                className="input"
              />
              <input
                name="user"
                placeholder="User"
                required
                className="input"
              />
              <input
                name="password"
                type="password"
                placeholder="Password"
                required
                className="input"
              />
              <div className="flex gap-2">
                <button type="submit" className="btn-indigo">
                  Save
                </button>
                <button
                  type="button"
                  onClick={() => setSheet(false)}
                  className="btn-gray"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
