import { useEffect, useState } from "react"

function App() {
  const [logs, setLogs] = useState<string[]>([])

  useEffect(() => {
    // Listen for process_log events from Go backend
    if (window.runtime) {
      window.runtime.EventsOn("process_log", (log: string) => {
        setLogs((prevLogs) => [...prevLogs, log]) // Append new log to list
      })
    }
  }, [])

  return (
    <div className="min-h-screen bg-white grid place-items-center mx-auto py-8">
      <div className="text-blue-900 text-2xl font-bold flex flex-col items-center space-y-4">
        <h1>Process Monitor</h1>
        
        {/* Process Logs Section */}
        <div className="w-full max-w-md h-64 overflow-y-auto bg-gray-100 p-4 rounded-lg border mt-4">
          <h2 className="text-lg font-semibold mb-2">Process Logs:</h2>
          <div className="text-sm text-gray-800">
            {logs.length > 0 ? (
              logs.map((log, index) => <div key={index} className="mb-1">{log}</div>)
            ) : (
              <p className="text-gray-500">No logs yet...</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
