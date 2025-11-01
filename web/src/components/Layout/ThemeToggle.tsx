import { SunIcon, MoonIcon } from "@heroicons/react/24/outline";
import useLocalStorage from "../../hooks/useLocalStorage";
import { useEffect } from "react";

export default function ThemeToggle() {
  const [dark, setDark] = useLocalStorage("dark", false);

  useEffect(() => {
    if (dark) document.documentElement.classList.add("dark");
    else document.documentElement.classList.remove("dark");
  }, [dark]);

  return (
    <button
      onClick={() => setDark(!dark)}
      className="p-2 rounded-lg bg-slate-800/50 hover:bg-slate-700/70 text-slate-300 hover:text-white transition"
      aria-label="Toggle dark mode"
    >
      {dark ? <SunIcon className="w-5 h-5" /> : <MoonIcon className="w-5 h-5" />}
    </button>
  );
}