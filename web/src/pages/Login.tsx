import { useState } from "react";
import axios from "axios";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await axios.post("/api/v1/login", { email, password });
      window.location.href = "/connectors";
    } catch {
      setError("Invalid credentials");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-900 to-slate-900">
      {/* subtle orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-gradient-radial from-indigo-500/20 to-transparent rounded-full blur-3xl animate-pulse" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-gradient-radial from-slate-500/20 to-transparent rounded-full blur-3xl animate-pulse delay-1000" />

      <form
        onSubmit={handleSubmit}
        className="relative z-10 w-full max-w-md p-8 space-y-6 bg-white/10 dark:bg-black/20 backdrop-blur-xl border border-white/20 dark:border-white/10 rounded-2xl shadow-2xl"
      >
        <div className="text-center">
          <h1 className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-indigo-300 to-slate-300">
            Welcome back
          </h1>
          <p className="mt-2 text-sm text-white/70">Sign in to SyncLoop</p>
        </div>

        {error && (
          <div className="bg-red-500/20 border border-red-400/30 text-red-200 text-sm px-4 py-2 rounded-lg">
            {error}
          </div>
        )}

        <div className="space-y-4">
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email address"
            type="email"
            required
            className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white placeholder-white/60 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-transparent transition"
          />
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password"
            type="password"
            required
            className="w-full px-4 py-3 rounded-lg bg-white/10 border border-white/20 text-white placeholder-white/60 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-transparent transition"
          />
        </div>

        <button
          type="submit"
          className="w-full py-3 rounded-lg bg-gradient-to-r from-indigo-500 to-slate-500 hover:from-indigo-600 hover:to-slate-600 text-white font-semibold shadow-lg hover:shadow-indigo-500/30 focus:outline-none focus:ring-2 focus:ring-indigo-400 transition transform hover:-translate-y-0.5"
        >
          Sign in
        </button>

        <p className="text-center text-xs text-white/50">
          By signing in you agree to our terms & privacy policy.
        </p>
      </form>
    </div>
  );
}