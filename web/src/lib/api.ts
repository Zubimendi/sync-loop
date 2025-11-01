import axios from "axios";
axios.defaults.baseURL = import.meta.env.VITE_API_URL;
axios.defaults.withCredentials = true; // cookies
export async function listConnectors() {
  const { data } = await axios.get("/api/v1/connectors");
  return data.connectors;
}
export async function createConnector(payload: any) {
  const { data } = await axios.post("/api/v1/connectors", payload);
  return data;
}
