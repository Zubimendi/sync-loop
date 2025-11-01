import { NavLink, useNavigate } from "react-router-dom";
import { CogIcon, CubeIcon, TableCellsIcon, ArrowRightIcon } from "@heroicons/react/24/outline";
import axios from "axios";

export default function Sidebar() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    await axios.post("/api/v1/logout"); // we’ll build this next
    document.cookie = "token=; path=/; max-age=0";
    navigate("/login");
  };

  const nav = [
    { name: "Connectors", href: "/connectors", icon: CubeIcon },
    { name: "Jobs", href: "/jobs", icon: TableCellsIcon },
    { name: "Settings", href: "/settings", icon: CogIcon },
  ];

  return (
    <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col">
      <div className="p-4 text-xl font-bold bg-gradient-to-r from-indigo-500 to-slate-500 text-transparent bg-clip-text">
        SyncLoop
      </div>
      <nav className="px-2 flex-1">
        {nav.map((i) => (
          <NavLink
            key={i.name}
            to={i.href}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition
              ${isActive ? "bg-indigo-600/20 text-indigo-300" : "text-slate-300 hover:bg-slate-800/50"}`
            }
          >
            <i.icon className="w-5 h-5" />
            {i.name}
          </NavLink>
        ))}
      </nav>

      <div className="p-2 border-t border-slate-800">
        <button
          onClick={handleLogout}
          className="flex items-center gap-3 w-full px-3 py-2 rounded-md text-sm font-medium text-slate-400 hover:bg-slate-800/50 hover:text-slate-200 transition"
        >
          <ArrowRightIcon className="w-5 h-5" />
          Logout
        </button>
      </div>
    </aside>
  );
}