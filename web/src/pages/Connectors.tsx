import { useState, useEffect } from "react";
import { PlusIcon } from "@heroicons/react/24/solid";
import ConnectorCard from "../components/Connector/Card";
import ConnectorSheet from "../components/Connector/Sheet";
import EmptyState from "../components/UI/EmptyState";
import { listConnectors } from "../lib/api";

export default function Connectors() {
  const [list, setList] = useState<any[]>([]);
  const [sheet, setSheet] = useState(false);

  const load = async () => setList(await listConnectors());

  useEffect(() => { load(); }, []);

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-800 dark:text-slate-100">Connectors</h1>
        <button onClick={() => setSheet(true)} className="btn-indigo">
          <PlusIcon className="w-5 h-5 inline mr-2" />Add Connector
        </button>
      </div>

      {list.length === 0 ? (
        <EmptyState onAdd={() => setSheet(true)} />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {list.map((c) => (
            <ConnectorCard key={c.id} conn={c} onClick={() => {}} />
          ))}
        </div>
      )}

      <ConnectorSheet open={sheet} onClose={() => setSheet(false)} onSuccess={load} />
    </div>
  );
}