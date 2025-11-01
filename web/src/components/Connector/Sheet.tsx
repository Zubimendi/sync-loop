import { Fragment, useEffect, useState } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XMarkIcon } from "@heroicons/react/24/outline";
import axios from "axios";

type Props = { open: boolean; onClose: () => void; onSuccess: () => void };

const portDefaults: Record<string, string> = {
  pg: "5432",
  mysql: "3306",
  s3: "443",
  excel: "",
  gsheets: "",
  rest: "",
};

export default function ConnectorSheet({ open, onClose, onSuccess }: Props) {
  const [form, setForm] = useState({
    name: "",
    type: "pg",
    host: "",
    port: portDefaults.pg,
    database: "",
    user: "",
    password: "",
  });

  // auto-switch port when type changes
  useEffect(() => {
    setForm((f) => ({ ...f, port: portDefaults[f.type] }));
  }, [form.type]);

  const handleSave = async () => {
    await axios.post("/api/v1/connectors", {
      name: form.name,
      type: form.type,
      config: { host: form.host, port: Number(form.port), database: form.database, user: form.user, password: form.password },
    });
    onSuccess();
    onClose();
  };

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child as={Fragment} enter="ease-in-out duration-500" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in-out duration-500" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black/30 backdrop-blur-sm" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10">
              <Transition.Child as={Fragment} enter="transform transition ease-in-out duration-500 sm:duration-700" enterFrom="translate-x-full" enterTo="translate-x-0" leave="transform transition ease-in-out duration-500 sm:duration-700" leaveFrom="translate-x-0" leaveTo="translate-x-full">
                <Dialog.Panel className="pointer-events-auto w-screen max-w-md">
                  <div className="flex h-full flex-col overflow-y-scroll bg-white/90 dark:bg-slate-900/90 backdrop-blur-xl py-6 shadow-2xl">
                    <div className="px-4 sm:px-6">
                      <div className="flex items-center justify-between">
                        <Dialog.Title className="text-lg font-semibold text-slate-900 dark:text-white">Add Connector</Dialog.Title>
                        <button type="button" className="rounded-md text-slate-400 hover:text-slate-500" onClick={onClose}>
                          <XMarkIcon className="w-6 h-6" />
                        </button>
                      </div>
                    </div>

                    <div className="relative mt-6 flex-1 px-4 sm:px-6">
                      <div className="space-y-4">
                        <input
                          value={form.name}
                          onChange={(e) => setForm({ ...form, name: e.target.value })}
                          placeholder="Name"
                          className="input-glass py-2 px-3 text-base"
                        />
                        <select
                          value={form.type}
                          onChange={(e) => setForm({ ...form, type: e.target.value })}
                          className="input-glass py-2 px-3 text-base"
                        >
                          <option value="pg">Postgres</option>
                          <option value="mysql">MySQL</option>
                          <option value="s3">S3</option>
                          <option value="excel">Excel</option>
                          <option value="gsheets">Google Sheets</option>
                        </select>
                        <input
                          value={form.host}
                          onChange={(e) => setForm({ ...form, host: e.target.value })}
                          placeholder="Host"
                          className="input-glass py-2 px-3 text-base"
                        />
                        <input
                          value={form.port}
                          onChange={(e) => setForm({ ...form, port: e.target.value })}
                          placeholder="Port"
                          className="input-glass py-2 px-3 text-base"
                        />
                        <input
                          value={form.database}
                          onChange={(e) => setForm({ ...form, database: e.target.value })}
                          placeholder="Database"
                          className="input-glass py-2 px-3 text-base"
                        />
                        <input
                          value={form.user}
                          onChange={(e) => setForm({ ...form, user: e.target.value })}
                          placeholder="User"
                          className="input-glass py-2 px-3 text-base"
                        />
                        <input
                          value={form.password}
                          onChange={(e) => setForm({ ...form, password: e.target.value })}
                          placeholder="Password"
                          type="password"
                          className="input-glass py-2 px-3 text-base"
                        />
                      </div>
                    </div>

                    <div className="px-4 sm:px-6 mt-6 flex gap-3">
                      <button type="button" onClick={handleSave} className="flex-1 py-2.5 rounded-lg bg-gradient-to-r from-indigo-500 to-pink-500 hover:from-indigo-600 hover:to-pink-600 text-white font-semibold shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5">
                        Save
                      </button>
                      <button type="button" onClick={onClose} className="flex-1 py-2.5 rounded-lg bg-white/10 hover:bg-white/20 dark:bg-slate-700/50 dark:hover:bg-slate-700/70 text-slate-800 dark:text-slate-200 border border-white/20 dark:border-slate-600 font-medium transition">
                        Cancel
                      </button>
                    </div>
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  );
}
