export default function Dashboard() {
return (
<div className="p-8"><h1 className="text-2xl font-bold">
<p className="mt-2 text-gray-600"><Ping />

)
}
function Ping() {
const [status, setStatus] = React.useState('…')
React.useEffect(() => {
fetch('/api/health').then(r => r.text()).then(setStatus).catch(() => setStatus('down'))
}, [])
return <span className={status === 'ok' ? 'text-green-600' : 'text-red-600'}>{status}
}
