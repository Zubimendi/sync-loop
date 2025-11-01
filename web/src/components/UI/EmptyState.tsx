export default function EmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <div className="text-center py-16">
      <div className="text-6xl mb-4">🗂️</div>
      <h3 className="text-xl font-semibold text-slate-700 dark:text-slate-300">No connectors yet</h3>
      <p className="text-slate-500 dark:text-slate-400 mt-2">Add your first source to get started.</p>
      <button onClick={onAdd} className="mt-4 btn-indigo">
        Add Connector
      </button>
    </div>
  );
}